package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var broadCast = make(map[string]chan bool)
var dataSocket = make(map[string][]interface{})

// Websocket connection
// ///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
var (
	upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	clients = make(map[*websocket.Conn][]string)
)

func connect_client(w http.ResponseWriter, r *http.Request) *websocket.Conn {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("upgrade:", err)
		return nil
	}
	return conn
}

func recive_message_client(conn *websocket.Conn) string {
	defer func() {
		if r := recover(); r != nil {
			log.Println("Client ngắt kết nối đột ngột:", conn.RemoteAddr())
			conn.Close()
			delete(clients, conn)
			return
		}
	}()

	_, messageClient, err := conn.ReadMessage()
	if err != nil {
		log.Println("Client mất kết nối tại:", conn.RemoteAddr())
		conn.Close()
		delete(clients, conn)
		return ""
	}

	log.Printf("Dữ liệu client %s là: %s", conn.RemoteAddr(), string(messageClient))

	return string(messageClient)
}

func send_message_client(conn *websocket.Conn, message interface{}) {
	defer func() {
		if r := recover(); r != nil {
			log.Println("Client ngắt kết nối đột ngột:", conn.RemoteAddr())
			conn.Close()
			delete(clients, conn)
		}
	}()
	err := conn.WriteJSON(message)
	fmt.Println("Dữ liệu gửi đến client là:", message)

	if err != nil {
		log.Println("Client đã mất kết nối tại địa chỉ:", conn.RemoteAddr())
		log.Panicln("Lỗi: ", err)
		conn.Close()
		delete(clients, conn)
	}
}

func handle_websocket(nameEvent string) {
	go func() {
		for <-broadCast[nameEvent] {
			for conn := range clients {

				data := dataSocket[nameEvent]
				dataJson := AnswerData{
					Detail:    data[0],
					Structure: dataStructure,
				}

				log.Printf("Gửi dữ liệu cho client: %s\n", conn.RemoteAddr())
				send_message_client(conn, dataJson)
			}
		}
	}()
}

func check_alive_client(conn *websocket.Conn) {
	go func() {
		recive_message_client(conn)
	}()
}

func handle_connection(nameEvent string, conn *websocket.Conn) {
	log.Printf("Client đã kết nối tại: %s", conn.RemoteAddr())
	clients[conn] = append(clients[conn], nameEvent)
	if broadCast[nameEvent] == nil {
		broadCast[nameEvent] = make(chan bool)
	}

	handle_websocket(nameEvent)
	check_alive_client(conn)
}
