package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"
	"unicode"
	"unicode/utf8"

	"github.com/joho/godotenv"
)

func get_api() string {
	err := godotenv.Load(".env")
	if err != nil {
		fmt.Println("lỗi không load được file .env")
	}
	return os.Getenv("API_KEY")
}

func check_err_request(err error) {
	if err != nil {
		fmt.Println("lỗi không thể gửi yêu cầu")
		fmt.Printf("cụ thể: %s\n", err)
	}
}

func get_attributes(attributes *map[string]bool) map[string]bool {
	*attributes = map[string]bool{
		"Define":       true,
		"PartOfSpeech": true,
	}
	return *attributes
}

func convert_struct_to_map(word interface{}) map[string]interface{} {
	result := make(map[string]interface{})

	v := reflect.ValueOf(word)
	if v.Kind() == reflect.Struct {
		for i := 0; i < v.NumField(); i++ {
			field := v.Type().Field(i)
			value := v.Field(i).Interface()
			result[field.Name] = value
		}
	}

	return result
}

func clean_define(define string) string {
	var result string
	in_tag := false
	for _, char := range define {
		if char == '<' {
			in_tag = true
		} else if char == '>' {
			in_tag = false
		} else if !in_tag {
			result += string(char)
		}
	}
	return result
}

func limit_definition(words map[string][]string, partOfSpeech string, definition string, limit int) {
	if len(words[partOfSpeech]) < limit {
		words[partOfSpeech] = append(words[partOfSpeech], definition)
	}
}

func classify_word(words []map[string]interface{}) map[string][]string {
	result := make(map[string][]string)
	for _, word := range words {
		partOfSpeech, okPos := word["PartOfSpeech"].(string)
		definition, okDef := word["Define"].(string)
		if okPos && okDef {
			limit_definition(result, partOfSpeech, definition, 5)
		}
	}
	return result
}

func write_file_translated(words map[string][]string) error {
	file, err := os.Create("./translate/trans.txt")
	if err != nil {
		return fmt.Errorf("lỗi khi tạo file trans.txt: %v", err)
	}
	defer file.Close()

	for pos, word := range words {
		for _, wd := range word {
			_, err := file.WriteString(fmt.Sprintf("%s: %s\n", pos, wd))
			if err != nil {
				return fmt.Errorf("lỗi khi ghi vào file trans.txt: %v", err)
			}
		}
	}
	return nil
}

func read_file_translate() (map[string][]string, error) {
	translatedFile, err := os.Open("./translate/trans_ed.txt")
	if err != nil {
		return nil, fmt.Errorf("lỗi khi mở file trans_ed.txt: %v", err)
	}
	defer translatedFile.Close()

	scanner := bufio.NewScanner(translatedFile)
	translatedData := make(map[string][]string)

	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.SplitN(line, ": ", 2)
		if len(parts) == 2 {
			pos, def := parts[0], parts[1]
			translatedData[pos] = append(translatedData[pos], def)
		}
	}
	return translatedData, nil
}

func update_words_map(words []map[string]interface{}, translatedData map[string][]string) map[string][]string {
	for i, word := range words {
		pos, ok := word["PartOfSpeech"].(string)
		if ok {
			if translatedDefs, exists := translatedData[pos]; exists && len(translatedDefs) > 0 {
				words[i]["Define"] = translatedDefs[0]
				translatedData[pos] = translatedDefs[1:]
			}
		}
	}
	return translatedData
}

func extract_replace(input string, extractedWords map[int]string, startCounter int) (string, int) {
    counter := startCounter

    if extractedWords == nil {
        extractedWords = make(map[int]string)
    }

    replaceFunc := func(s string) string {
        word := s[1 : len(s)-1]
        extractedWords[counter] = word
        result := fmt.Sprintf("(%d)", counter)
        counter++
        return result
    }

    result := regexp.MustCompile(`\(.*?\)`).ReplaceAllStringFunc(input, replaceFunc)
    return result, counter
}

