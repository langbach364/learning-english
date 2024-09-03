package main

import (
	"fmt"
	"os"
	"path/filepath"
)


const socketPath = "./tmp/translation_complete.sock"

func main() {
	fmt.Println("Bắt đầu chương trình")

	tmpDir := filepath.Dir(socketPath)
	os.MkdirAll(tmpDir, os.ModePerm)

	if _, err := os.Stat(socketPath); err == nil {
		os.Remove(socketPath)
	}

	go run_script("./translate/auto_translate.sh")

	fmt.Println("Đang xử lý từ...")
	define_word("remake")

	fmt.Println("Chương trình kết thúc")
	fmt.Println()
}
