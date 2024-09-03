package main

import (
	"fmt"
	"net"
	"os"
	"os/exec"
)


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

// Đợi luồng công cụ xử lý
func wait_tool_complete(socketPath string) error {
    listener, err := net.Listen("unix", socketPath)
    if err != nil {
        return fmt.Errorf("lỗi khi tạo socket: %v", err)
    }
    defer listener.Close()

    conn, err := listener.Accept()
    if err != nil {
        return fmt.Errorf("lỗi khi chấp nhận kết nối: %v", err)
    }
    defer conn.Close()

    buffer := make([]byte, 1024)
    _, err = conn.Read(buffer)
    if err != nil {
        return fmt.Errorf("lỗi khi đọc từ socket: %v", err)
    }

    return nil
}