func handle_map_vi(words *[]map[string]interface{}, extractedWords *map[int]string) map[string][]string {
	classified := classify_word(*words)
    if *extractedWords == nil {
        *extractedWords = make(map[int]string)
    }
    counter := 0
    
	for pos, definitions := range classified {
		for i, def := range definitions {
			if pos == "idiom" {
				classified[pos][i], counter = extract_replace(clean_define(def), *extractedWords, counter)
			} else {
				classified[pos][i] = clean_define(def)
			}
		}
	}
	
	err := write_file_translated(classified)

	if err != nil {
		fmt.Println(err)
		return classified
	}

	time.Sleep(3 * time.Second)

	translatedData, err := read_file_translate()
	if err != nil {
		fmt.Println(err)
		return classified
	}

	return update_words_map(*words, translatedData)
}

func handle_map_en(words *[]map[string]interface{}) map[string][]string {
	classified := classify_word(*words)
	for pos, definitions := range classified {
		for i, def := range definitions {
			classified[pos][i] = clean_define(def)
		}
	}
	return classified
}

func check_value_map(definitions []string) bool {
	for _, def := range definitions {
		if def != "" {
			return true
		}
	}
	return false
}

func format_output(pos string, definitions []string) string {
	pos = strings.TrimSuffix(pos, ".")
	if len(pos) > 0 {
		r, size := utf8.DecodeRuneInString(pos)
		pos = string(unicode.ToUpper(r)) + pos[size:]
	}
	result := fmt.Sprintf("%s:\n", pos)
	for _, def := range definitions {
		if def != "" {
			result += fmt.Sprintf("- %s\n", def)
		}
	}
	return strings.TrimSpace(result)
}

func replace_numbers(input string, extractedWords map[int]string) string {
    replaceFunc := func(s string) string {
        numStr := s[1 : len(s)-1]
        if num, err := strconv.Atoi(numStr); err == nil {
            if word, ok := extractedWords[num]; ok {
                return "(" + word + ")"
            }
        }
        return s
    }
    
    return regexp.MustCompile(`\(\d+\)`).ReplaceAllStringFunc(input, replaceFunc)
}

func output_word(words map[string][]string, extractedWords map[int]string) {
    isFirst := true
    for pos, definitions := range words {
        if len(definitions) > 0 && check_value_map(definitions) {
            if !isFirst {
                fmt.Println("-------------------------------------")
            }
            replacedDefinitions := make([]string, len(definitions))
            for i, def := range definitions {
                replacedDefinitions[i] = replace_numbers(def, extractedWords)
            }
            formatted := format_output(pos, replacedDefinitions)
            fmt.Println(formatted)
            isFirst = false
        }
    }
}


func fetch_word(url string) ([]map[string]interface{}, error) {
	response, err := http.Get(url)
	check_err_request(err)
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	var word_definitions []word_definitions
	err = json.Unmarshal(body, &word_definitions)
	if err != nil {
		return nil, err
	}

	definitions := make([]map[string]interface{}, len(word_definitions))
	for attribute, text := range word_definitions {
		definitions[attribute] = convert_struct_to_map(text)
	}
	return definitions, nil
}

func define_word(word string) {
	api_key := get_api()
	url := fmt.Sprintf("https://api.wordnik.com/v4/word.json/%s/definitions?limit=300&includeRelated=false&sourceDictionaries=ahd-5&useCanonical=false&includeTags=false&api_key=%s", word, api_key)

	definitions, err := fetch_word(url)
	check_err_request(err)

	attributes := make(map[string]bool)
	attributes = get_attributes(&attributes)

	var extractedWords map[int]string
	result := handle_map_vi(&definitions, &extractedWords)

	fmt.Printf("Từ: %s\n\n", word)
	output_word(result, extractedWords)
}

func main() {
	define_word("make")
}
