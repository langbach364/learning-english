package main

import (
	"bufio"
	"fmt"
	"regexp"
)

func process_line(line string) string {
	re := regexp.MustCompile(`Phần (\d+)`)
	match := re.FindStringSubmatch(line)

	if len(match) > 1 {
		return match[1]
	}

	return ""
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
    re := regexp.MustCompile(`\*\s+(.+?):\s*(.*)`)
    match := re.FindStringSubmatch(line)

    if len(match) < 2 {
        return "", ""
    }

    key := match[1]
    value := ""
    if len(match) > 2 {
        value = match[2]
    }

    return key, value
}

func handler_data() map[string]map[string][]string {
	pathFile, err := read_file("./sourcegraph-cody/answer.txt")
	check_err(err)
	defer pathFile.Close()

	scanner := bufio.NewScanner(pathFile)

	data := make(map[string]map[string][]string)
	saveNumberPart := ""
	saveKey2 := ""

	for scanner.Scan() {
		line := scanner.Text()
		numberPart := process_line(line)

		if numberPart != "" {
			data[numberPart] = make(map[string][]string)
			saveNumberPart = numberPart
		} else {
			checkCase := check_key_or_value(line)
			switch checkCase {
			case 1:
				{
					key, value := get_key_value(line)
					data[saveNumberPart][key] = append(data[saveNumberPart][key], value)
					if value == "" {
						saveKey2 = key
					}
				}
			case 2:
				{
					data[saveNumberPart][saveKey2] = append(data[saveNumberPart][saveKey2], line)
				}
			case 3:
				{
					data[saveNumberPart][saveKey2] = append(data[saveNumberPart][saveKey2], line)
				}
			}
		}
	}

	return data
}

func print_data(data map[string]map[string][]string) {
	for key, value := range data {
		fmt.Println("Part:", key)
		for key2, value2 := range value {
			fmt.Println("Key:", key2)
			for _, value3 := range value2 {
				fmt.Println("Value:", value3)
			}
		}
	}
}
