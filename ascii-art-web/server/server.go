package server

import (
	"asciiartweb/handler"
	"log"
	"net/http"
)

func RegisterHandlers() {
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))
	http.HandleFunc("/", handler.ServeTemplate)
	http.HandleFunc("/ascii-art", handler.HandleAsciiArt)
}

func StartServer() {
	RegisterHandlers()

	log.Print("Server running on :http://localhost:8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}
}
