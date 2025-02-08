## sourcegraph-cody
Thư mục chứa script tự động trả lời và các prompt cho AI. Chứa log câu trả lời AI để dễ debug hơn

## translate
Chứa log bản dịch tiếng Anh sang Việt và script tự động dịch. Hiện giờ thì không cần nữa tuy nhiên trong tương lai để tối ưu có lẽ sẽ cần. Tôi nghĩ thế chứ hiện giờ vẫn chưa có ý tưởng sử dụng

## api.go
File này chứa chứa các đoạn code về logic xử lý api. Do là tôi sản phẩm này được tôi phát triển cả frontend và backend nên đặt phương thức là post hết tuy nhiên nếu bạn muốn code api thì dùng đúng phương thức để xác định từng loại api để người code frontend dễ nhận biết các xử lý từng api. **Tôi đã từng bị góp ý (nói tránh nói giảm)** 

#### Cấu trúc truy vấn API
Đầu tiên hãy lấy token từ đường dẫn /login với cấu trúc

```bash 
{
  "Username": "username",
  "Password": "password"
}
```
Cách tạo thì nó nằm trong thư mục **/enviroment/.env**
Sau khi gửi truy vấn API sẽ trả về kết quả với dạng 

``` bash
{
    "token": "eyJhbxxxx..."
}
```
**Lưu ý đây là dạng token động vì vậy nó sẽ tự thay đổi sau một khoảng thời gian**

Sau khi đã lấy xong token hãy truy cập các API như **/word**
#### Lưu ý để thêm phần xác thực API thì trong thành phần https của bạn phải thêm key = Authorization, value = Bearer + mã token sau khi cấu hình header xong thì hãy gọi api với dạng


```bash
{
  "Data": "nội dung cần tra từ điển hoặc phân tích câu"
}
```

Sau khi gửi truy vấn thì sẽ trả về dạng

```bash
Success 
```
Nếu mà sai thì nó sẽ trả về status lỗi (phần này tôi chưa có làm trường hợp sai chỉ in ra vài trường hợp lỗi như token sai hoặc in ra trong log của server)

Với websocket **/ChatCody**
#### Lưu ý để thêm phần xác thực websocket thì trong thành phần wss của bạn phải thêm key = Sec-WebSocket-Protocol, value = mã token sau khi cấu hình header. Hãy chuyển đổi giao thức từ https sang wss.

Dữ liệu sẽ tự gửi qua cho frontend sau khi thực hiện gọi truy vấn api **/word** với dạng

```bash 
   Detail interface{}
   Structure string 
```

Ví dụ output

``` 
{
    "detail": {
        "* Cụm động từ:": {
            "+ Che đậy hoặc làm mờ bằng mực đen (14)": {
                "EN": [
                    "  + Ví dụ: They blacked out the confidential information (14) (EN)"
                ],
                "VI": [
                    "  + Ví dụ: Họ bôi đen những thông tin mật (14) (VI)"
                ]
            },
            "+ Cấm phổ biến thông tin, đặc biệt qua kiểm duyệt (15)": {
                "EN": [
                    "  + Ví dụ: The government blacks out the news (15) (EN)"
                ],
                "VI": [
                    "  + Ví dụ: Chính phủ kiểm duyệt tin tức (15) (VI)"
                ]
            },
            "+ Tạm thời mất ý thức hoặc trí nhớ (13)": {
                "EN": [
                    "  + Ví dụ: He blacked out from heat stroke (13) (EN)"
                ],
                "VI": [
                    "  + Ví dụ: Anh ấy ngất đi vì say nắng (13) (VI)"
                ]
            },
            "+ Tắt hết đèn để tránh máy bay địch phát hiện (16)": {
                "EN": [
                    "  + Ví dụ: The city blacks out during air raids (16) (EN)"
                ],
                "VI": [
                    "  + Ví dụ: Thành phố tắt đèn phòng không (16) (VI)"
                ]
            }
        },
        "* Danh từ:": {
            "+ Chất màu hoặc thuốc nhuộm có màu đen (6)": {
                "EN": [
                    "  + Ví dụ: Black dye (6) (EN)"
                ],
                "VI": [
                    "  + Ví dụ: Thuốc nhuộm màu đen (6) (VI)"
                ]
            },
            "+ Màu đen, màu tối nhất không phản chiếu ánh sáng (5)": {
                "EN": [
                    "  + Ví dụ: Black is my favorite color (5) (EN)"
                ],
                "VI": [
                    "  + Ví dụ: Màu đen là màu yêu thích của tôi (5) (VI)"
                ]
            },
            "+ Quần áo màu đen, đặc biệt là đồ tang (8)": {
                "EN": [
                    "  + Ví dụ: Wearing black for mourning (8) (EN)"
                ],
                "VI": [
                    "  + Ví dụ: Mặc đồ đen để tang (8) (VI)"
                ]
            },
            "+ Sự vắng mặt hoàn toàn của ánh sáng; bóng tối (7)": {
                "EN": [
                    "  + Ví dụ: The room was in black (7) (EN)"
                ],
                "VI": [
                    "  + Ví dụ: Bóng tối bao trùm căn phòng (7) (VI)"
                ]
            }
        },
        "* Nội động từ:": {
            "+ Làm cho đen (9)": {
                "EN": [
                    "  + Ví dụ: The smoke blacks the wall (9) (EN)"
                ],
                "VI": [
                    "  + Ví dụ: Khói làm tường nhà đen (9) (VI)"
                ]
            },
            "+ Trở nên đen (12)": {
                "EN": [
                    "  + Ví dụ: His skin blacks in the sun (12) (EN)"
                ],
                "VI": [
                    "  + Ví dụ: Da anh ấy đen đi vì nắng (12) (VI)"
                ]
            },
            "+ Tẩy chay như một hành động công đoàn (11)": {
                "EN": [
                    "  + Ví dụ: Workers black the company (11) (EN)"
                ],
                "VI": [
                    "  + Ví dụ: Công nhân tẩy chay công ty (11) (VI)"
                ]
            },
            "+ Đánh xi đen (10)": {
                "EN": [
                    "  + Ví dụ: To black the shoes (10) (EN)"
                ],
                "VI": [
                    "  + Ví dụ: Đánh xi giày đen (10) (VI)"
                ]
            }
        },
        "* Tính từ:": {
            "+ Có màu đen, phản chiếu ít ánh sáng và không có màu sắc chủ đạo (1)": {
                "EN": [
                    "  + Ví dụ: The black shirt (1) (EN)"
                ],
                "VI": [
                    "  + Ví dụ: Chiếc áo đen (1) (VI)"
                ]
            },
            "+ Có ít hoặc không có ánh sáng (2)": {
                "EN": [
                    "  + Ví dụ: The black sky (2) (EN)"
                ],
                "VI": [
                    "  + Ví dụ: Bầu trời đen tối (2) (VI) "
                ]
            },
            "+ Thuộc về chủng tộc có màu da nâu đến đen, đặc biệt là người gốc Phi (3)": {
                "EN": [
                    "  + Ví dụ: Black people (3) (EN)"
                ],
                "VI": [
                    "  + Ví dụ: Người da đen (3) (VI)"
                ]
            },
            "+ Thuộc về nhóm dân tộc Mỹ gốc Phi có màu da sẫm (4)": {
                "EN": [
                    "  + Ví dụ: Black American culture (4) (EN)"
                ],
                "VI": [
                    "  + Ví dụ: Văn hóa người Mỹ da đen (4) (VI)"
                ]
            }
        }
    },
    "structure": "WordClass"
}
```

