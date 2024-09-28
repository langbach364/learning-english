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

func word_file(filePath string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "POST":
			{
				body, err := io.ReadAll(r.Body)
				check_err(err)

				var word Word
				err = json.Unmarshal(body, &word)
				check_err(err)

				wrFile, err := write_file(filePath)
				check_err(err)
				defer wrFile.Close()

				_, err = wrFile.WriteString(word.Data)
				fmt.Println("Từ đã được in trong file")
				check_err(err)
			}
		default:
			fmt.Println("Method không được sử dụng")
		}
	}
}

func listen_word_file(filePath string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "POST":
			{
				body, err := io.ReadAll(r.Body)
				check_err(err)

				var word Word
				err = json.Unmarshal(body, &word)
				check_err(err)

				wr, err := write_file(filePath)
				check_err(err)
				defer wr.Close()

				_, err = wr.WriteString(word.Data)
				fmt.Println("Từ đã được in trong file")
				check_err(err)
			}
		default:
			fmt.Println("Method không được sử dụng")
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

func path_file() map[string]string {
	return map[string]string{
		"word":       "../Middleware/word.txt",
		"listen_word": "../Middleware/listen.txt",
	}
}

func muxtiplexer_router(router *http.ServeMux) {
	data := path_file()
	router.HandleFunc("/word", word_file(data["word"]))
	router.HandleFunc("/listen_word", listen_word_file(data["listen_word"]))
}

func create_server() {
	router := http.NewServeMux()
	muxtiplexer_router(router)

	server := http.Server{
		Addr:    ":7089",
		Handler: enable_middleware_cors(router),
	}
	log.Fatal(server.ListenAndServe())
	server.ListenAndServe()
}
