package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
)

var (
	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
)

func Handle_web_socket(c echo.Context) error {
	ws, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		log.Printf("❌ Lỗi nâng cấp kết nối: %v", err)
		return err
	}
	defer ws.Close()

	log.Println("🔌 Kết nối WebSocket mới được thiết lập")

	for {
		messageType, msg, err := ws.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("❌ Lỗi đọc message: %v", err)
			} else {
				log.Printf("🔌 Kết nối WebSocket đóng: %v", err)
			}
			break
		}

		if messageType != websocket.TextMessage {
			log.Printf("⚠️ Nhận được message không phải dạng text: %d", messageType)
			continue
		}

		log.Printf("📥 Nhận được message: %s", string(msg))

		var message WebSocketMessageRevised
		if err := json.Unmarshal(msg, &message); err != nil {
			log.Printf("❌ Lỗi parse JSON: %v. Dữ liệu nhận được: %s", err, string(msg))
			sendErrorResponse(ws, "Lỗi định dạng JSON: "+err.Error()+". Vui lòng kiểm tra cấu trúc {handle: '...', data: { 'new_words': [...], 'old_words': [...]}}")
			continue
		}

		switch message.Handle {
		case "start":
			handle_start_process(ws, message.Data)
		default:
			log.Printf("🚫 Handle không được hỗ trợ: %s", message.Handle)
			sendErrorResponse(ws, "Không hỗ trợ handle này: "+message.Handle)
		}
	}

	log.Println("🔌 Đóng kết nối WebSocket")
	return nil
}

func handle_start_process(ws *websocket.Conn, data map[string][]string) {
	log.Println("🚀 Bắt đầu xử lý dữ liệu từ vựng (mới/cũ)")

	if data == nil {
		log.Println("❌ Dữ liệu từ vựng rỗng")
		sendErrorResponse(ws, "Dữ liệu từ vựng không được rỗng")
		return
	}

	newWords, newExists := data["newWords"]
	oldWords, oldExists := data["oldWords"]

	if !newExists || len(newWords) == 0 {
		log.Println("❌ Không tìm thấy hoặc không có 'newWords' hợp lệ")
		sendErrorResponse(ws, "Thiếu hoặc không có dữ liệu 'newWords'")
	}
	if !oldExists || len(oldWords) == 0 {
		log.Println("❌ Không tìm thấy hoặc không có 'oldWords' hợp lệ")
		sendErrorResponse(ws, "Thiếu hoặc không có dữ liệu 'oldWords'")
	}

	if (!newExists || len(newWords) == 0) && (!oldExists || len(oldWords) == 0) {
		log.Println("❌ Không có dữ liệu từ vựng hợp lệ nào (cả mới và cũ)")
		sendErrorResponse(ws, "Không có dữ liệu từ vựng hợp lệ nào được cung cấp")
		return
	}

	log.Printf("📊 Dữ liệu nhận được: %d từ mới, %d từ cũ", len(newWords), len(oldWords))

	processedData := data

	sendProgressResponse(ws, "Đang xử lý dữ liệu với Cody...")

	model := "anthropic::2024-10-22::claude-3-7-sonnet-extended-thinking"

	test_generator(processedData, model)

	sendProgressResponse(ws, "Đang phân tích kết quả từ Cody...")

	result := handler_data()

	response := WebSocketResponse{
		Status:  "success",
		Message: "Xử lý dữ liệu thành công",
		Data:    result,
	}

	send_response(ws, response)
	log.Println("✅ Hoàn thành xử lý dữ liệu từ vựng (mới/cũ)")
}

func sendErrorResponse(ws *websocket.Conn, message string) {
	response := WebSocketResponse{
		Status:  "error",
		Message: message,
	}
	send_response(ws, response)
}

func sendProgressResponse(ws *websocket.Conn, message string) {
	response := WebSocketResponse{
		Status:  "progress",
		Message: message,
	}
	send_response(ws, response)
}

func send_response(ws *websocket.Conn, response WebSocketResponse) {
	responseJSON, err := json.Marshal(response)
	if err != nil {
		log.Printf("❌ Lỗi chuyển đổi response thành JSON: %v", err)
		return
	}

	if err := ws.WriteMessage(websocket.TextMessage, responseJSON); err != nil {
		log.Printf("❌ Lỗi gửi response: %v", err)
	} else {
		log.Printf("📤 Đã gửi response: %s", string(responseJSON))
	}
}

func Setup_web_socket(e *echo.Echo) {
	e.GET("/ws", Handle_web_socket)
	log.Println("📡 Đã thiết lập endpoint WebSocket tại /ws")
}
