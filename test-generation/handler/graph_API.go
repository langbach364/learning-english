package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"reflect"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	"github.com/graphql-go/graphql"
	"github.com/graphql-go/handler"
)

// Thiết lập kết nối với database
// Establish database connection
func connectDB() (*sql.DB, error) {
	return sql.Open("mysql", "root:@ztegc4df9f4e@tcp(localhost:3306)/learned_vocabulary")
}

func get_graphQL_type(fieldType string) graphql.Type {
	switch fieldType {
	case "string":
		return graphql.String
	case "int":
		return graphql.Int
	case "datetime":
		return graphql.DateTime
	case "float":
		return graphql.Float
	case "boolean":
		return graphql.Boolean
	default:
		return graphql.String
	}
}

// Định nghĩa schema database với các bảng và kiểu dữ liệu của trường
// Define database schema with tables and their field types
func data_table() map[string]map[string]string {
	data := map[string]map[string]string{
		"vocabulary": {
			"word":        "string",
			"frequency":   "int",
			"error_count": "int",
		},
		"schedule": {
			"id":   "int",
			"time": "datetime",
			"word": "string",
		},
	}
	return data
}

// Tạo định nghĩa các trường GraphQL từ schema database
// Generate GraphQL field definitions from database schema
func create_graphQL_fields(tableInfo map[string]string) graphql.Fields {
	fields := graphql.Fields{}
	for fieldName, fieldType := range tableInfo {
		fields[fieldName] = &graphql.Field{
			Type: get_graphQL_type(fieldType),
		}
	}
	return fields
}

// Lấy tên các bảng từ schema database đã định nghĩa
// Retrieve table names from defined database schema
func get_table_name(tableInfo map[string]string) string {
	tables := data_table()
	for name, info := range tables {
		if reflect.DeepEqual(info, tableInfo) {
			return name
		}
	}
	return ""
}

// Xây dựng chuỗi truy vấn SQL từ cấu hình bảng
// Construct SQL query string from table configuration
func build_query(tableInfo map[string]string, tableName string, limit int) string {
	columns := make([]string, 0)
	for colName := range tableInfo {
		columns = append(columns, colName)
	}
	return fmt.Sprintf("SELECT %s FROM %s LIMIT %d",
		strings.Join(columns, ","),
		tableName,
		limit,
	)
}

// Chuyển đổi dữ liệu từ các dòng database sang cấu trúc map
// Transform database rows into map structure
func process_rows(rows *sql.Rows) ([]map[string]interface{}, error) {
	var result []map[string]interface{}
	columns, _ := rows.Columns()
	count := len(columns)
	values := make([]interface{}, count)
	valuePtrs := make([]interface{}, count)

	for rows.Next() {
		for i := range columns {
			valuePtrs[i] = &values[i]
		}
		rows.Scan(valuePtrs...)

		entry := make(map[string]interface{})
		for i, col := range columns {
			val := values[i]
			if byteArray, ok := val.([]byte); ok {
				entry[col] = string(byteArray)
			} else {
				entry[col] = val
			}
		}
		result = append(result, entry)
	}
	return result, nil
}

// Định nghĩa GraphQL resolver cho việc truy xuất dữ liệu
// Define GraphQL query resolver for data selection
func select_data(tableInfo map[string]string, limit int) *graphql.Field {
	return &graphql.Field{
		Type: graphql.NewList(graphql.NewObject(
			graphql.ObjectConfig{
				Name:   get_table_name(tableInfo) + "Type",
				Fields: create_graphQL_fields(tableInfo),
			},
		)),
		Args: graphql.FieldConfigArgument{
			"limit": &graphql.ArgumentConfig{
				Type: graphql.Int,
			},
		},
		Resolve: func(p graphql.ResolveParams) (interface{}, error) {
			if val, ok := p.Args["limit"].(int); ok {
				limit = val
			}

			tableName := get_table_name(tableInfo)

			db, err := connectDB()
			if err != nil {
				return nil, err
			}
			defer db.Close()

			query := build_query(tableInfo, tableName, limit)
			rows, err := db.Query(query)
			if err != nil {
				return nil, err
			}
			defer rows.Close()

			return process_rows(rows)
		},
	}
}

func enable_graph_sql(port, pattern string) {
	tables := data_table()
	limitQuery := 5
	fields := graphql.Fields{}

	for tableName, tableInfo := range tables {
		fields[tableName] = select_data(tableInfo, limitQuery)
	}

	rootQuery := graphql.ObjectConfig{Name: "RootQuery", Fields: fields}
	schemaConfig := graphql.SchemaConfig{Query: graphql.NewObject(rootQuery)}
	schema, err := graphql.NewSchema(schemaConfig)
	if err != nil {
		log.Fatal("Lỗi khi tạo schema:", err)
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
