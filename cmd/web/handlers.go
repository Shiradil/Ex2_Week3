package main

import (
	"encoding/json"
	"github.com/go-resty/resty/v2"
	"net/http"
	"strings"
)

func (app *application) handleRequest(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		app.Template.Execute(w, nil)
		return
	}

	r.ParseForm()
	question := r.FormValue("question")
	client := resty.New()

	// Filtration
	var thematicFilter = []string{"tourism", "travel", "vacation"}

	containsTheme := false
	for _, theme := range thematicFilter {
		if strings.Contains(strings.ToLower(question), theme) {
			containsTheme = true
			break
		}
	}

	if !containsTheme {
		app.Template.Execute(w, map[string]interface{}{"Response": "Your request was declined because your question is not related to the vision of the touristic company"})
		return
	}

	app.Logger.Printf("Received question: %s", question)

	// API call
	response, err := client.R().
		SetAuthToken(apiKey).
		SetHeader("Content-Type", "application/json").
		SetBody(map[string]interface{}{
			"model": "gpt-3.5-turbo",
			"messages": []interface{}{
				map[string]interface{}{"role": "user", "content": question},
			},
			"max_tokens": 1024,
		}).
		Post(apiEndpoint)

	if err != nil {
		app.Logger.Printf("Error while sending the request: %v", err)
		http.Error(w, "Failed to send request to OpenAI", http.StatusInternalServerError)
		return
	}

	// Log the response body for debugging
	app.Logger.Printf("Response Body: %s", response.Body())

	var data map[string]interface{}
	err = json.Unmarshal(response.Body(), &data)
	if err != nil {
		app.Logger.Printf("Error while decoding JSON response: %v", err)
		http.Error(w, "Failed to decode JSON response", http.StatusInternalServerError)
		return
	}

	choices, ok := data["choices"].([]interface{})
	if !ok || len(choices) == 0 {
		app.Logger.Println("Invalid or no choices in response")
		http.Error(w, "No response choices found", http.StatusInternalServerError)
		return
	}

	choice, ok := choices[0].(map[string]interface{})
	if !ok {
		app.Logger.Println("Invalid format for choice")
		http.Error(w, "Invalid response format", http.StatusInternalServerError)
		return
	}

	message, ok := choice["message"].(map[string]interface{})
	if !ok {
		app.Logger.Println("Invalid format for message")
		http.Error(w, "Invalid response message format", http.StatusInternalServerError)
		return
	}

	content, ok := message["content"].(string)
	if !ok {
		app.Logger.Println("Content is not a string")
		http.Error(w, "Content format error", http.StatusInternalServerError)
		return
	}

	app.Template.Execute(w, map[string]interface{}{"Response": content})
}
