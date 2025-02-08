## auto_install.sh
Lúc tôi xài arch thì cần tải các công cụ này để dịch thuật lại các dữ liệu lấy ở từ điển thông qua API tuy nhiên thì hiện tại có lẽ không cần sử dụng nhưng tôi vẫn để đây vì biết đâu để tối ưu khả năng xử lý tôi cần dùng tới

## models.txt
Chứa các thuộc tính để sử dụng mô hình AI từ nguồn [Sourcegraph](https://sourcegraph.com/)
Cách tải và triển khai dạng cody cli thì có thể tham khảo qua thư mục /learning-english/docker/be.dockerfile

## data
Thư mục này chứa các dữ liệu hồi lúc mới lên ý tưởng xây dự án này tuy nhiên do vượt quá dự đoán số vốn từ vựng vì thế tôi đã thay thế data từ [Wordnik](https://www.wordnik.com/) và dùng AI từ [Sourcegraph](https://sourcegraph.com/) để lấy từ vựng và lấy ví dụ. Tuy nhiên tôi vẫn để ở đó vì hiện giờ tôi đang có ý tưởng trong việc lấy data của từ thông qua các lần xử lý rồi lưu lại về sau truy vấn lại. Nó giống tương tự như cache lưu dữ liệu cho lần truy vấn sau không cần xử lý lại tuy nhiên tôi sẽ thực hiện ý tưởng đó sau khi đã hoàn thành bản beta cho dự án

## docker
Thư mục chứa file cấu hình và triển khai, môi trường test, sản phẩm trên các server. Vẫn chưa hoàn toàn tự động hóa việc triển khai môi trường test và môi trường sản phẩm. Có lẽ trong tương lai tôi sẽ làm hoàn toàn tự động hóa cho việc triển khai dễ dàng hơn

## frontend
Chứa các đoạn code giao diện thôi và tôi tách khá nhiều thành phần vì dùng AI để code hoàn toàn giao diện lmao. Tôi tự nhận nếu là về giao diện tôi sẽ code thua cả AI vì thế tôi đã tận dụng sức mạnh AI để code giao diện vì thế có những đoạn code nó khá là kinh khủng và khó hiểu xin thông cảm

## handler
Thư mục chứa các phần xử lý luồng rồi logic xử lý của hệ thống do chính tay tôi code. Nó còn thiếu rất nhiều đoạn log chi tiết để in thông tin xử lý luồng hệ thống nếu có bug hoặc lưu các thông tin xử lý

## middleware
Chứa các file txt đây là những dữ liệu từ client gửi về cho hệ thống xử lý. Trong tương lai tôi sẽ cập nhật chuyển đổi sang json để tăng độ chính xác và cấu trúc rõ ràng hơn nhưng vì hồi đầu xây dựng tôi chưa nghĩ tới @_@

