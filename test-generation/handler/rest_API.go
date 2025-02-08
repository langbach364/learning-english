package main

import (
	"database/sql"
	"log"
	"net/http"
	"time"

	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/echo/v4"
)

var (
	clientWords []infoWord
)

func connect_db() (*sql.DB, error) {
	db, err := sql.Open("mysql", "root:@ztegc4df9f4e@tcp(localhost:3306)/learned_vocabulary")
	if err != nil {
		return nil, err
	}
	return db, nil
}

func Send_Word(e *echo.Echo, pattern string, wordLearned int) error {
	e.POST(pattern, func(c echo.Context) error {
		data := generate_word(wordLearned)
		if data == nil {
			return c.JSON(http.StatusOK, map[string]string{
				"message": "No data available",
			})
		}
		return c.JSON(http.StatusOK, data)
	})
	
	e.OPTIONS(pattern, func(c echo.Context) error {
		return c.NoContent(http.StatusOK)
	})
	
	return nil
}
func Get_Schedule(e *echo.Echo, pattern string) error {
	e.POST(pattern, func(c echo.Context) error {
		var target TargetDate
		if err := c.Bind(&target); err != nil {
			log.Printf("Lỗi binding: %v", err)
			return c.JSON(http.StatusBadRequest, map[string]string{
				"message": "Lỗi khi nhận dữ liệu",
			})
		}

		parsedDate, err := time.Parse("2006-01-02", target.Date)
		if err != nil {
			log.Printf("Lỗi parse date: %v", err)
			return c.JSON(http.StatusBadRequest, map[string]string{
				"message": "Định dạng ngày không hợp lệ",
			})
		}

		words, err := get_schedule(parsedDate)
		if err != nil {
			log.Printf("Lỗi get_schedule: %v", err)
			return c.JSON(http.StatusInternalServerError, map[string]string{
				"message": "Lỗi khi lấy lịch học",
			})
		}

		return c.JSON(http.StatusOK, map[string]interface{}{
			"words": words,
		})
	})
	return nil
}

func Get_Word(e *echo.Echo, pattern string) error {
	e.POST(pattern, func(c echo.Context) error {
		var receivedWords []infoWord

		if err := c.Bind(&receivedWords); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{
				"message": "Lỗi khi nhận dữ liệu",
			})
		}

		clientWords = append(clientWords, receivedWords...)
		for _, word := range receivedWords {
			wordChannel <- word
			log.Println("Đã gửi từ vào channel:", word.Word)
		}

		return c.JSON(http.StatusOK, map[string]string{
			"message": "Đã gửi dữ liệu thành công",
		})
	})
	return nil
}

func Get_Statistic(e *echo.Echo, pattern string) error {
	e.POST(pattern, func(c echo.Context) error {
		var t Timege
		if err := c.Bind(&t); err != nil {
			log.Printf("Lỗi binding JSON: %v", err)
			return c.JSON(http.StatusBadRequest, map[string]string{
				"message": "Lỗi khi nhận dữ liệu",
			})
		}

		parsedDate, err := time.Parse("2006-01-02", t.Date)
		if err != nil {
			log.Printf("Lỗi parse date: %v", err)
			return c.JSON(http.StatusBadRequest, map[string]string{
				"message": "Định dạng ngày không hợp lệ",
			})
		}

		data, err := get_vocabulary_stastics(TimeRange(t.Range), parsedDate)
		if err != nil {
			log.Printf("Lỗi get_vocabulary_stastics: %v", err)
			return c.JSON(http.StatusInternalServerError, map[string]string{
				"message": "Lỗi khi lấy thống kê",
			})
		}

		return c.JSON(http.StatusOK, data)
	})
	return nil
}
func enable_rest(port string) *echo.Echo {
	e := echo.New()

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete},
		AllowHeaders:     []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
		AllowCredentials: true,
	}))

	serverAddr := ":" + port
	println("REST Server is running at", serverAddr)

	go e.Start(serverAddr)
	return e
}