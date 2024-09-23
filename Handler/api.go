package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/rs/cors"
)

func check_err(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func word_file(filePath string) http.HandlerFunc{
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "POST":
			{
				body, err := io.ReadAll(r.Body)
				check_err(err)

				var word Word
				err = json.Unmarshal(body, &word)
				check_err(err)

				go middleware_Word(filePath)
				wrFile, err := write_file(filePath)
				check_err(err)
				defer wrFile.Close()

				_, err = wrFile.WriteString(word.Data)
				fmt.Println("Từ đã được in trong file")
				check_err(err)
			}
		default: fmt.Println("Method không được sử dụng")
		}
	}
}

func enable_middleware_cors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		Cors := cors.New(cors.Options{
			AllowedHeaders:   []string{"Accept", "Accept-Language", "Content-Language", "Content-Type"},
			AllowedMethods:   []string{"POST"},
			AllowedOrigins:   []string{"*"},
			AllowCredentials: true,
			Debug:            true,
		})
		Cors.ServeHTTP(w, r, next.ServeHTTP)
	})
}

func muxtiplexer_router(router *http.ServeMux) {
	filePath := "../Middleware/word.txt"
	router.HandleFunc("/word", word_file(filePath))
}

func create_server() {
	router := http.NewServeMux()
	muxtiplexer_router(router)

	server := http.Server{
		Addr:    ":5050",
		Handler: enable_middleware_cors(router),
	}
	log.Fatal(server.ListenAndServe())
	server.ListenAndServe()
}