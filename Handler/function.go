package main

import (
	"fmt"
	"net"
	"os"
	"os/exec"
	"time"
)


func write_file(fileName string) (*os.File, error) {
	file, err := os.Create(fileName)
	if err != nil {
		return nil, fmt.Errorf("lỗi khi tạo file: %v", err)
	}
	time.Sleep(100 * time.Millisecond) 
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

func create_socket(socketPath string) error {
    fmt.Printf("Đã tạo socket %s\n", socketPath)
    listener, err := net.Listen("unix", socketPath)
    if err != nil {
        return fmt.Errorf("lỗi khi tạo socket: %v", err)
    }
	fmt.Println("Socket đã được tạo thành công")
    defer listener.Close()

	fmt.Println("Socket đã được tự động xóa khi đóng listener.Close")
    return nil
}