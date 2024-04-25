package main

import (
	"html/template"
	"log"
	"os"
)

type application struct {
	Logger   *log.Logger
	Template *template.Template
}

func NewApplication() *application {
	file, err := os.OpenFile("history.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal("Cannot open history log file:", err)
	}

	logger := log.New(file, "LOG: ", log.Ldate|log.Ltime|log.Lshortfile)
	tmpl := template.Must(template.ParseFiles("ui/templates/index.html"))

	return &application{
		Logger:   logger,
		Template: tmpl,
	}
}
