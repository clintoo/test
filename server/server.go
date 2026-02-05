package server

import (
	"asciiartweb/handler"
	"fmt"
	"log"
	"net/http"
)

const serverPort = ":8080"

// RegisterHandlers sets up all HTTP routes and static file serving
func RegisterHandlers() {
	// Serve static files (CSS, images, etc.)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))
	
	// Main page route
	http.HandleFunc("/", handler.ServeTemplate)
	
	// ASCII art generation route
	http.HandleFunc("/ascii-art", handler.HandleAsciiArt)
}

// StartServer initializes and starts the HTTP server
func StartServer() {
	RegisterHandlers()

	fmt.Printf("Server running on http://localhost%s\n", serverPort)
	
	err := http.ListenAndServe(serverPort, nil)
	if err != nil {
		log.Fatalf("Server failed to start: %v\n", err)
	}
}
