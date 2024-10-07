package main

import (
	"bufio"
	"fmt"
	"strings"
)

func add_data(data map[string][]string) bool {
	file, err := write_file("./sourcegraph-cody/data.txt")
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
	file, err := write_file("./sourcegraph-cody/model.txt")
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

func chat_cody(data map[string][]string, model string) {
	start_chat(data, model, "./sourcegraph-cody/cody.sh")

	file, err := read_file("./sourcegraph-cody/answer.txt")
	if err != nil {
		fmt.Println("Lỗi khi đọc file: ", err)
	}
	defer file.Close()
}

func get_key_line(line string) string {
	x := strings.Index(line, ":")
	return line[:x+1]
}

func get_value_line(line string) string {
	x := strings.Index(line, ":")
	return strings.TrimSpace(line[x+1:])
}

func process_line(line string) int {
	if len(line) == 0 {
        return -1
    }

	switch line[0] {
	case '*':
		{
			return 0
		}
	case '+':
		{
			return 1
		}
	}
	return -1
}

func handler_data(data map[string][]string, line string, saveKey *string) {
	switch process_line(line) {
	case 0:
		{
			key := get_key_line(line)
			value := get_value_line(line)

			if len(value) > 0 {
				data[key] = append(data[key], value)
				return
			}
			*saveKey = key
		}
	case 1:
		{
			data[*saveKey] = append(data[*saveKey], line)
		}
	default:
		{
			{
				fmt.Println("Lỗi khi xử lý dữ liệu")
			}
		}
	}
}

func skip_line(line string) bool {
	return strings.Contains(line, "[") && strings.Contains(line, "]")
}

func data_structure() map[string][]string{
	pathFile, err := read_file("./sourcegraph-cody/answer.txt")
	check_err(err)
	defer pathFile.Close()

	scanner := bufio.NewScanner(pathFile)
	data := make(map[string][]string)
	var saveKey string

	for scanner.Scan() {
		line := scanner.Text()
		if skip_line(line) {
			continue
		}
		handler_data(data, line, &saveKey)
	}
	return data
}
