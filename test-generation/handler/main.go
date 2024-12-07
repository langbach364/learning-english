package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

func main() {
	// Khởi động GraphQL server
	go enable_graphQL(":8080", "graphql", 10)

	// Khởi động REST API server
	go enable_rest("8081", "/words")

	// Khởi động scheduling system
	scheduling_word()

	// Đợi servers khởi động
	time.Sleep(2 * time.Second)

	// Test data để gửi qua REST API
	testWords := []infoWord{
		{Word: "book", Frequency: 5, ErrorCount: 1},
		{Word: "study", Frequency: 3, ErrorCount: 2},
		{Word: "learn", Frequency: 0, ErrorCount: 4},
	}

	// Test gửi từ qua REST API
	for _, word := range testWords {
		jsonData, err := json.Marshal(word)
		if err != nil {
			log.Printf("Lỗi marshal JSON: %v", err)
			continue
		}

		resp, err := http.Post("http://localhost:8081/words",
			"application/json",
			bytes.NewBuffer(jsonData))

		if err != nil {
			log.Printf("Lỗi gửi request: %v", err)
			continue
		}

		if resp.StatusCode == http.StatusOK {
			fmt.Printf("Thêm từ thành công: %s\n", word.Word)
		} else {
			fmt.Printf("Lỗi khi thêm từ %s: %d\n", word.Word, resp.StatusCode)
		}
		resp.Body.Close()

		// Đợi một chút giữa các request
		time.Sleep(500 * time.Millisecond)
	}

	// Đợi để xem kết quả
	time.Sleep(2 * time.Second)

	select {}
}