## chat_cody.go
File này tôi sẽ nói tóm tắt cụ thể lọc dữ liệu lại cho sạch đã lấy từ wordnik sau đó thêm vài chức năng chuyển đổi model cho AI và đảm bảo AI trả lời đúng cấu trúc bản thân muốn, rồi tự động chạy script.

## function.go
Chứa những đoạn hàm mà hay dùng thường xuyên nói chung là đoạn nào thường xài nhiều chức năng đó thì tôi code về 1 hàm cho khỏi code lại

## go.mod && go.sum
Chứa các module, thư viện đã import vào project

## json_web_token (jwt) 
Đoạn code này tạo token động cho các api chủ yếu nhằm tránh tình trạng og nào lấy api spam lung tung dẫn đến server bị crack. (Nói chung thì nếu muốn bảo mật hơn dùng thêm session tạo thêm nhiều thuộc tính để xác thực). Như tôi lười quá nên chỉ dùng mỗi token động thôi. Hãy kham khảo thêm từ nguồn [github](https://github.com/golang-jwt/jwt)

## listen.txt
File này là file đọc thành tiếng cho các chuỗi nhưng do nó chỉ hoạt động trên mỗi server thôi vì hiện tại tôi đã public nên tạm thời để đó và chuyển sang sử dụng api của [responsivevoice](https://responsivevoice.org/) để mỗi trình duyệt client có thể sử dụng chức năng chuyển đổi text sang âm thanh (chủ yếu tạo chức năng nghe từ vựng mỗi khi client trỏ vào từ vựng)

## look_up_word
Kết hợp lại từ việc gọi API wordnik lấy dữ liệu rồi lọc dữ liệu sau đó đưa về cho AI trả lời (AI chủ yếu lấy các ví dụ của mỗi từ vựng vì có vài từ vựng bên wordnik không có ví dụ)

## struct_data.go
Cứ hiểu đơn giản là tôi nhóm cấu trúc dữ liệu lại cho gọn với đặt sang 1 file riêng để dễ tìm thôi. Ông nào rành về hướng đối tượng là sẽ rõ

## struct_json_api.go
Tương tự như **struct_data.go** chỉ khác là nó là cầu nối giữa backend và frontend

## struct_json.go
Tương tự như **struct_data.go** chỉ khác là nó cầu nối giữa backend của tôi với wordnik (vì tôi có xài api của wordnik nữa, nói cách khác backend của tôi đang là client gọi api từ wordnik)

## web_socket.go
Định nghĩa tương tự như socket tuy nhiên khác một tý là nó là cầu nối giao tiếp giữa server và client thôi tức là cách hoạt động của nó sẽ là bên **A mỗi khi xử lý dữ liệu xong là tự động lập tức gửi dữ liệu cho bên B** khác với code bình thường ở chỗ A gửi dữ liệu cho B xong B lại phải báo lại bên A.