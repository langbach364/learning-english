package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

var wordChannel = make(chan infoWord)

func rules(word infoWord, currentTime time.Time) time.Time {
	x := word.Frequency - word.ErrorCount
	if x > 0 {
		return currentTime.AddDate(0, 0, x)
	}
	return currentTime.AddDate(0, 0, 1)
}

func execute_graphQL_mutation(mutation string) error {
	payload := map[string]string{
		"query": mutation,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("lỗi marshal JSON: %v", err)
	}

	resp, err := http.Post("http://localhost:8080/graphql",
		"application/json",
		bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("lỗi gọi API: %v", err)
	}
	defer resp.Body.Close()
	return nil
}

func rules_scheduling(words []infoWord) {
	currentTime := time.Now()

	for _, word := range words {
		nextTime := rules(word, currentTime)

		mutation := fmt.Sprintf(`mutation {
            insert {
                schedule(
                    time: "%s",
                    word: "%s"
                )
            }
        }`, nextTime.Format("2006-01-02"), word.Word)

		if err := execute_graphQL_mutation(mutation); err != nil {
			log.Printf("Lỗi tạo lịch cho từ %s: %v", word.Word, err)
		}
	}
}

func process_word(word infoWord) {
	vocabularyMutation := fmt.Sprintf(`mutation {
        insert {
            vocabulary(
                word: "%s",
                frequency: %d,
                error_count: %d
            )
        }
    }`, word.Word, word.Frequency, word.ErrorCount)

	if err := execute_graphQL_mutation(vocabularyMutation); err != nil {
		log.Printf("Lỗi thêm từ vào vocabulary: %v", err)
		return
	}

	rules_scheduling([]infoWord{word})
}

func scheduling_word() {
	go func() {
		for word := range wordChannel {
			process_word(word)
		}
	}()
}

func get_schedule() ([]string, error) {
    db, err := connect_db()
    if err != nil {
        return nil, fmt.Errorf("không thể kết nối đến cơ sở dữ liệu: %v", err)
    }
    defer db.Close()

    currentDate := time.Now().Format("2006-01-02")

    query := fmt.Sprintf(`
        SELECT word 
        FROM schedule 
        WHERE time = '%s'`, currentDate)

    rows, err := db.Query(query)
    if err != nil {
        return nil, fmt.Errorf("không thể lấy danh sách từ: %v", err)
    }
    defer rows.Close()

    var words []string
    for rows.Next() {
        var word string
        if err := rows.Scan(&word); err != nil {
            return nil, fmt.Errorf("không thể đọc từ: %v", err)
        }
        words = append(words, word)
    }

    if err := rows.Err(); err != nil {
        return nil, fmt.Errorf("lỗi trong quá trình đọc dữ liệu: %v", err)
    }

    return words, nil
}