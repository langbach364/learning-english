## list_option_callAPI
File này chứa các thông tin của các thuộc tính có thể gọi API wordnik

## sourcegraph-cody
Thư mục này chứa script tự động trả lời với các file cấu trúc như dữ liệu, câu trả lời, form output (form mẫu trả lởi), form gốc, model chọn, prompt để cho AI trả lời theo cấu trúc

## chat_cody.go
File này tôi sẽ nói tóm tắt cụ thể lọc dữ liệu lại cho sạch đã lấy từ wordnik sau đó thêm vài chức năng chuyển đổi model cho AI và đảm bảo AI trả lời đúng cấu trúc bản thân muốn, rồi tự động chạy script.

## create_schedule.go
File này chứa logic xử lý cho việc tạo lịch tự động

## function.go
Chứa những đoạn hàm mà hay dùng thường xuyên nói chung là đoạn nào thường xài nhiều chức năng đó thì tôi code về 1 hàm cho khỏi code lại

## get_word.go
File này chứa các đoạn code xử lý việc lấy từ vựng thông qua API wordnik

## go.mod && go.sum
Chứa các module, thư viện đã import vào project

## graph_API.go
File này chứa các đoạn code xử lý cho API theo dạng graphql

#### Dạng câu select

```
query {
  select {
    ${table}(limit: ${number}) {
        ${columns}
        ...
    }
  }
}
```

#### Dạng câu insert
```
mutation {
  insert {
    ${table}($columns: ${values}, .... )
  }
}
```
#### Dạng câu delete
```
mutation {
  delete {
    ${table}(${columm}: ${value})
  }
}

```
Với dự án hiện tại thì chỉ dùng 1 tham số word cho tất cả bảng nếu muốn cấu hình lại thì hãy tìm đến đoạn

```
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

```

hãy cấu hình thêm tham số ở Args nếu muốn

```
Args: graphql.FieldConfigArgument{
			"word": &graphql.ArgumentConfig{
				Type: graphql.NewNonNull(graphql.String),
			},
		},
```

## Hiện chưa có cập nhật đầy đủ cho API này vì đang thiếu xác thực API

#### Cập nhật sau

## rest_API.go
File này chứa các đoạn code xử lý cho API theo dạng restful

Thư mục chỉ chứa logic chạy code tuy nhiên các endpoint của API được cấu hình trong **./main.go** 

## Hiện chưa có cập nhật đầy đủ cho API này vì đang thiếu xác thực API

#### Cập nhật sau

## struct_data.go
File này chứa các struct dữ liệu để lưu trữ dữ liệu từ API wordnik

## struct_json,go
File này chứa cấu trúc dữ liệu để cho API wordnik và cấu trúc truy vấn api thông qua json từ frontend

## vocabulary_statistics.go
File này chứa thống kê từ vựng đã học trong khoảng thời gian là ngày, tháng, năm