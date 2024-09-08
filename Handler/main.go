package main

import (
	"fmt"
)

func main() {
	data := result_definitions("make")
	fmt.Println(chat_cody(data, "anthropic/claude-3-5-sonnet-20240620", "./tmp/cody.sock"))
}
