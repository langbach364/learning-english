package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/fsnotify/fsnotify"
)

func check_err(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func write_file(fileName string) (*os.File, error) {
	file, err := os.Create(fileName)
	if err != nil {
		return nil, fmt.Errorf("lỗi khi tạo file: %v", err)
	}
	return file, nil
}

func read_file(fileName string) (*os.File, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return nil, fmt.Errorf("lỗi khi mở file: %v", err)
	}
	return file, nil
}

func run_script(scriptName string) {
	cmd := exec.Command("bash", scriptName)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		fmt.Printf("Lỗi khi chạy %s: %v\n", scriptName, err)
	}
}

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