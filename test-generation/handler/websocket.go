package main

import (
	"bytes"
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

func HandleWebSocket(c echo.Context) error {
	ws, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		log.Printf("❌ Lỗi nâng cấp kết nối: %v", err)
		return err
	}
	defer ws.Close()

	log.Println("🔌 Kết nối WebSocket mới được thiết lập")

	for {
		_, msg, err := ws.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("❌ Lỗi đọc message: %v", err)
			}
			break
		}

		// Ghi log message nhận được
		log.Printf("📥 Nhận được message: %s", string(msg))

		// Làm sạch dữ liệu JSON
		cleanMsg := bytes.TrimSpace(msg)

		var message WebSocketMessage
		if err := json.Unmarshal(cleanMsg, &message); err != nil {
			log.Printf("❌ Lỗi parse JSON: %v", err)
			sendErrorResponse(ws, "Lỗi định dạng JSON: "+err.Error())
			continue
		}

		switch message.Handle {
		case "start":
			handleStartProcess(ws, message.Data)
		default:
			sendErrorResponse(ws, "Không hỗ trợ handle này")
		}
	}

	return nil
}

func handleStartProcess(ws *websocket.Conn, data []map[string]string) {
	log.Println("🚀 Bắt đầu xử lý dữ liệu từ vựng")

	if len(data) == 0 {
		sendErrorResponse(ws, "Không có dữ liệu từ vựng")
		return
	}

	processedData := make(map[string][]string)

	for _, item := range data {
		if word, exists := item["word"]; exists && word != "" {
			if _, ok := processedData["word"]; !ok {
				processedData["word"] = []string{}
			}
			processedData["word"] = append(processedData["word"], word)
		}
	}

	if len(processedData["word"]) == 0 {
		sendErrorResponse(ws, "Không tìm thấy từ vựng hợp lệ")
		return
	}

	sendProgressResponse(ws, "Đang xử lý dữ liệu với Cody...")

	model := "openai::2024-02-01::gpt-4o"
	chat_cody(processedData, model)

	sendProgressResponse(ws, "Đang phân tích kết quả từ Cody...")

	result := handler_data()

	response := WebSocketResponse{
		Status:  "success",
		Message: "Xử lý dữ liệu thành công",
		Data:    result,
	}

	sendResponse(ws, response)
	log.Println("✅ Hoàn thành xử lý dữ liệu từ vựng")
}

func sendErrorResponse(ws *websocket.Conn, message string) {
	response := WebSocketResponse{
		Status:  "error",
		Message: message,
	}
	sendResponse(ws, response)
}

func sendProgressResponse(ws *websocket.Conn, message string) {
	response := WebSocketResponse{
		Status:  "progress",
		Message: message,
	}
	sendResponse(ws, response)
}

func sendResponse(ws *websocket.Conn, response WebSocketResponse) {
	responseJSON, err := json.Marshal(response)
	if err != nil {
		log.Printf("❌ Lỗi chuyển đổi response thành JSON: %v", err)
		return
	}

	if err := ws.WriteMessage(websocket.TextMessage, responseJSON); err != nil {
		log.Printf("❌ Lỗi gửi response: %v", err)
	}
}

func SetupWebSocket(e *echo.Echo) {
	e.GET("/ws", HandleWebSocket)
	log.Println("📡 Đã thiết lập endpoint WebSocket tại /ws")
}
