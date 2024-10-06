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

func process_line(line string, data *AnswerData, currentKey *string) {
    if len(line) == 0 {
        return
    }

    content := strings.TrimSpace(line)

    switch {
    case strings.HasPrefix(content, "[Câu]:"):
        *currentKey = "Câu"
        if data.Details == nil {
            data.Details = make(map[string][]string)
        }
        value := strings.TrimSpace(strings.TrimPrefix(content, "[Câu]:"))
        if value != "" {
            data.Details[*currentKey] = append(data.Details[*currentKey], value)
        }
    case strings.HasPrefix(content, "* "):
        parts := strings.SplitN(content[2:], ":", 2)
        *currentKey = strings.TrimSpace(parts[0])
        if data.Details == nil {
            data.Details = make(map[string][]string)
        }
        if len(parts) > 1 {
            value := strings.TrimSpace(parts[1])
            if value != "" {
                data.Details[*currentKey] = append(data.Details[*currentKey], value)
            }
        }
    case strings.HasPrefix(content, "+ "):
        value := strings.TrimSpace(strings.TrimPrefix(content, "+ "))
        if *currentKey != "" && value != "" {
            data.Details[*currentKey] = append(data.Details[*currentKey], value)
        }
    default:
        if *currentKey != "" && content != "" {
            data.Details[*currentKey] = append(data.Details[*currentKey], content)
        }
    }
}

func data_structure() AnswerData {
    pathFile, err := read_file("./sourcegraph-cody/answer.txt")
    check_err(err)

    scanner := bufio.NewScanner(pathFile)
    var data AnswerData
    var currentKey string

    data.Details = make(map[string][]string)

    for scanner.Scan() {
        line := scanner.Text()
        process_line(line, &data, &currentKey)
    }

    if err := scanner.Err(); err != nil {
        fmt.Println("Lỗi khi đọc file:", err)
    }

    return data
}
