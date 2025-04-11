package main

import (
	"fmt"
)


func add_data(data map[string][]string) bool {
	file, err := write_file("./sourcegraph-cody/test-generator/data.txt")
	if err != nil {
		fmt.Println("Lỗi khi tạo file: ", err)
		return false
	}
	defer file.Close()

	file.WriteString("&&&\n")
	for key, values := range data {
		file.WriteString(key + ":\n\n")
		for _, value := range values {
			file.WriteString("- " + value + "\n\n")
		}
		file.WriteString("-------------------------------------\n\n")
	}
	file.WriteString("&&&")
	return true
}

func add_model(model string) bool {
	file, err := write_file("./sourcegraph-cody/test-generator/model.txt")
	if err != nil {
		fmt.Println("Lỗi khi tạo file: ", err)
		return false
	}
	defer file.Close()

	_, err = file.WriteString(model)
	if err != nil {
		fmt.Println("Lỗi khi ghi vào file: ", err)
		return false
	}

	return true
}

func start_chat(data map[string][]string, model, scriptName string) bool {
	chan_data := make(chan bool)
	chan_model := make(chan bool)

	go func() {
		if !add_data(data) {
			chan_data <- false
			fmt.Println("Lỗi khi thêm dữ liệu vào cody")
			return
		}
		chan_data <- true
	}()

	go func() {
		if !add_model(model) {
			chan_model <- false
			fmt.Println("Lỗi khi thêm model vào cody")
			return
		}
		chan_model <- true
	}()

	if <-chan_data && <-chan_model {
		run_script(scriptName)
		return true
	}
	return false
}

func test_generator(data map[string][]string, model string) {
	start_chat(data, model, "./sourcegraph-cody/test-generator/cody.sh")

	file, err := read_file("./sourcegraph-cody/test-generator/answer.txt")
	if err != nil {
		fmt.Println("Lỗi khi đọc file: ", err)
	}
	defer file.Close()
}
