package main

import (
	"bufio"
	"fmt"
)

func add_data(data map[string][]string) {
	file, err := write_file("./sourcegraph-cody/data.txt")
	if err != nil {
		fmt.Println("Lỗi khi tạo file: ", err)
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
}


func add_model(model string)  {
	file, err := write_file("./sourcegraph-cody/model.txt")
	if err != nil {
		fmt.Println("Lỗi khi tạo file: ", err)
	}
	defer file.Close()

	_, err = file.WriteString(model)
	if err != nil {
		fmt.Println("Lỗi khi ghi vào file: ", err)
	}

}

func chat_cody(data map[string][]string, model string, pathSocket string) string {
	add_data(data)
	add_model(model)
	
	create_socket(pathSocket)
	run_script("./sourcegraph-cody/cody.sh")

	file, err := read_file("./sourcegraph-cody/answer.txt")
	if err != nil {
		fmt.Println("Lỗi khi đọc file: ", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		answer := scanner.Text()
		return answer
	}
	return ""
}
