package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/rs/cors"
)

var TOKEN string

func init() {
	TOKEN = load_API_key("TOKEN")
}

func clean_token(token string) string {
	token = strings.Trim(token, "\"")
	return token
}

func validate_token(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := clean_token(r.Header.Get("Authorization"))

		if token == "" {
			http.Error(w, "Token không được để trống", http.StatusUnauthorized)
			return
		}

		if token != TOKEN {
			http.Error(w, "Token không hợp lệ", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	}
}

func check_err(err error) {
	if err != nil {
		fmt.Println("Lỗi")
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

func websocket_cody(nameEvent string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		conn := connect_client(w, r)

		if conn == nil {
			return
		}
		handle_connection(nameEvent, conn)
	}
}

func enable_middleware_cors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		Cors := cors.New(cors.Options{
			AllowedHeaders:   []string{"Accept", "Accept-Language", "Content-Language", "Content-Type"},
			AllowedMethods:   []string{"POST", "GET", "OPTIONS", "PUT", "DELETE"},
			AllowedOrigins:   []string{"*"},
			AllowCredentials: true,
			Debug:            true,
		})
		Cors.ServeHTTP(w, r, next.ServeHTTP)
	})
}

func path_file() map[string]string {

	return map[string]string{
		"word": "../middleware/word.txt",
		//	"listen_word": "../middleware/listen.txt",
		"answer_cody": "./sourcegraph-cody/answer.txt",
	}
}

func muxtiplexer_router(router *http.ServeMux) {
	data := path_file()
	router.HandleFunc("/word", validate_token(write_word_file_api(data["word"])))
	// router.HandleFunc("/listen_word", write_word_file_api(data["listen_word"]))
}

func muxtiplexer_websocket(router *http.ServeMux) {
	router.HandleFunc("/ChatCody", validate_token(websocket_cody("ChatCody")))
}

func create_server() {
	router := http.NewServeMux()

	muxtiplexer_router(router)
	muxtiplexer_websocket(router)

	server := http.Server{
		Addr:    ":7089",
		Handler: enable_middleware_cors(router),
	}

	log.Fatal(server.ListenAndServe())
}
