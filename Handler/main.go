package main

import (
)

func main() {
    data := result_definitions("cat")
	
    chat_cody(data, "anthropic/claude-3-5-sonnet-20240620", "./tmp/cody.sock")
}
