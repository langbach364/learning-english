package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

func find_next_question_index(questionDir string) int {
	var nextIndex int = 1

	files, err := os.ReadDir(questionDir)
	if err == nil {
		for _, file := range files {
			fileName := file.Name()
			if strings.HasPrefix(fileName, "question(") && strings.HasSuffix(fileName, ").txt") {
				numStr := strings.TrimPrefix(fileName, "question(")
				numStr = strings.TrimSuffix(numStr, ").txt")
				if num, err := strconv.Atoi(numStr); err == nil && num >= nextIndex {
					nextIndex = num + 1
				}
			}
		}
	}
	return nextIndex
}

func save_question_to_file(questionDir string, index int, question string) error {
	questionFile := fmt.Sprintf("%s/question(%d).txt", questionDir, index)
	return os.WriteFile(questionFile, []byte(question), 0644)
}

func read_answer_from_file(answerDir string, index int) (string, error) {
	answerFile := fmt.Sprintf("%s/answer(%d).txt", answerDir, index)
	answer, err := os.ReadFile(answerFile)
	if err != nil {
		return "", err
	}
	return string(answer), nil
}

func check_answer_file_exists(answerDir string) bool {
	_, err := os.Stat(answerDir + "/answer.txt")
	return err == nil
}

func chat_with_cody(question string) string {
	questionDir := "./sourcegraph-cody/chat-cody/question"
	answerDir := "./sourcegraph-cody/chat-cody/answer"

	if !check_answer_file_exists(answerDir) {
		fmt.Println("Khởi tạo file answer.txt...")
		run_script("./sourcegraph-cody/chat-cody/cody.sh")
	}

	nextIndex := find_next_question_index(questionDir)

	err := save_question_to_file(questionDir, nextIndex, question)
	if err != nil {
		fmt.Println("Lỗi khi lưu câu hỏi:", err)
		return ""
	}

	run_script("./sourcegraph-cody/chat-cody/cody.sh")

	answer, err := read_answer_from_file(answerDir, nextIndex)
	if err != nil {
		fmt.Println("Lỗi khi đọc câu trả lời:", err)
		return ""
	}

	return answer
}
