package main

import (
	"encoding/json"
	"fmt"
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

	// In kết quả ra màn hình
	fmt.Fprintln(w, "Các từ ngẫu nhiên được tìm thấy:")
	fmt.Fprintln(w, "----------------------------------------")
	for _, word := range randomWords {
		fmt.Fprintf(w, "Từ: %s\n", word.Word)
		fmt.Fprintln(w, "----------------------------------------")
	}
	return randomWords
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

	return words
}
