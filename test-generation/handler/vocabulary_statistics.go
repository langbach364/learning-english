package main

import (
	"log"
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

func get_vocabulary_stastics(timeRange TimeRange, date time.Time) ([]VocabStats, error) {
    formattedDate := date.Format("2006-01-02")
    log.Printf("Bắt đầu lấy thống kê với timeRange=%s, date=%s", timeRange, formattedDate)

    db, err := connect_db()
    if err != nil {
        log.Printf("Lỗi kết nối DB: %v", err)
        return nil, err
    }
    defer db.Close()

    query := queries[timeRange]
    args := prepare_args(timeRange, date)

    log.Printf("Thực thi query: %s với args: %v", query, args)

    rows, err := db.Query(query, args...)
    if err != nil {
        log.Printf("Lỗi query DB: %v", err)
        return nil, err
    }
    defer rows.Close()

    var stats []VocabStats
    for rows.Next() {
        var stat VocabStats
        var dateStr string

        if err := rows.Scan(&dateStr, &stat.WordLearned, &stat.Wrong); err != nil {
            log.Printf("Lỗi scan row: %v", err)
            return nil, err
        }

        stat.Date, err = time.Parse("2006-01-02", dateStr)
        if err != nil {
            log.Printf("Lỗi parse date %s: %v", dateStr, err)
            return nil, err
        }

        stats = append(stats, stat)
        log.Printf("Đã thêm stat: Date=%s, WordLearned=%d, Wrong=%d", 
            stat.Date.Format("2006-01-02"), stat.WordLearned, stat.Wrong)
    }

    log.Printf("Hoàn thành, trả về %d records", len(stats))
    return stats, nil
}
