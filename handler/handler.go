package handler

import (
	asciiart "asciiartweb/ascii-art"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	// "strings"
)

type PageData struct {
	Results       template.HTML
	InputText     string
	SelectedStyle string
	Error         string
}

var allowedBanners = map[string]bool{
	"standard":   true,
	"shadow":     true,
	"thinkertoy": true,
}

// renderTemplate renders an HTML template with the given data
func renderTemplate(w http.ResponseWriter, templateName string, data PageData) error {
	lp := filepath.Join("templates", "layout.html")
	fp := filepath.Join("templates", templateName)

	// Check if template file exists
	info, err := os.Stat(fp)
	if err != nil {
		if os.IsNotExist(err) {
			return err
		}
		return err
	}
	
	if info.IsDir() {
		return os.ErrNotExist
	}

	tmpl, err := template.ParseFiles(lp, fp)
	if err != nil {
		return err
	}

	return tmpl.ExecuteTemplate(w, "layout", data)
}

// renderError renders an error page with the appropriate status code
func renderError(w http.ResponseWriter, status int, templateName string) {
	w.WriteHeader(status)
	err := renderTemplate(w, templateName, PageData{})
	if err != nil {
		log.Printf("Error rendering %s: %v\n", templateName, err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

// ServeTemplate handles GET / - serves the main page
func ServeTemplate(w http.ResponseWriter, r *http.Request) {
	// Redirect /index.html to /
	if r.URL.Path == "/index.html" {
		http.Redirect(w, r, "/", http.StatusMovedPermanently)
		return
	}

	// Only accept root path
	if r.URL.Path != "/" {
		renderError(w, http.StatusNotFound, "notFound.html")
		return
	}

	// Only accept GET method
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	data := PageData{
		Results:       "",
		InputText:     "",
		SelectedStyle: "standard",
		Error:         "",
	}

	err := renderTemplate(w, "index.html", data)
	if err != nil {
		log.Printf("Error rendering index: %v\n", err)
		renderError(w, http.StatusInternalServerError, "serverError.html")
	}
}

// validateInput checks if the input contains only printable ASCII characters
func validateInput(text string) bool {
	for _, char := range text {
		// Allow newlines and printable ASCII (32-126)
		if char != '\n' && char != '\r' && (char < 32 || char > 126) {
			return false
		}
	}
	return true
}

// HandleAsciiArt handles POST /ascii-art - generates ASCII art
func HandleAsciiArt(w http.ResponseWriter, r *http.Request) {
	// Only accept POST method
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	// Parse form data
	if err := r.ParseForm(); err != nil {
		log.Printf("Error parsing form: %v\n", err)
		renderError(w, http.StatusBadRequest, "notFound.html")
		return
	}

	userText := r.FormValue("text")
	banner := r.FormValue("fontstyle")

	// Validate banner selection
	if !allowedBanners[banner] {
		log.Printf("Invalid banner selected: %s\n", banner)
		renderError(w, http.StatusBadRequest, "notFound.html")
		return
	}

	// Validate input text (only printable ASCII + newlines)
	if !validateInput(userText) {
		log.Println("Input contains invalid characters")
		renderError(w, http.StatusBadRequest, "notFound.html")
		return
	}

	// Check if banner file exists
	bannerPath := filepath.Join("ascii-art", "banners", banner+".txt")
	if _, err := os.Stat(bannerPath); os.IsNotExist(err) {
		log.Printf("Banner file not found: %s\n", bannerPath)
		renderError(w, http.StatusNotFound, "notFound.html")
		return
	}

	// Generate ASCII art
	asciiArtResult := asciiart.AsciiArt(userText, bannerPath)
	
	// Check if generation failed
	if asciiArtResult == "" && userText != "" {
		log.Println("ASCII art generation failed")
		renderError(w, http.StatusInternalServerError, "serverError.html")
		return
	}

	data := PageData{
		Results:       template.HTML(asciiArtResult),
		InputText:     userText,
		SelectedStyle: banner,
		Error:         "",
	}

	// Render result (200 OK by default)
	err := renderTemplate(w, "index.html", data)
	if err != nil {
		log.Printf("Error rendering ASCII result: %v\n", err)
		renderError(w, http.StatusInternalServerError, "serverError.html")
	}
}
