package main

import (
	"log"
	"time"

	"github.com/labstack/echo/v4"
)

func learn_word(e *echo.Echo, wordLearned int) {
	data := generate_word(wordLearned)
	Send_Word(e, "/learn_word", data)
}

func create_schedule(e *echo.Echo) {
	Get_Word(e, "/create_schedule")
}

func revise_word(e *echo.Echo) {
	data, err := get_schedule()
	if err != nil {
		log.Printf("Lỗi lấy lịch: %v", err)
		return
	}
	Send_Word(e, "/revise_word", data)
}

func main() {
	wordLearned := 5
	reviewWord := 10
	scheduling_word(reviewWord)
	enable_graphQL(":8080", "graph", wordLearned)

	rest := enable_rest("8081")
	time.Sleep(1 * time.Second)
	learn_word(rest, wordLearned)

	
	create_schedule(rest)
	revise_word(rest)
	select {}
}
