package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/graphql-go/graphql"
	"github.com/graphql-go/handler"
	"github.com/iancoleman/orderedmap"
)

var (
	dbPool *sql.DB
	tables map[string]*orderedmap.OrderedMap
)

func init() {
	tables = init_tables()
}

func init_tables() map[string]*orderedmap.OrderedMap {
	data := make(map[string]*orderedmap.OrderedMap)

	vocabulary := orderedmap.New()
	vocabulary.Set("word", "string")
	vocabulary.Set("frequency", "int")
	vocabulary.Set("error_count", "int")
	data["vocabulary"] = vocabulary

	schedule := orderedmap.New()
	schedule.Set("id", "int")
	schedule.Set("time", "datetime")
	schedule.Set("word", "string")
	data["schedule"] = schedule

	return data
}

func init_DB() error {
	db, err := sql.Open("mysql", "root:@ztegc4df9f4e@tcp(localhost:3306)/learned_vocabulary")
	if err != nil {
		return err
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * time.Minute)

	dbPool = db
	return nil
}

func get_graphQL_type(fieldType string) graphql.Type {
	typeMap := map[string]graphql.Type{
		"string":   graphql.String,
		"int":      graphql.Int,
		"datetime": graphql.DateTime,
		"float":    graphql.Float,
		"boolean":  graphql.Boolean,
	}

	if t, exists := typeMap[fieldType]; exists {
		return t
	}
	return graphql.String
}

func create_graphQL_fields(tableInfo *orderedmap.OrderedMap) graphql.Fields {
	fields := graphql.Fields{}
	for _, key := range tableInfo.Keys() {
		fieldType, _ := tableInfo.Get(key)
		fields[key] = &graphql.Field{
			Type: get_graphQL_type(fieldType.(string)),
		}
	}
	return fields
}

func get_table_name(tableInfo *orderedmap.OrderedMap) string {
	for name, info := range tables {
		if compare_table_info(tableInfo, info) {
			return name
		}
	}
	return ""
}

func compare_table_info(t1, t2 *orderedmap.OrderedMap) bool {
	for _, key := range t1.Keys() {
		val1, _ := t1.Get(key)
		val2, _ := t2.Get(key)
		if val1 != val2 {
			return false
		}
	}
	return true
}

func build_query(tableInfo *orderedmap.OrderedMap, tableName string, limit int, operation string) string {
	columns := tableInfo.Keys()
	operation = strings.ToUpper(strings.TrimSpace(operation))

	switch operation {
	case "DELETE":
		return fmt.Sprintf("DELETE FROM %s WHERE word = ?", tableName)
	case "INSERT":
		placeholders := strings.Repeat("?,", len(columns))
		placeholders = placeholders[:len(placeholders)-1]
		return fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)",
			tableName,
			strings.Join(columns, ","),
			placeholders)
	default: // SELECT
		return fmt.Sprintf("SELECT %s FROM %s LIMIT %d",
			strings.Join(columns, ","),
			tableName,
			limit)
	}
}

func process_rows(rows *sql.Rows) ([]map[string]interface{}, error) {
	columns, err := rows.Columns()
	if err != nil {
		return nil, err
	}

	values := make([]interface{}, len(columns))
	valuePtrs := make([]interface{}, len(columns))
	result := make([]map[string]interface{}, 0)

	for i := range columns {
		valuePtrs[i] = &values[i]
	}

	for rows.Next() {
		if err := rows.Scan(valuePtrs...); err != nil {
			return nil, err
		}

		entry := make(map[string]interface{})
		for i, col := range columns {
			if val, ok := values[i].([]byte); ok {
				entry[col] = string(val)
			} else {
				entry[col] = values[i]
			}
		}
		result = append(result, entry)
	}

	return result, nil
}

