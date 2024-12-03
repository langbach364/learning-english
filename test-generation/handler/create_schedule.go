package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func get_data(query string) {
	url := "http://localhost:8080/graphql"

	requestBody, _ := json.Marshal(map[string]string{
        "query": query,
    })

	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer resp.Body.Close()

	reponseBody, _ := io.ReadAll(resp.Body)
	fmt.Println(string(reponseBody))
}

func schedule(limitQuery int, ) {
	query := fmt.Sprintf(`{
        vocabulary(limit: %d) {
            word
            frequency
            error_count
        }
    }`, limitQuery)

	get_data(query)


}