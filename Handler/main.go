package main

import (
	"fmt"
)


const socketPath = "./tmp/translation_complete.sock"

func main() {
	fmt.Println("Bắt đầu chương trình")

	create_socket(socketPath, "./translate/auto_translate.sh")

	fmt.Println("Đang xử lý từ...")
	define_word("the")

	fmt.Println("Chương trình kết thúc")
	fmt.Println()
}