func select_data(tableInfo *orderedmap.OrderedMap, limit int) *graphql.Field {
	return &graphql.Field{
		Type: graphql.NewList(graphql.NewObject(graphql.ObjectConfig{
			Name:   get_table_name(tableInfo) + "Type",
			Fields: create_graphQL_fields(tableInfo),
		})),
		Args: graphql.FieldConfigArgument{
			"limit": &graphql.ArgumentConfig{Type: graphql.Int},
		},
		Resolve: func(p graphql.ResolveParams) (interface{}, error) {
			if val, ok := p.Args["limit"].(int); ok {
				limit = val
			}

			query := build_query(tableInfo, get_table_name(tableInfo), limit, "SELECT")
			rows, err := dbPool.Query(query)
			if err != nil {
				return nil, err
			}
			defer rows.Close()

			return process_rows(rows)
		},
	}
}

func delete_data(tableInfo *orderedmap.OrderedMap) *graphql.Field {
	args := make(graphql.FieldConfigArgument)
	for _, key := range tableInfo.Keys() {
		fieldType, _ := tableInfo.Get(key)
		args[key] = &graphql.ArgumentConfig{Type: get_graphQL_type(fieldType.(string))}
	}

	return &graphql.Field{
		Type: graphql.Boolean,
		Args: args,
		Resolve: func(p graphql.ResolveParams) (interface{}, error) {
			query := build_query(tableInfo, get_table_name(tableInfo), 0, "DELETE")
			result, err := dbPool.Exec(query, p.Args["word"])
			if err != nil {
				return false, err
			}
			affected, _ := result.RowsAffected()
			return affected > 0, nil
		},
	}
}

func insert_data(tableInfo *orderedmap.OrderedMap) *graphql.Field {
	args := make(graphql.FieldConfigArgument)
	for _, key := range tableInfo.Keys() {
		fieldType, _ := tableInfo.Get(key)
		args[key] = &graphql.ArgumentConfig{Type: get_graphQL_type(fieldType.(string))}
	}

	return &graphql.Field{
		Type: graphql.Boolean,
		Args: args,
		Resolve: func(p graphql.ResolveParams) (interface{}, error) {
			fields := tableInfo.Keys()
			query := build_query(tableInfo, get_table_name(tableInfo), 0, "INSERT")

			values := make([]interface{}, len(fields))
			for i, field := range fields {
				values[i] = p.Args[field]
			}

			result, err := dbPool.Exec(query, values...)
			if err != nil {
				return false, err
			}
			affected, _ := result.RowsAffected()
			return affected > 0, nil
		},
	}
}

func enable_graphQL(port, pattern string, limitQuery int) {
	if err := init_DB(); err != nil {
		log.Fatal("Không thể kết nối database:", err)
	}
	defer dbPool.Close()

	fields := graphql.Fields{
    "select": &graphql.Field{
        Type: graphql.NewObject(graphql.ObjectConfig{
            Name: "SelectType",
            Fields: graphql.Fields{
                "vocabulary": select_data(tables["vocabulary"], limitQuery),
                "schedule": select_data(tables["schedule"], limitQuery),
            },
        }),
    },
    "delete": &graphql.Field{
        Type: graphql.NewObject(graphql.ObjectConfig{
            Name: "DeleteType",
            Fields: graphql.Fields{
                "vocabulary": delete_data(tables["vocabulary"]),
                "schedule": delete_data(tables["schedule"]),
            },
        }),
    },
    "insert": &graphql.Field{
        Type: graphql.NewObject(graphql.ObjectConfig{
            Name: "InsertType", 
            Fields: graphql.Fields{
                "vocabulary": insert_data(tables["vocabulary"]),
                "schedule": insert_data(tables["schedule"]),
            },
        }),
    },
}

	schema, err := graphql.NewSchema(graphql.SchemaConfig{
		Query: graphql.NewObject(graphql.ObjectConfig{
			Name:   "RootQuery",
			Fields: fields,
		}),
	})

	if err != nil {
		log.Fatal("Lỗi khởi tạo schema:", err)
	}

	h := handler.New(&handler.Config{
		Schema:   &schema,
		Pretty:   true,
		GraphiQL: true,
	})

	http.Handle("/"+pattern, h)
	fmt.Printf("Server đang chạy tại http://localhost%s/%s\n", port, pattern)
	log.Fatal(http.ListenAndServe(port, nil))
}
