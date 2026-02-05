package handler

import (
	asciiart "asciiartweb/ascii-art"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
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

	// Check if template file exists and is not a directory
	info, err := os.Stat(fp)
	if err != nil {
		return err
	}
	if info.IsDir() {
		return os.ErrNotExist
	}

	// Parse and execute template
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

// renderEmptyInputError renders the main page with an error message for empty input
func renderEmptyInputError(w http.ResponseWriter, banner string) {
	data := PageData{
		Results:       "",
		InputText:     "",
		SelectedStyle: banner,
		Error:         "Please enter some text to convert to ASCII art.",
	}
	if err := renderTemplate(w, "index.html", data); err != nil {
		log.Printf("Error rendering template: %v\n", err)
		renderError(w, http.StatusInternalServerError, "serverError.html")
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

// validateInput checks if the input contains only printable ASCII characters and newlines
func validateInput(text string) bool {
	for _, char := range text {
		// Allow newlines (\n, \r) and printable ASCII (space to tilde: 32-126)
		isNewline := char == '\n' || char == '\r'
		isPrintable := char >= 32 && char <= 126
		isBackslash := char == '\\'
		
		if !isNewline && !isPrintable && !isBackslash {
			return false
		}
	}
	return true
}

// interpretEscapeSequences converts escape sequences like \n, \t, etc. to their actual characters
func interpretEscapeSequences(text string) string {
	// Replace common escape sequences
	text = strings.ReplaceAll(text, "\\n", "\n")
	text = strings.ReplaceAll(text, "\\t", "\t")
	text = strings.ReplaceAll(text, "\\r", "\r")
	return text
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

	userTextRaw := r.FormValue("text")
	banner := r.FormValue("fontstyle")

	// Validate that banner is selected (set default if empty)
	if banner == "" {
		banner = "standard"
	}

	// Check if input is empty or only whitespace
	if strings.TrimSpace(userTextRaw) == "" {
		renderEmptyInputError(w, banner)
		return
	}

	// Interpret escape sequences like \n, \t, etc. for rendering only
	userText := interpretEscapeSequences(userTextRaw)

	// Validate banner selection
	if !allowedBanners[banner] {
		log.Printf("Invalid banner selected: %s\n", banner)
		renderError(w, http.StatusBadRequest, "notFound.html")
		return
	}

	// Validate input text (only printable ASCII + newlines)
	if !validateInput(userTextRaw) {
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
	
	// Check if generation failed (empty result when input was not empty)
	if asciiArtResult == "" && strings.TrimSpace(userText) != "" {
		log.Println("ASCII art generation failed")
		renderError(w, http.StatusInternalServerError, "serverError.html")
		return
	}

	data := PageData{
		Results:       template.HTML(asciiArtResult),
		InputText:     userTextRaw,
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
