package main

import (
	"bufio"
	"fmt"
	"regexp"
)

func process_line(line string) (string, string) {
	re := regexp.MustCompile(`Phần (\d+):\s*(.+)`)
	match := re.FindStringSubmatch(line)

	if len(match) > 2 {
		key := "Phần " + match[1]
		value := match[2]
		return key, value
	}

	return "", ""
}

func check_key_or_value(line string) int {
	var char string

	if len(line) >= 2 {
		char = line[:2]
	}

	re := regexp.MustCompile(`\s+`)
	char = re.ReplaceAllString(char, "")

	switch char {
	case "*":
		return 1
	case "+":
		return 2
	case "->":
		return 3
	}

	return 0
}

func get_key_value(line string) (string, string) {
	re := regexp.MustCompile(`\*\s+([\w-]+\s+\([\w-]+\)):\s+(.+)\.`)
	match := re.FindStringSubmatch(line)

	if len(match) > 2 {
		key := match[1]
		value := match[2]
		return key, value
	}

	return "", ""
}

func add_data_processed(check int, data *map[string][]Word, line string) map[string][]Word {
	numberPart, content := process_line(line)

	switch check {
	case 0:
		if numberPart != "" {
			key := numberPart + content
			(*data)[key] = append((*data)[key], Word{})
		}
	case 1:
		key, value := get_key_value(line)
		(*data)[key] = append((*data)[key], Word{WordMeaning: []string{value}})
	case 2:
		key, value := get_key_value(line)
		(*data)[key] = append((*data)[key], Word{Remember: []string{value}})
	case 3:
		key, value := get_key_value(line)
		(*data)[key] = append((*data)[key], Word{Remember: []string{value}})
	}

	return *data
}

func handler_data() map[string][]Word{
	pathFile, err := read_file("./sourcegraph-cody/answer.txt")
	check_err(err)
	defer pathFile.Close()

	scanner := bufio.NewScanner(pathFile)
	dataWord := make(map[string][]Word)
	for scanner.Scan() {
		line := scanner.Text()

		check := check_key_or_value(line)
		dataWord = add_data_processed(check, &dataWord, line)
	}

	return dataWord
}

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
