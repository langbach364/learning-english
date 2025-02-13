package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
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

func find_available_day(db *sql.DB, startTime time.Time, limitQuery int) time.Time {
	currentTime := startTime
	for {
		countQuery := fmt.Sprintf(`
            SELECT COUNT(*) 
            FROM schedule 
            WHERE time = '%s'`,
			currentTime.Format("2006-01-02"))

		var wordCount int
		row := db.QueryRow(countQuery)
		if err := row.Scan(&wordCount); err != nil {
			log.Printf("Lỗi đếm số từ: %v", err)
			return currentTime
		}

		if wordCount < limitQuery {
			return currentTime
		}
		currentTime = currentTime.AddDate(0, 0, 1)
	}
}

func execute_graphQL_mutation(mutation string) error {
	payload := map[string]string{
		"query": mutation,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Lỗi marshal JSON: %v", err)
		return err
	}

	resp, err := http.Post("http://localhost:8081/graph",
		"application/json",
		bytes.NewBuffer(jsonData))
	if err != nil {
		log.Printf("Lỗi gọi API: %v", err)
		return err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	log.Printf("GraphQL response: %s", string(body))

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("GraphQL error: %s", string(body))
	}

	return nil
}

func process_word(word infoWord, limitQuery int) {
	vocabularyMutation := fmt.Sprintf(`
	mutation {
		insert {
			vocabulary(word: "%s", frequency: %d, error_count: %d)
		}
	}`, word.Word, word.Frequency, word.ErrorCount)

	if err := execute_graphQL_mutation(vocabularyMutation); err != nil {
		log.Printf("Lỗi thêm từ vào vocabulary: %v", err)
		return
	}
	log.Printf("Đã thêm từ vào database: %s", word.Word)

	db, err := connect_db()
	if err != nil {
		log.Printf("Lỗi kết nối database: %v", err)
		return
	}
	defer db.Close()

	nextTime := find_available_day(db, rules(word, time.Now()), limitQuery)

	var maxPriority int
	err = db.QueryRow(`
        SELECT COALESCE(MAX(priority), 0) 
        FROM schedule 
        WHERE time = ?`,
		nextTime.Format("2006-01-02")).Scan(&maxPriority)
	if err != nil {
		log.Printf("Lỗi lấy max priority: %v", err)
		return
	}

	var maxID int64
	err = db.QueryRow("SELECT COALESCE(MAX(id), 0) FROM schedule").Scan(&maxID)
	if err != nil {
		log.Printf("Lỗi lấy max id: %v", err)
		return
	}
	newID := maxID + 1

	scheduleMutation := fmt.Sprintf(`
	mutation {
		insert {
			schedule(id: %d, time: "%s", word: "%s", priority: %d)
		}
	}`, newID, nextTime.Format("2006-01-02"), word.Word, maxPriority+1)

	if err := execute_graphQL_mutation(scheduleMutation); err != nil {
		log.Printf("Lỗi thêm lịch học: %v", err)
		return
	}
	log.Printf("Đã sắp xếp từ %s vào trong database với bản schedule", word.Word)
}

func scheduling_word(limitQuery int) {
	go func() {
		for word := range wordChannel {
			process_word(word, limitQuery)
		}
	}()
}

func get_schedule(targetDate time.Time) ([]string, error) {
	formattedTime := targetDate.Format("2006-01-02")
	db, err := connect_db()
	if err != nil {
		return nil, fmt.Errorf("không thể kết nối đến cơ sở dữ liệu: %v", err)
	}
	defer db.Close()

	query := fmt.Sprintf(`
        SELECT word 
        FROM schedule 
        WHERE time = '%s'
        ORDER BY priority ASC`, formattedTime)

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

	return words, nil
}
