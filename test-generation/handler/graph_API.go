package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
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

// Tạo kết nối database
// Create database connection
func init_DB() error {
	db, err := sql.Open("mysql", "root:@ztegc4df9f4e@tcp(localhost:3306)/learned_vocabulary")
	if err != nil {
		return err
	}

	db.SetConnMaxLifetime(5 * time.Minute)
	dbPool = db
	return nil
}

// Lấy cấu trúc schema của các bảng trong database
// Get table schema from database
func get_table_schema(db *sql.DB, dbName string) (map[string]*orderedmap.OrderedMap, error) {
	data := make(map[string]*orderedmap.OrderedMap)

	query := `SELECT TABLE_NAME 
             FROM INFORMATION_SCHEMA.TABLES 
             WHERE TABLE_SCHEMA = ?`

	rows, err := db.Query(query, dbName)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var tableName string
		if err := rows.Scan(&tableName); err != nil {
			return nil, err
		}

		colQuery := `SELECT COLUMN_NAME, DATA_TYPE 
                    FROM INFORMATION_SCHEMA.COLUMNS 
                    WHERE TABLE_SCHEMA = ? AND TABLE_NAME = ?`

		colRows, err := db.Query(colQuery, dbName, tableName)
		if err != nil {
			return nil, err
		}
		defer colRows.Close()

		tableInfo := orderedmap.New()
		fields := make(map[string]string)

		for colRows.Next() {
			var colName, dataType string
			if err := colRows.Scan(&colName, &dataType); err != nil {
				return nil, err
			}

			switch dataType {
			case "int", "tinyint", "smallint", "mediumint", "bigint":
				fields[colName] = "int"
			case "varchar", "text", "char":
				fields[colName] = "string"
			case "float", "double", "decimal":
				fields[colName] = "float"
			case "datetime", "timestamp":
				fields[colName] = "datetime"
			default:
				fields[colName] = "string"
			}
		}

		for k, v := range fields {
			tableInfo.Set(k, v)
		}
		data[tableName] = tableInfo
	}

	return data, nil
}

// Chuyển đổi kiểu dữ liệu sang kiểu GraphQL tương ứng
// Convert data type to corresponding GraphQL type
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

// Tạo các trường GraphQL từ thông tin bảng
// Create GraphQL fields from table information
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

// Xây dựng câu truy vấn SQL dựa trên operation
// Build SQL query based on operation type
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

// Xử lý kết quả truy vấn từ database
// Process query results from database
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

// Tạo resolver cho operation SELECT
// Create resolver for SELECT operation
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
			rows, err := dbPool.Query(query)
			if err != nil {
				return nil, err
			}

			defer rows.Close()
			return process_rows(rows)
		},
	}
}

// Tạo resolver cho operation DELETE
// Create resolver for DELETE operation
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
			result, err := dbPool.Exec(query)
			if err != nil {
				return false, err
			}
			affected, _ := result.RowsAffected()
			return affected > 0, nil
		},
	}
}

// Tạo resolver cho operation INSERT
// Create resolver for INSERT operation
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

// Tạo schema operation cho GraphQL
// Create schema operation for GraphQL
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

// Tạo các trường operation cho GraphQL
// Create operation fields for GraphQL
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

// Khởi tạo và cấu hình GraphQL server
// Initialize and configure GraphQL server
func enable_graphQL(port, pattern string, limitQuery int) {
	if err := init_DB(); err != nil {
		log.Fatal("Không thể kết nối database:", err)
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		<-c
		if dbPool != nil {
			dbPool.Close()
		}
		os.Exit(0)
	}()

	var err error
	tables, err = get_table_schema(dbPool, "learned_vocabulary")
	if err != nil {
		log.Fatal("Không thể lấy schema:", err)
	}

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
	go func() {
		if err := http.ListenAndServe(port, nil); err != nil {
			log.Printf("GraphQL server error: %v", err)
		}
	}()
}
