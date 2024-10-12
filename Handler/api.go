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

func write_word_file_api(filePath string) http.HandlerFunc {
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

				w.WriteHeader(http.StatusOK)
				w.Write([]byte("Success"))
			}
		default:
			fmt.Println("Method không được sử dụng")
		}
	}
}

func read_file_word_api() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "POST":
			{
				data := data_structure()
				jsonData, err := json.MarshalIndent(data, "", "    ")
				check_err(err)

				fmt.Println(string(jsonData))
				response := map[string]string{
					"data": string(jsonData),
				}

				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).SetIndent("", "    ")
				json.NewEncoder(w).Encode(response)
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
		"word":        "../Middleware/word.txt",
		"listen_word": "../Middleware/listen.txt",
		"read_word":   "./sourcegraph-cody/answer.txt",
	}
}

func muxtiplexer_router(router *http.ServeMux) {
	data := path_file()
	router.HandleFunc("/word", write_word_file_api(data["word"]))
	router.HandleFunc("/listen_word", write_word_file_api(data["listen_word"]))
	router.HandleFunc("/read_word", read_file_word_api())
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
