package handler

import (
	asciiart "asciiartweb/ascii-art"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

type PageData struct {
	Results       string
	InputText     string
	SelectedStyle string
}

// template renderer
func renderTemplate(w http.ResponseWriter, templateName string, data PageData) error {
	lp := filepath.Join("templates", "layout.html")
	fp := filepath.Join("templates", templateName)

	info, err := os.Stat(fp)
	if err != nil || info.IsDir() {
		return err
	}

	tmpl, err := template.ParseFiles(lp, fp)
	if err != nil {
		return err
	}

	return tmpl.ExecuteTemplate(w, "layout", data)
}

// error renderer
func renderError(w http.ResponseWriter, status int, templateName string) {
	w.WriteHeader(status)
	err := renderTemplate(w, templateName, PageData{})
	if err != nil {
		log.Println("Error rendering error page:", err)
	}
}

// Serve the main page
func ServeTemplate(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/index.html" {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	if r.URL.Path != "/" {
		renderError(w, http.StatusNotFound, "notFound.html")
		return
	}

	// Default page load
	data := PageData{
		Results:       "",
		InputText:     "",
		SelectedStyle: "standard",
	}

	err := renderTemplate(w, "index.html", data)
	if err != nil {
		log.Println("Error rendering index:", err)
		renderError(w, http.StatusInternalServerError, "serverError.html")
	}
}

// Handle form submission to generate ASCII art
func HandleAsciiArt(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		renderError(w, http.StatusInternalServerError, "serverError.html")
		return
	}

	userText := r.FormValue("text")
	banner := r.FormValue("fontstyle")

	allowedStyles := map[string]bool{
		"standard":   true,
		"shadow":     true,
		"thinkertoy": true,
	}

	if !allowedStyles[banner] {
		renderError(w, http.StatusNotFound, "notFound.html")
		return
	}

	bannerPath := filepath.Join("ascii-art", "banners", banner+".txt")
	asciiArtResult := asciiart.AsciiArt(userText, bannerPath)

	data := PageData{
		Results:       asciiArtResult,
		InputText:     userText,
		SelectedStyle: banner,
	}

	if err := renderTemplate(w, "index.html", data); err != nil {
		log.Println("Error rendering ASCII result:", err)
		renderError(w, http.StatusInternalServerError, "serverError.html")
	}
}
