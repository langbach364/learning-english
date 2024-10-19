package main

import (
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

var (
	upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	clients      = make(map[*websocket.Conn]bool)
	clientsMutex = &sync.Mutex{}
)

func connect_client(w http.ResponseWriter, r *http.Request) *websocket.Conn {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Lỗi khi thực hiện kết nối WebSocket:", err)
		return nil
	}
	return conn
}

func send_message(conn *websocket.Conn, message []byte) {
	clientsMutex.Lock()
	defer clientsMutex.Unlock()

	err := conn.WriteMessage(websocket.TextMessage, message)
	if err != nil {
		log.Println("Lỗi khi gửi tin nhắn:", err)
		conn.Close()
		delete(clients, conn)
	}
}

func handler_boardcast(key string, boardcast *map[string]chan bool, jsonData []byte, conn *websocket.Conn) {
	if jsonData == nil {
		log.Println("Không tìm thấy dữ liệu JSON")
		return
	}
	
	send_message(conn, jsonData)
	(*boardcast)[key] <- false
}

func boardcast_socket(conn *websocket.Conn, boardcast map[string]chan bool, data_json map[string][]byte) {
	for {
		for key, ch := range boardcast {
			select {
			case <-ch:
				switch key {
				case "AnsCody":
					{
						jsonData := data_json[key]
						handler_boardcast(key, &boardcast, jsonData, conn)
					}
				case "ListenWord":
					{
						jsonData := data_json[key]
						handler_boardcast(key, &boardcast, jsonData, conn)
					}
				}
				default: 
			}
		}
	}
}
func synthetic_websocket(w http.ResponseWriter, r *http.Request) {
	conn := connect_client(w, r)
	if conn == nil {
		return
	}

	clientsMutex.Lock()
    clients[conn] = true
    clientsMutex.Unlock()

	go boardcast_socket(conn, broadcast, data_json)
}