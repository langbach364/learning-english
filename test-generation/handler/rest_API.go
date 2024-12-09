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

func Send_Word(e *echo.Echo, pattern string, data interface{}) error {
	e.POST(pattern, func(c echo.Context) error {
		return c.JSON(http.StatusOK, data)
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
		}

		return c.JSON(http.StatusOK, map[string]string{
			"message": "Đã gửi dữ liệu thành công",
		})
	})
	return nil
}

func enable_rest(port string) *echo.Echo {
	e := echo.New()

	serverAddr := ":" + port
	println("REST Server is running at", serverAddr)

	go e.Start(serverAddr)
	return e
}