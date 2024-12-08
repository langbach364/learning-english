package main

import (
    "fmt"
    "log"
    "time"
)

func insertSampleData() error {
    db, err := connect_db()
    if err != nil {
        return err
    }
    defer db.Close()

    sampleData := []struct {
        date        string
        wordLearned int
        wrong       int
    }{
        {"2024-03-01", 10, 2},
        {"2024-03-01", 15, 3},
        {"2024-02-15", 30, 8},
        {"2024-02-16", 22, 5},
        {"2023-12-01", 45, 10},
        {"2023-11-01", 38, 7},
    }

    query := `INSERT INTO vocabulary_statistics (time, word_learned, wrong) VALUES (?, ?, ?)`
    
    for _, data := range sampleData {
        _, err := db.Exec(query, data.date, data.wordLearned, data.wrong)
        if err != nil {
            return err
        }
    }
    return nil
}

func printStats(timeRange TimeRange, stats []VocabStats) {
    title := map[TimeRange]string{
        Daily:   "Thống kê theo ngày",
        Monthly: "Thống kê theo tháng",
        Yearly:  "Thống kê theo năm",
    }

    fmt.Printf("\n%s:\n", title[timeRange])
    if len(stats) == 0 {
        fmt.Println("Không có dữ liệu")
        return
    }

    fmt.Println("┌──────────────┬──────────────┬───────┐")
    fmt.Println("│  Thời gian   │ Từ đã học    │  Sai  │")
    fmt.Println("├──────────────┼──────────────┼───────┤")

    timeFormat := map[TimeRange]string{
        Daily:   "02/01/2006",
        Monthly: "01/2006",
        Yearly:  "2006",
    }

    for _, stat := range stats {
        fmt.Printf("│ %s │ %12d │ %5d │\n",
            stat.Date.Format(timeFormat[timeRange]),
            stat.WordLearned,
            stat.Wrong)
    }
    fmt.Println("└──────────────┴──────────────┴───────┘")
}

func checkStatistics(year int, month int, day int) {
    date := time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.Local)
    searchDate := date.Format("02/01/2006")
    fmt.Printf("\n=== Thống kê cho ngày %s ===\n", searchDate)

    timeRanges := []TimeRange{Daily, Monthly, Yearly}
    for _, timeRange := range timeRanges {
        stats, err := get_data(timeRange, date)
        if err != nil {
            // Thêm ngày đang tìm kiếm vào thông báo lỗi
            log.Printf("Lỗi thống kê %v cho ngày %s: %v", timeRange, searchDate, err)
            continue
        }
        printStats(timeRange, stats)
    }
}


func main() {
    if err := insertSampleData(); err != nil {
        log.Fatalf("Lỗi khi chèn dữ liệu mẫu: %v", err)
    }
    fmt.Println("Đã chèn dữ liệu mẫu thành công!")

    var year, month, day int
    fmt.Print("\nNhập năm: ")
    fmt.Scan(&year)
    fmt.Print("Nhập tháng: ")
    fmt.Scan(&month)
    fmt.Print("Nhập ngày: ")
    fmt.Scan(&day)
    
    checkStatistics(year, month, day)
}
