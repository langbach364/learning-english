package main

import (
	"bufio"
	"fmt"
	"regexp"
	"strings"

	"github.com/elliotchance/orderedmap"
)

var saveKey2 string
var DataStructure = map[string][]string{
	"word":     {"Phần 1", "Phần 2"},
	"sentence": {"Phần 3", "Phần 4"},
}

func check_part_type(numberPart string) string {
	for dataType, parts := range DataStructure {
		for _, part := range parts {
			if part == numberPart {
				return dataType
			}
		}
	}
	return ""
}

func merge_data(dest map[string]map[string][]string, src map[string]*orderedmap.OrderedMap) {
	for key, value := range src {
		if _, exists := dest[key]; !exists {
			dest[key] = make(map[string][]string)
		}
		iter := value.Front()
		for iter != nil {
			k := iter.Key.(string)
			v := iter.Value.([]string)
			dest[key][k] = append(dest[key][k], v...)
			iter = iter.Next()
		}
	}
}

func process_line(line string) string {
	re := regexp.MustCompile(`^Phần (\d+):?`)
	match := re.FindStringSubmatch(line)

	if len(match) > 1 {
		return "Phần " + match[1]
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
	re := regexp.MustCompile(`\*\s+([^:]+):\s*(.*)`)
	match := re.FindStringSubmatch(line)

	if len(match) < 2 {
		return "", ""
	}

	key := strings.TrimSpace(match[1])
	value := ""
	if len(match) > 2 {
		value = strings.TrimSpace(match[2])
	}

	return key, value
}

func handler_word(numberPart *string, line string, saveNumberPart string) map[string]*orderedmap.OrderedMap {
	check := check_part_type(*numberPart)
	if check == "sentence" {
		return nil
	}

	data := make(map[string]*orderedmap.OrderedMap)

	if (*numberPart) != "" {
		data[(*numberPart)] = orderedmap.NewOrderedMap()
		saveNumberPart = *numberPart
		*numberPart = ""
	}

	if _, exists := data[saveNumberPart]; !exists {
		data[saveNumberPart] = orderedmap.NewOrderedMap()
	}

	checkCase := check_key_or_value(line)
	switch checkCase {
	case 1:
		{
			key, value := get_key_value(line)
			if key != "" {
				saveKey2 = key
				if value != "" {
					values := []string{value}
					data[saveNumberPart].Set(key, values)
				} else {
					data[saveNumberPart].Set(key, []string{})
				}
			}
		}
	case 2:
		{
			if saveKey2 != "" {
				if values, exists := data[saveNumberPart].Get(saveKey2); exists {
					data[saveNumberPart].Set(saveKey2, append(values.([]string), line))
				} else {
					data[saveNumberPart].Set(saveKey2, []string{line})
				}
			}
		}
	case 3:
		{
			if saveKey2 != "" {
				if values, exists := data[saveNumberPart].Get(saveKey2); exists {
					data[saveNumberPart].Set(saveKey2, append(values.([]string), line))
				} else {
					data[saveNumberPart].Set(saveKey2, []string{line})
				}
			}
		}
	}

	return data
}

func handler_sentence(numberPart *string, line string) map[string]*orderedmap.OrderedMap {
	check := check_part_type(*numberPart)
	if check == "word" {
		return nil
	}

	data := make(map[string]*orderedmap.OrderedMap)
	data[*numberPart] = orderedmap.NewOrderedMap()

	data[*numberPart].Set("vietnamese", []string{})
	data[*numberPart].Set("english", []string{})

	checkCase := check_key_or_value(line)

	switch checkCase {
	case 1:
		{
			_, value := get_key_value(line)
			if values, exists := data[*numberPart].Get("english"); exists {
				data[*numberPart].Set("english", append(values.([]string), value))
			} else {
				data[*numberPart].Set("english", []string{value})
			}
		}
	case 2:
		{
			re := regexp.MustCompile(`\+\s+(.+)`)
			match := re.FindStringSubmatch(line)
			if len(match) > 1 {
				value := match[1]
				if values, exists := data[*numberPart].Get("vietnamese"); exists {
					data[*numberPart].Set("vietnamese", append(values.([]string), value))
				} else {
					data[*numberPart].Set("vietnamese", []string{value})
				}
			}
		}
	}
	return data
}

func handler_data() map[string]map[string][]string {
	pathFile, err := read_file("./sourcegraph-cody/test-generator/answer.txt")
	check_err(err)
	defer pathFile.Close()

	scanner := bufio.NewScanner(pathFile)
	data := make(map[string]map[string][]string)

	var currentPart string
	var saveNumberPart string

	for scanner.Scan() {
		line := scanner.Text()
		part := process_line(line)

		if part != "" {
			currentPart = part
			saveNumberPart = currentPart
			if _, exists := data[currentPart]; !exists {
				data[currentPart] = make(map[string][]string)
			}
			continue
		}

		merge_data(data, handler_word(&currentPart, line, saveNumberPart))
		// merge_data(data, handler_sentence(&currentPart, line))
	}

	return data
}

func printf_log_data(data map[string]map[string][]string) {
	fmt.Println("\n=== PRINTING DATA STRUCTURE ===")
	fmt.Printf("Tổng số phần: %d\n\n", len(data))

	for part, content := range data {
		fmt.Printf("┌─────────────────────────────────────────────┐\n")
		fmt.Printf("│ PHẦN: %-39s │\n", part)
		fmt.Printf("├─────────────────────────────────────────────┤\n")
		fmt.Printf("│ Số lượng khóa: %-29d │\n", len(content))
		fmt.Printf("└─────────────────────────────────────────────┘\n")

		for key, values := range content {
			fmt.Printf("  • Khóa: \"%s\"\n", key)
			fmt.Printf("    Số lượng giá trị: %d\n", len(values))

			if len(values) == 0 {
				fmt.Println("    Giá trị: (không có)")
			} else {
				fmt.Println("    Giá trị:")
				for i, value := range values {
					fmt.Printf("      %d. %s\n", i+1, value)
				}
			}
			fmt.Println()
		}
		fmt.Println("------------------------------------------------")
	}

	// Thống kê tổng quát
	totalKeys := 0
	totalValues := 0
	for _, content := range data {
		totalKeys += len(content)
		for _, values := range content {
			totalValues += len(values)
		}
	}

	fmt.Printf("\n=== THỐNG KÊ TỔNG QUÁT ===\n")
	fmt.Printf("Tổng số phần: %d\n", len(data))
	fmt.Printf("Tổng số khóa: %d\n", totalKeys)
	fmt.Printf("Tổng số giá trị: %d\n", totalValues)
	fmt.Println("===========================\n")
}
