package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"
)

func check_text(text string) bool {
	words := strings.Fields(text)
	return len(words) <= 1
}

func middleware_Word(filePath string) {
	fileChanged := make(chan bool)
	var checkText string
	go watch_file(filePath, fileChanged)

	for range fileChanged {
		func() {
			file, err := os.Open(filePath)
			if err != nil {
				log.Printf("Không thể mở file: %v", err)
				return
			}
			defer file.Close()

			scanner := bufio.NewScanner(file)
			var text string

			for scanner.Scan() {
				text = scanner.Text()
			}

			if text == "" || checkText == text {
				return
			}

			checkText = text

			if err := scanner.Err(); err != nil {
				log.Printf("Lỗi khi đọc file: %v", err)
				return
			}

			data := make(map[string][]string)
			data["Câu"] = []string{text}
			if check_text(text) {
				data = result_definitions(text)
				fmt.Println(text)
			}
			chat_cody(data, "openai/gpt-4o")
			dataStructure := data_structure()

			if dataSocket["ChatCody"] == nil {
				dataSocket["ChatCody"] = make([]interface{}, 0)
			}

			dataSocket["ChatCody"] = []interface{}{dataStructure}
			log.Print(dataSocket["ChatCody"])
			broadCast["ChatCody"] <- true
		}()	}
}


func middleware_listen_word(filePath string) {
	var processMutex sync.Mutex
	fileChanged := make(chan bool)
	go watch_file(filePath, fileChanged)

	for range fileChanged {
		if !processMutex.TryLock() {
			continue
		}
		go func() {
			defer processMutex.Unlock()

			file, err := os.Open(filePath)
			if err != nil {
				log.Printf("Không thể mở file: %v", err)
				return
			}
			defer file.Close()

			scanner := bufio.NewScanner(file)
			var text string
			for scanner.Scan() {
				text = scanner.Text()
			}

			if text == "" {
				return
			}

			get_data("LangBach", "en")
			fmt.Println("Đã có âm thanh")
		}()
	}
}

func main() {
	data := path_file()

	go middleware_Word(data["word"])
	go middleware_listen_word(data["listen_word"])
	create_server()

	// data := data_structure()
	// fmt.Println(data)
}
