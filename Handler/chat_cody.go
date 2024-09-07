package main

import (
	
	"fmt"
)

func add_data(data map[string]string) {
	file, err := write_file("./sourcegraph-cody/data.txt")
	if err != nil {
		fmt.Println("Lỗi khi tạo file: ", err)
	}
	defer file.Close()

	for key, value := range data {
		_, err = file.WriteString(fmt.Sprintf("%s: %s\n", key, value))
		if err != nil {
			fmt.Println("Lỗi khi ghi vào file: ", err)
		}
	}
}

func add_model(model string)  {
	file, err := write_file("./sourcegraph-cody/answer.txt")
	if err != nil {
		fmt.Println("Lỗi khi tạo file: ", err)
	}
	defer file.Close()

	_, err = file.WriteString(model)
	if err != nil {
		fmt.Println("Lỗi khi ghi vào file: ", err)
	}

}

// func chat_cody(data map[string]string, model string, pathSocket string) string {
// 	add_data(data)
// 	add_model(model)
	
// 	create_socket(pathSocket, "./sourcegraph-cody/cody.sh")
// 	file, err := read_file("./sourcegraph-cody/answer.txt")
// 	if err != nil {
// 		fmt.Println("Lỗi khi đọc file: ", err)
// 	}
// 	defer file.Close()

// 	scanner := bufio.NewScanner(file)
// 	for scanner.Scan() {
// 		answer := scanner.Text()
// 		return answer
// 	}
// 	return ""

// }
