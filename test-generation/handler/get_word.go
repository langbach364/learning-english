package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

var APIKey string

type RandomWord struct {
	Id   int    `json:"id"`
	Word string `json:"word"`
}

func load_API_key(nameAPI string) string {
	if err := godotenv.Load("../../enviroment/.env"); err != nil {
		fmt.Println("Lỗi: Không thể tải file .env")
	}
	return os.Getenv(nameAPI)
}

func init() {
	APIKey = load_API_key("API_wordnik")
}

func random_word(w http.ResponseWriter, field FieldsWord) []RandomWord {
	baseURL := "https://api.wordnik.com/v4/words.json/randomWords"

	params := url.Values{}
	params.Add("api_key", APIKey)

	wordParams := RandomWordParams{
		HasDictionaryDef:    field.HasDictionaryDef,
		MinLength:           field.MinLength,
		MaxLength:           field.MaxLength,
		Limit:               field.Limit,
		IncludePartOfSpeech: field.IncludePartOfSpeech,
	}

	if wordParams.HasDictionaryDef {
		params.Add("hasDictionaryDef", "true")
	}
	if wordParams.IncludePartOfSpeech != "" {
		params.Add("includePartOfSpeech", wordParams.IncludePartOfSpeech)
	}
	if wordParams.MinLength > 0 {
		params.Add("minLength", fmt.Sprintf("%d", wordParams.MinLength))
	}
	if wordParams.MaxLength > 0 {
		params.Add("maxLength", fmt.Sprintf("%d", wordParams.MaxLength))
	}
	if wordParams.Limit > 0 {
		params.Add("limit", fmt.Sprintf("%d", wordParams.Limit))
	}

	fullURL := baseURL + "?" + params.Encode()

	resp, err := http.Get(fullURL)
	if err != nil {
		fmt.Fprintf(w, "Lỗi khi gọi API: %v\n", err)
		return nil
	}
	defer resp.Body.Close()

	var randomWords []RandomWord
	if err := json.NewDecoder(resp.Body).Decode(&randomWords); err != nil {
		fmt.Fprintf(w, "Lỗi khi đọc dữ liệu: %v\n", err)
		return nil
	}

	fmt.Fprintln(w, "Các từ ngẫu nhiên được tìm thấy:")
	fmt.Fprintln(w, "----------------------------------------")
	for _, word := range randomWords {
		fmt.Fprintf(w, "Từ: %s\n", word.Word)
		fmt.Fprintln(w, "----------------------------------------")
	}
	return randomWords
}

func check_duplicate_words(words []RandomWord) []RandomWord {
	db, err := sql.Open("mysql", "root:@ztegc4df9f4e@tcp(localhost:3306)/learned_vocabulary")
	if err != nil {
		log.Printf("Lỗi kết nối database: %v", err)
		return words
	}
	defer db.Close()

	for i, word := range words {
		for {
			var exists bool
			err := db.QueryRow("SELECT EXISTS(SELECT 1 FROM vocabulary WHERE word = ?)", word.Word).Scan(&exists)
			if err != nil {
				log.Printf("Lỗi truy vấn database: %v", err)
				break
			}

			if !exists {
				break
			}

			w := httptest.NewRecorder()
			field := FieldsWord{
				HasDictionaryDef: true,
				MinLength:        3,
				MaxLength:        10,
				Limit:            1,
			}
			newWords := random_word(w, field)
			if len(newWords) > 0 {
				words[i] = newWords[0]
				word = newWords[0]
			}
		}
	}

	return words
}

func generate_word(limitWord int) []RandomWord {
	w := httptest.NewRecorder()
	listInclue := []string{
		"noun",
		"adjective",
		"verb",
		"adverb",
		"interjection",
		"pronoun",
		"preposition",
		"verb-transitive",
		"verb-intransitive",
		"past-participle",
		"noun-posessive",
		"imperative",
		"noun-plural",
		"definite-article",
		"conjunction",
		"auxiliary-verb",
		"article",
	}

	field := FieldsWord{
		HasDictionaryDef:    true,
		IncludePartOfSpeech: strings.Join(listInclue, ","),
		MinLength:           3,
		MaxLength:           10,
		Limit:               limitWord,
	}

	words := random_word(w, field)

	words = check_duplicate_words(words)
	
	return words
}
