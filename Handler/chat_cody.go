package main

import (
	"bufio"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"unicode"
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

func check_data_structure(line string) int {
	startChar := strings.Index(line, "[")
	endChar := strings.Index(line, "]")

	switch line[startChar+1 : endChar] {
	case "Câu":
		{
			return 1
		}
	case "Từ loại":
		{
			return 2
		}
	}
	return -1
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
	line = strings.TrimSpace(line)
	return len(line) == 0
}

func check_key_or_value(data string) int {
	re := regexp.MustCompile(`\((\w+)\)`)
	matches := re.FindAllStringSubmatch(data, -1)

	if len(matches) > 1 {
		return 1 // value
	}
	return 2 // key
}

func get_number_string(data string) string {
	re := regexp.MustCompile(`\((\w+)\)`)
	matches := re.FindAllStringSubmatch(data, -1)

	if len(matches) > 0 && len(matches[0]) > 1 {
  		for i, match := range matches {
  			if _, err := strconv.Atoi(match[1]); err == nil {
  				return matches[i][1]
  			}
  		}
	}
	return ""
}

func get_language_code_string(data string) string {
	re := regexp.MustCompile(`\((\w+)\)`)
	matches := re.FindAllStringSubmatch(data, -1)

	if len(matches) > 1 && len(matches[1]) > 1 {
		for _, match := range matches {
			if strings.IndexFunc(match[1], unicode.IsLetter) != -1 {
                return match[1]
            }
		}
	}
	return ""
}

func handler_data_word_class(data map[string]map[string]map[string][]string, line string, saveKey *string) {
	switch process_line(line) {
	case 0:
		{
			value := get_value_line(line)
			if len(value) == 0 {
				key := get_key_line(line)
				data[key] = make(map[string]map[string][]string)
				*saveKey = key
			}
		}
	case 1:
		{
			check := check_key_or_value(line)
			value := line
			if check == 2 {
				data[*saveKey][value] = make(map[string][]string, 0)
			} else {
				for key1, value1 := range data {
					for key2 := range value1 {
						indexKey := get_number_string(key2)
						indexValue := get_number_string(value)
						if indexKey == indexValue {
							codeLa := get_language_code_string(value)
							data[key1][key2][codeLa] = append(data[key1][key2][codeLa], value)
						}
					}
				}
			}
		}
	}
}

func data_structure() DataStructure {
	pathFile, err := read_file("./sourcegraph-cody/answer.txt")
	check_err(err)
	defer pathFile.Close()

	scanner := bufio.NewScanner(pathFile)
	var result DataStructure
	var data DataStructure
	check := 0
	var saveKey string

	for scanner.Scan() {
		line := scanner.Text()

		if skip_line(line) {
			continue
		}
		if check == 0 {
			check = check_data_structure(line)
			switch check {
			case 1:
				data = make(map[string][]string)
			case 2:
				data = make(map[string]map[string]map[string][]string)
			}
		}

		switch check {
		case 1:
			handler_data(data.(map[string][]string), line, &saveKey)
		case 2:
			handler_data_word_class(data.(map[string]map[string]map[string][]string), line, &saveKey)
		default:
			fmt.Println("Không đúng trường hợp nào")
		}
	}

	result = data
	return result
}