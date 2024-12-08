package main

import (
	"database/sql"
	"net/http"

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


func enable_rest(port, pattern string) {
	e := echo.New()

	e.POST(pattern, func(c echo.Context) error {
		word := new(infoWord)
		if err := c.Bind(word); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{
				"error": "Invalid data format",
			})
		}
		clientWords = append(clientWords, *word)
		wordChannel <- *word
		return c.JSON(http.StatusOK, word)
	})

	serverAddr := ":" + port
	println("REST Server is running at", serverAddr)

	e.Start(serverAddr)
}
