package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strings"

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

func setTables(nameTable *orderedmap.OrderedMap, fields map[string]string) *orderedmap.OrderedMap {
	for k, v := range fields {
		nameTable.Set(k, v)
	}
	return nameTable
}

func init_tables() map[string]*orderedmap.OrderedMap {
	data := make(map[string]*orderedmap.OrderedMap)

	vocabulary := orderedmap.New()
	fieldsVocabulary := map[string]string{
		"word":        "string",
		"frequency":   "int",
		"error_count": "int",
	}
	setTables(vocabulary, fieldsVocabulary)
	data["vocabulary"] = vocabulary

	schedule := orderedmap.New()
	fieldsSchedule := map[string]string{
		"id":   "string",
		"time": "string",
		"word": "string",
	}
	setTables(schedule, fieldsSchedule)
	data["schedule"] = schedule

	return data
}

func init_DB() error {
	db, err := sql.Open("mysql", "root:@ztegc4df9f4e@tcp(localhost:3306)/learned_vocabulary")
	if err != nil {
		return err
	}
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

func build_query(tableInfo *orderedmap.OrderedMap, tableName string, limit int, operation, word string) string {
	columns := tableInfo.Keys()
	operation = strings.ToUpper(strings.TrimSpace(operation))

	switch operation {
	case "DELETE":
		return fmt.Sprintf("DELETE FROM %s WHERE word = '%s'", tableName, word)
	case "INSERT":
		placeholders := strings.Repeat("?,", len(columns))
		placeholders = placeholders[:len(placeholders)-1]
		return fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)",
			tableName,
			strings.Join(columns, ","),
			placeholders)
	default:
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

func select_data(tableName string, tableInfo *orderedmap.OrderedMap, limit int) *graphql.Field {
	return &graphql.Field{
		Type: graphql.NewList(graphql.NewObject(graphql.ObjectConfig{
			Name:   fmt.Sprintf("%sSelectType", tableName),
			Fields: create_graphQL_fields(tableInfo),
		})),
		Args: graphql.FieldConfigArgument{
			"limit": &graphql.ArgumentConfig{Type: graphql.Int},
		},
		Resolve: func(p graphql.ResolveParams) (interface{}, error) {
			if val, ok := p.Args["limit"].(int); ok {
				limit = val
			}

			query := build_query(tableInfo, tableName, limit, "SELECT", "")
			fmt.Println(query)

			rows, err := dbPool.Query(query)
			if err != nil {
				return nil, err
			}

			defer rows.Close()
			return process_rows(rows)
		},
	}
}

func delete_data(tableName string, tableInfo *orderedmap.OrderedMap) *graphql.Field {
	return &graphql.Field{
		Type: graphql.Boolean,
		Args: graphql.FieldConfigArgument{
			"word": &graphql.ArgumentConfig{
				Type: graphql.NewNonNull(graphql.String),
			},
		},
		Resolve: func(p graphql.ResolveParams) (interface{}, error) {
			word := p.Args["word"].(string)

			query := build_query(tableInfo, tableName, 0, "DELETE", word)
			fmt.Println(query)

			result, err := dbPool.Exec(query)
			if err != nil {
				return false, err
			}

			affected, _ := result.RowsAffected()
			return affected > 0, nil
		},
	}
}

func insert_data(tableName string, tableInfo *orderedmap.OrderedMap) *graphql.Field {
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

			query := build_query(tableInfo, tableName, 0, "INSERT", "")
			fmt.Println(query)

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

func create_schema_operation(operation string, resolvers graphql.Fields) *graphql.Field {
	return &graphql.Field{
		Type: graphql.NewObject(graphql.ObjectConfig{
			Name:   fmt.Sprintf("%sOperation", operation),
			Fields: resolvers,
		}),
		Resolve: func(p graphql.ResolveParams) (interface{}, error) {
			return map[string]interface{}{}, nil
		},
	}
}

func create_operation_fields(operation string, tables map[string]*orderedmap.OrderedMap, limitQuery int) graphql.Fields {
	fields := graphql.Fields{}

	for tableName, tableInfo := range tables {
		switch operation {
		case "Select":
			fields[tableName] = select_data(tableName, tableInfo, limitQuery)
		case "Delete":
			fields[tableName] = delete_data(tableName, tableInfo)
		case "Insert":
			fields[tableName] = insert_data(tableName, tableInfo)
		}
	}
	return fields
}

func enable_graphQL(port, pattern string, limitQuery int) {
	if err := init_DB(); err != nil {
		log.Fatal("Không thể kết nối database:", err)
	}
	defer dbPool.Close()

	fieldsQuery := graphql.Fields{
		"select": create_schema_operation("Select", create_operation_fields("Select", tables, limitQuery)),
	}

	fieldsMutation := graphql.Fields{
		"delete": create_schema_operation("Delete", create_operation_fields("Delete", tables, 0)),
		"insert": create_schema_operation("Insert", create_operation_fields("Insert", tables, 0)),
	}

	schemaRoot, err := graphql.NewSchema(graphql.SchemaConfig{
		Query: graphql.NewObject(graphql.ObjectConfig{
			Name:   "RootQuery",
			Fields: fieldsQuery,
		}),
		Mutation: graphql.NewObject(graphql.ObjectConfig{
			Name:   "RootMutation",
			Fields: fieldsMutation,
		}),
	})

	if err != nil {
		log.Fatal("Lỗi khởi tạo schema:", err)
	}

	h := handler.New(&handler.Config{
		Schema:   &schemaRoot,
		Pretty:   true,
		GraphiQL: true,
	})

	http.Handle("/"+pattern, h)
	fmt.Printf("Server đang chạy tại http://localhost%s/%s\n", port, pattern)
	log.Fatal(http.ListenAndServe(port, nil))
}
