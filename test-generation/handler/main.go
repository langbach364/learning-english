package main

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"log"
	"time"
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
	log.Println("⚙️ Đang cấu hình logging...")

	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "📝 [${time_rfc3339}] ${method} ${uri} ${status} (${latency_human})\n",
	}))

	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			log.Printf("➡️ Nhận request mới: %s %s", c.Request().Method, c.Request().URL)

			start := time.Now()
			err := next(c)

			log.Printf("✅ Hoàn thành: %s %s - Status: %d - Thời gian: %v",
				c.Request().Method,
				c.Request().URL,
				c.Response().Status,
				time.Since(start))

			return err
		}
	})
	log.Println("✅ Cấu hình logging hoàn tất")
}

func main() {
	// wordLearned := 5
	// reviewWord := 10

	// log.Println("🎯 Khởi động ứng dụng với cấu hình:")
	// log.Printf("📚 Số từ học mới: %d", wordLearned)
	// log.Printf("🔄 Số từ ôn tập: %d", reviewWord)

	// log.Println("🚀 Khởi động server...")

	// scheduling_word(reviewWord)
	// log.Println("📅 Đã khởi tạo lịch học")

	// rest := enable_rest("8081")
	// setupLogging(rest)

	// enable_graphQL(":8081", "graph", wordLearned)
	// log.Println("🎯 GraphQL server đã sẵn sàng")

	// time.Sleep(1 * time.Second)

	// learn_word(rest, wordLearned)
	// log.Println("📖 Đã cấu hình học từ mới")

	// create_schedule(rest)
	// log.Println("📅 Đã tạo lịch học")

	// revise_word(rest)
	// log.Println("🔄 Đã cấu hình ôn tập")

	// get_statistics(rest)
	// log.Println("📊 Đã cấu hình thống kê")

	// log.Println("✨ Server đã sẵn sàng phục vụ")
	// select {}

	data := handler_data()
	print_data(data)
}
