package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func open_file(path string) (*os.File, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("lỗi khi mở file %s: %v", path, err)
	}
	return file, nil
}

func read_text(file *os.File) map[string]bool {
	text_words := make(map[string]bool)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		text_words[strings.ToLower(strings.TrimSpace(scanner.Text()))] = true
	}
	return text_words
}

func create_lack(path string) (*os.File, error) {
	file, err := os.Create(path)
	if err != nil {
		return nil, fmt.Errorf("lỗi khi tạo file %s: %v", path, err)
	}
	return file, nil
}

func process_verb_line(line string, text_words map[string]bool, lack_word_file *os.File, skip_words map[string]bool) {
	verbs := strings.Split(strings.TrimSuffix(line, ":"), "/")
	for _, verb := range verbs {
		verb = strings.ToLower(strings.TrimSpace(verb))
		verb = strings.ReplaceAll(verb, "(uk)", "")
		verb = strings.ReplaceAll(verb, "(us)", "")
		if !skip_words[verb] && !text_words[verb] {
			fmt.Fprintln(lack_word_file, verb)
		}
	}
}

func scan_verbs(verbs_file *os.File, text_words map[string]bool, lack_word_file *os.File) {
	verb_scanner := bufio.NewScanner(verbs_file)
	var skip_words = map[string]bool{
		"ngữ cảnh": true,
		"ví dụ":    true,
	}

	for verb_scanner.Scan() {
		line := verb_scanner.Text()
		if strings.HasSuffix(line, ":") {
			process_verb_line(line, text_words, lack_word_file, skip_words)
		}
	}
}
