package main

import (
	"time"
)

type VocabStats struct {
	Date        time.Time
	WordLearned int
	Wrong       int
}

type TimeRange string

const (
	Daily   TimeRange = "daily"
	Monthly TimeRange = "monthly"
	Yearly  TimeRange = "yearly"
)

var queries = map[TimeRange]string{
	Daily: `SELECT 
		DATE(time) as date,
		SUM(word_learned) as total_learned,
		SUM(wrong) as total_wrong
	FROM vocabulary_statistics 
	WHERE DATE(time) = ?
	GROUP BY DATE(time)`,

	Monthly: `SELECT 
		DATE_FORMAT(MIN(time), '%Y-%m-%d') as date,
		SUM(word_learned) as total_learned,
		SUM(wrong) as total_wrong
	FROM vocabulary_statistics 
	WHERE YEAR(time) = ? AND MONTH(time) = ?
	GROUP BY YEAR(time), MONTH(time)`,

	Yearly: `SELECT 
		DATE_FORMAT(MIN(time), '%Y-%m-%d') as date,
		SUM(word_learned) as total_learned,
		SUM(wrong) as total_wrong
	FROM vocabulary_statistics 
	WHERE YEAR(time) = ?
	GROUP BY YEAR(time)`,
}

func prepare_args(timeRange TimeRange, date time.Time) []interface{} {
	switch timeRange {
	case Daily:
		return []interface{}{date.Format("2006-01-02")}
	case Monthly:
		return []interface{}{date.Year(), int(date.Month())}
	case Yearly:
		return []interface{}{date.Year()}
	default:
		return nil
	}
}

func get_data(timeRange TimeRange, date time.Time) ([]VocabStats, error) {
	db, err := connect_db()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	query := queries[timeRange]
	args := prepare_args(timeRange, date)

	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var stats []VocabStats
	for rows.Next() {
		var stat VocabStats
		var dateStr string

		if err := rows.Scan(&dateStr, &stat.WordLearned, &stat.Wrong); err != nil {
			return nil, err
		}

		if stat.Date, err = time.Parse("2006-01-02", dateStr); err != nil {
			return nil, err
		}

		stats = append(stats, stat)
	}

	return stats, nil
}
