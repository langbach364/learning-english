package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"regexp"
	"strings"
	"strconv"
	"sync"
	"unicode"
	"unicode/utf8"

	"github.com/joho/godotenv"
)


// Hàm chạy đầu tiên trước cả hàm main()
/////////////////////////////////////////////////////////////////////////////////////
// Hàm khởi tạo đầu tiên và lấy API key từ file .env
/////////////////////////////////////////////////////////////////////////////////////
// Hàm phụ trợ 
/////////////////////////////////////////////////////////////////////////////////////
var (
	APIKey    string
	regexPool sync.Pool
)

func init() {
	APIKey = load_API_key("API_wordnik")
	regexPool = sync.Pool{
		New: func() interface{} {
			return regexp.MustCompile(`\(.*?\)`)
		},
	}
}

func load_API_key(nameAPI string) string {
	if err := godotenv.Load("./enviroment/.env"); err != nil {
		fmt.Println("Lỗi: Không thể tải file .env")
	}
	return os.Getenv(nameAPI)
}


// Kiểm tra các đói tượng trong slice có chứa item hay không
func contains(slice []string, item string) bool { 
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

// Ghi hoa chữ đầu và thêm "." vào cuối chuỗi (đã xử lý trường hợp ký tự khác bảng mã ascii)
func normalize_key(key string) string {
	key = strings.TrimSuffix(key, ".")
	if len(key) > 0 {
		r, size := utf8.DecodeRuneInString(key)
		key = string(unicode.ToUpper(r)) + key[size:]
	}
	return key
}

// Giá trị có cùng tên key sau khi đã đưa về key trùng
func group_attributes(attributes *map[string][]string) map[string][]string {
	groupedAttributes := make(map[string][]string)
	for key, values := range *attributes {
		normalizedKey := normalize_key(key)
		groupedAttributes[normalizedKey] = append(groupedAttributes[normalizedKey], values...)
	}
	return groupedAttributes
}

// Thay thế các con số bằng các từ ngữ đã lưu trữ trước đó trong biến extractedWords
func replace_numbers(input string, extractedWords map[int]string) string {
	regex := regexPool.Get().(*regexp.Regexp)
	defer regexPool.Put(regex)

	return regex.ReplaceAllStringFunc(input, func(s string) string {
		numStr := s[1 : len(s)-1]
		if num, err := strconv.Atoi(numStr); err == nil {
			if word, ok := extractedWords[num]; ok {
				return "(" + word + ")"
			}
		}
		return s
	})
}

// Liểm tra giá trị có phải rỗng hay không
func check_value_map(definitions []string) bool {
	for _, def := range definitions {
		if def != "" {
			return true
		}
	}
	return false
}

// Làm sạch sẽ lại dữ liệu rồi in ra 
func format_output(pos string, definitions []string) string {
	var result strings.Builder
	result.WriteString(fmt.Sprintf("%s:\n", pos))
	for _, def := range definitions {
		if def != "" {
			result.WriteString(fmt.Sprintf("- %s\n", def))
		}
	}
	return strings.TrimSpace(result.String())
}


/////////////////////////////////////////////////////////////////////////////////////


// Các bước để chạy chương trình
/////////////////////////////////////////////////////////////////////////////////////
// Lấy các định nghĩa của từ API Wordnik
func fetch_word_definitions(word string) ([]WordDefinition, error) {
	url := fmt.Sprintf("https://api.wordnik.com/v4/word.json/%s/definitions?limit=300&includeRelated=false&sourceDictionaries=ahd-5&useCanonical=false&includeTags=false&api_key=%s", word, APIKey)

	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("lỗi khi gửi yêu cầu: %v", err)
	}
	defer resp.Body.Close()

	var definitions []WordDefinition
	decoder := json.NewDecoder(resp.Body)

	if err := decoder.Decode(&definitions); err != nil {
		return nil, fmt.Errorf("lỗi khi giải mã JSON: %v", err)
	}

	return definitions, nil
}

// Phân loại loại từ theo định nghĩa của từ giới hạn mỗi từ loại là 5 định nghĩa 
func classify_definitions(definitions []WordDefinition) map[string][]string {
	result := make(map[string][]string)
	for _, def := range definitions {
		if len(result[def.PartOfSpeech]) < 5 {
			result[def.PartOfSpeech] = append(result[def.PartOfSpeech], def.Define)
		}
	}
	return result
}

// Làm sạch lại dữ liệu khi đã lấy được dữ liệu từ api
func clean_definition(define string) string {
	var result strings.Builder
	inTag := false
	for _, char := range define {
		if char == '<' {
			inTag = true
		} else if char == '>' {
			inTag = false
		} else if !inTag {
			result.WriteRune(char)
		}
	}
	return result.String()
}

