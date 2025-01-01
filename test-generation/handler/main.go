package main

import (
	"time"
	"log"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func learn_word(e *echo.Echo, wordLearned int) {
	
	Send_Word(e, "/learn_word", wordLearned)
}

func create_schedule(e *echo.Echo) {
	Get_Word(e, "/create_schedule")
}

func revise_word(e *echo.Echo) {
	Get_Schedule(e, "/revise_word")
}

func get_statistics(e *echo.Echo) {
	Get_Statistic(e, "/get_statistics")
}

func setupLogging(e *echo.Echo) {
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "[${time_rfc3339}] ${method} ${uri} ${status} (${latency_human})\n",
	}))

	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			log.Printf("➡️ Nhận request mới: %s %s", c.Request().Method, c.Request().URL)
			
			start := time.Now()
			err := next(c)
			
			log.Printf("✅ Hoàn thành xử lý: %s %s - Status: %d - Thời gian xử lý: %v", 
				c.Request().Method, 
				c.Request().URL,
				c.Response().Status,
				time.Since(start))
				
			return err
		}
	})
}

func main() {
	wordLearned := 5
	reviewWord := 10
	
	log.Println("🚀 Khởi động server...")
	
	scheduling_word(reviewWord)
	enable_graphQL(":8080", "graph", wordLearned)

	rest := enable_rest("8081")
	setupLogging(rest)
	
	time.Sleep(1 * time.Second)

	learn_word(rest, wordLearned)
	create_schedule(rest)
	revise_word(rest)
	get_statistics(rest)
	
	log.Println("✨ Server đã sẵn sàng phục vụ")
	select {}
}
