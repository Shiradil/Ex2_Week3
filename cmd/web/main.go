package main

import (
	"encoding/json"
	"fmt"
	"github.com/go-resty/resty/v2"
	"html/template"
	"log"
	"net/http"
)

const apiKey = ""
const apiEndpoint = "https://api.openai.com/v1/chat/completions"

var tmpl = template.Must(template.ParseFiles("ui/templates/index.html"))

func main() {
	http.HandleFunc("/", handleRequest)
	fmt.Println("Server is running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func handleRequest(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		tmpl.Execute(w, nil)
		return
	}

	r.ParseForm()
	question := r.FormValue("question")
	client := resty.New()
	response, err := client.R().
		SetAuthToken(apiKey).
		SetHeader("Content-Type", "application/json").
		SetBody(map[string]interface{}{
			"model": "gpt-3.5-turbo",
			"messages": []interface{}{
				map[string]interface{}{"role": "user", "content": question},
			},
			"max_tokens": 512,
		}).
		Post(apiEndpoint)

	if err != nil {
		log.Printf("Error while sending the request: %v", err)
		http.Error(w, "Failed to send request to OpenAI", http.StatusInternalServerError)
		return
	}

	// Log the response body for debugging
	log.Printf("Response Body: %s", response.Body())

	var data map[string]interface{}
	err = json.Unmarshal(response.Body(), &data)
	if err != nil {
		log.Printf("Error while decoding JSON response: %v", err)
		http.Error(w, "Failed to decode JSON response", http.StatusInternalServerError)
		return
	}

	choices, ok := data["choices"].([]interface{})
	if !ok || len(choices) == 0 {
		log.Println("Invalid or no choices in response")
		http.Error(w, "No response choices found", http.StatusInternalServerError)
		return
	}

	choice, ok := choices[0].(map[string]interface{})
	if !ok {
		log.Println("Invalid format for choice")
		http.Error(w, "Invalid response format", http.StatusInternalServerError)
		return
	}

	message, ok := choice["message"].(map[string]interface{})
	if !ok {
		log.Println("Invalid format for message")
		http.Error(w, "Invalid response message format", http.StatusInternalServerError)
		return
	}

	content, ok := message["content"].(string)
	if !ok {
		log.Println("Content is not a string")
		http.Error(w, "Content format error", http.StatusInternalServerError)
		return
	}

	tmpl.Execute(w, map[string]interface{}{"Response": content})
}