// Cập nhật lại các định nghĩa đã được dịch
func update_definitions(definitions []WordDefinition, translatedData map[string][]string) map[string][]string {
	for i := range definitions {
		pos := definitions[i].PartOfSpeech
		if translatedDefs, exists := translatedData[pos]; exists && len(translatedDefs) > 0 {
			definitions[i].Define = translatedDefs[0]
			translatedData[pos] = translatedDefs[1:]
		}
	}
	return translatedData
}

// Qúa trình trích xuất từ không cần dịch dựa trên từ loại và cấu trúc trích xuất dựa trên return regexp.MustCompile(`\(.*?\)`) ở hàm init()
func process_untranslated(classified *map[string][]string, extractedWords *map[int]string) {
	counter := 0
	class := []string{"idiom", "intransitive verb"}

	for pos, definitions := range *classified {
		for i, def := range definitions {
			if contains(class, pos) {
				(*classified)[pos][i], counter = extract_and_replace(clean_definition(def), *extractedWords, counter)
			} else {
				(*classified)[pos][i] = clean_definition(def)
			}
		}
	}
}

// Trích xuất các từ đã lưu trước đó và thay thế chúng bằng các số thứ tự tương ứng (để kiểm soát các từ không cần dịch)
func extract_and_replace(input string, extractedWords map[int]string, startCounter int) (string, int) {
	regex := regexPool.Get().(*regexp.Regexp)
	defer regexPool.Put(regex)

	counter := startCounter
	return regex.ReplaceAllStringFunc(input, func(s string) string {
		word := s[1 : len(s)-1]
		extractedWords[counter] = word
		result := fmt.Sprintf("(%d)", counter)
		counter++
		return result
	}), counter
}
/////////////////////////////////////////////////////////////////////////////////////

// Hàm ghép các đoạn xử lý với bản muốn dịch
func handle_vietnamese_map(definitions *[]WordDefinition, extractedWords *map[int]string) map[string][]string {
	classified := classify_definitions(*definitions)

	if *extractedWords == nil {
		*extractedWords = make(map[int]string)
	}

	process_untranslated(&classified, extractedWords)

	if err := write_translation_file(classified); err != nil {
		fmt.Println(err)
		return classified
	}

	if err := wait_tool_complete(socketPath); err != nil {
        fmt.Println(err)
        return classified
    }

	translatedData, err := read_translated_file()
	if err != nil {
		fmt.Println(err)
		return classified
	}

	translatedData = group_attributes(&translatedData)
	return update_definitions(*definitions, translatedData)
}

// Hàm ghép lại các đoạn xử lý với bản không muốn dịch
func handle_english_map(definitions *[]WordDefinition) map[string][]string {
	classified := classify_definitions(*definitions)
	for pos, defs := range classified {
		for i, def := range defs {
			classified[pos][i] = clean_definition(def)
		}
	}
	return group_attributes(&classified)
}
/////////////////////////////////////////////////////////////////////////////////////

/////////////////////////////////////////////////////////////////////////////////////
// In ra txt cho công cụ dịch xử lý
func write_translation_file(words map[string][]string) error {
	file, err := write_file("./translate/trans.txt")
	if err != nil {
		return fmt.Errorf("lỗi khi tạo file: %v", err)
	}
	defer file.Close()

	for pos, definitions := range words {
		for _, def := range definitions {
			if _, err := fmt.Fprintf(file, "%s: %s\n", pos, def); err != nil {
				return fmt.Errorf("lỗi khi ghi vào file trans.txt: %v", err)
			}
		}
	}
	return nil
}

// Đọc file txt đã được công cụ xử lý
func read_translated_file() (map[string][]string, error) {
	file, err := read_file("./translate/trans_ed.txt")
	if err != nil {
		return nil, fmt.Errorf("lỗi khi đọc file: %v", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	translatedData := make(map[string][]string)

	for scanner.Scan() {
		parts := strings.SplitN(scanner.Text(), ": ", 2)
		if len(parts) == 2 {
			pos, def := parts[0], parts[1]
			translatedData[pos] = append(translatedData[pos], def)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("lỗi khi đọc file: %v", err)
	}

	return translatedData, nil
}

// Xuất ra các từ đã cấu trúc rồi
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

// Tổng hợp lại 
func define_word(word string) {
	definitions, err := fetch_word_definitions(word)
	if err != nil {
		fmt.Printf("Lỗi khi lấy định nghĩa: %v\n", err)
		return
	}

	var extractedWords map[int]string
	result := handle_vietnamese_map(&definitions, &extractedWords)

	fmt.Printf("Từ: %s\n\n", word)
	output_word(result, extractedWords)
}
/////////////////////////////////////////////////////////////////////////////////////