package main

import (
	"log"
	"net/http"
)

const apiKey = "sk-proj-UchxBLJyVrfedWH7JJFlT3BlbkFJ2yasZndADMp5uxSS3ji0"
const apiEndpoint = "https://api.openai.com/v1/chat/completions"

func main() {
	app := NewApplication()

	http.HandleFunc("/", app.handleRequest)
	app.Logger.Println("Server is running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
