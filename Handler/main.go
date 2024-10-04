package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/fsnotify/fsnotify"
)

func watch_file(fileName string, fileChanged chan<- bool) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	err = watcher.Add(filepath.Dir(fileName))
	if err != nil {
		log.Fatal(err)
	}

	for {
		select {
		case event, ok := <-watcher.Events:
			if !ok {
				return
			}
			if event.Op&fsnotify.Write == fsnotify.Write && event.Name == fileName {
				fileChanged <- true
			}
		case err, ok := <-watcher.Errors:
			if !ok {
				return
			}
			log.Println("error:", err)
		}
	}
}

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
			chat_cody(data, "openai/gpt-4o", "./tmp/cody.sock")
		}()
	}
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

func print_structured_data(data map[string][]string) {
    for key, values := range data {
        fmt.Printf("[key] %s: ", key)
        for _, value := range values {
            fmt.Printf("[value] %s\n", value)
        }
        fmt.Println()
    }
}


func main() {
	data := path_file()
	go middleware_Word(data["word"])
	go middleware_listen_word(data["listen_word"])
    create_server()
}
