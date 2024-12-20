package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/rs/cors"
)

type contextKey string

const (
	UserIDKey contextKey = "userID"
)

var (
	DEFAULT_USERNAME = load_API_key("USERNAME")
	DEFAULT_PASSWORD = load_API_key("PASSWORD")
)

func check_err(err error) {
	if err != nil {
		fmt.Println("Lỗi")
		log.Fatal(err)
	}
}

func send_error_response(w http.ResponseWriter, status int, message string, code string) {
	response := ErrorResponse{
		Status:  status,
		Message: message,
		Code:    code,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(response)
}

func login_handler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		send_error_response(w, http.StatusMethodNotAllowed, "Method không được hỗ trợ", "METHOD_NOT_ALLOWED")
		return
	}

	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		send_error_response(w, http.StatusBadRequest, "Dữ liệu không hợp lệ", "INVALID_REQUEST")
		return
	}

	if req.Username == DEFAULT_USERNAME && req.Password == DEFAULT_PASSWORD {
		token, err := generate_token(req.Username)
		if err != nil {
			send_error_response(w, http.StatusInternalServerError, "Lỗi tạo token", "TOKEN_GENERATION_ERROR")
			return
		}

		response := LoginResponse{Token: token}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	} else {
		send_error_response(w, http.StatusUnauthorized, "Sai thông tin đăng nhập", "INVALID_CREDENTIALS")
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

func auth_middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Request path: %s", r.URL.Path)

		if r.URL.Path == "/login" {
			next.ServeHTTP(w, r)
			return
		}

		if r.URL.Path == "/ChatCody" && strings.HasPrefix(r.Header.Get("Upgrade"), "websocket") {
			next.ServeHTTP(w, r)
			return
		}

		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			send_error_response(w, http.StatusUnauthorized, "Không tìm thấy token xác thực", "AUTH_MISSING_TOKEN")
			return
		}

		bearerToken := strings.Split(authHeader, " ")
		if len(bearerToken) != 2 || !strings.EqualFold(bearerToken[0], "Bearer") {
			send_error_response(w, http.StatusUnauthorized, "Token không đúng định dạng", "AUTH_INVALID_FORMAT")
			return
		}

		claims, err := validate_token(bearerToken[1])
		if err != nil {
			send_error_response(w, http.StatusUnauthorized, "Token không hợp lệ", "AUTH_INVALID_TOKEN")
			return
		}

		ctx := context.WithValue(r.Context(), UserIDKey, claims.UserID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func enable_middleware_cors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		Cors := cors.New(cors.Options{
			AllowedHeaders:   []string{"Authorization", "Content-Type"},
			AllowedMethods:   []string{"POST", "GET", "OPTIONS"},
			AllowedOrigins:   []string{"*"},
			AllowCredentials: true,
		})
		Cors.ServeHTTP(w, r, next.ServeHTTP)
	})
}

func path_file() map[string]string {
	return map[string]string{
		"word":        "../middleware/word.txt",
		"answer_cody": "./sourcegraph-cody/answer.txt",
	}
}

func muxtiplexer_router(router *http.ServeMux) {
	router.HandleFunc("/login", login_handler)
	data := path_file()
	router.HandleFunc("/word", write_word_file_api(data["word"]))
}

func muxtiplexer_websocket(router *http.ServeMux) {
	router.Handle("/ChatCody", websocket_cody("ChatCody"))
}

func create_server() {
	router := http.NewServeMux()
	muxtiplexer_router(router)
	muxtiplexer_websocket(router)

	authHandler := auth_middleware(router)
	corsHandler := enable_middleware_cors(authHandler)

	server := http.Server{
		Addr:    ":7089",
		Handler: corsHandler,
	}

	log.Fatal(server.ListenAndServe())
}
