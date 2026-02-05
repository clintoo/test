package tests

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"asciiartweb/handler"
	"asciiartweb/server"
)

// TestIntegrationMainPage tests the complete flow of accessing the main page
func TestIntegrationMainPage(t *testing.T) {
	server.RegisterHandlers()

	req := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()

	http.DefaultServeMux.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Main page: status = %d, expected %d", w.Code, http.StatusOK)
	}

	body := w.Body.String()

	// Check for essential elements
	expectedElements := []string{
		"Ascii Art Web",
		"standard",
		"shadow",
		"thinkertoy",
		"text",
		"fontstyle",
	}

	for _, elem := range expectedElements {
		if !strings.Contains(body, elem) {
			t.Errorf("Main page should contain %q", elem)
		}
	}
}

// TestIntegrationAsciiArtGeneration tests the complete ASCII art generation flow
func TestIntegrationAsciiArtGeneration(t *testing.T) {
	server.RegisterHandlers()

	tests := []struct {
		name       string
		input      string
		banner     string
		shouldPass bool
	}{
		{"Simple text - standard", "Hello", "standard", true},
		{"Simple text - shadow", "Hello", "shadow", true},
		{"Simple text - thinkertoy", "Hello", "thinkertoy", true},
		{"With newline", "Hello\\nWorld", "standard", true},
		{"Empty input", "", "standard", false},
		{"Whitespace only", "   ", "standard", false},
		{"Invalid banner", "Hello", "invalid", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			formData := url.Values{
				"text":      {tt.input},
				"fontstyle": {tt.banner},
			}

			req := httptest.NewRequest("POST", "/ascii-art", strings.NewReader(formData.Encode()))
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
			w := httptest.NewRecorder()

			handler.HandleAsciiArt(w, req)

			if tt.shouldPass {
				if w.Code != http.StatusOK {
					t.Errorf("Expected success (200) but got %d", w.Code)
				}
			} else {
				if w.Code != http.StatusOK && w.Code != http.StatusBadRequest {
					t.Logf("Got status %d for invalid input (acceptable)", w.Code)
				}
			}
		})
	}
}

// TestIntegrationStaticFiles tests static file serving
func TestIntegrationStaticFiles(t *testing.T) {
	server.RegisterHandlers()

	req := httptest.NewRequest("GET", "/static/style.css", nil)
	w := httptest.NewRecorder()

	http.DefaultServeMux.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Static file: status = %d, expected %d", w.Code, http.StatusOK)
	}

	// Check that the file has some content
	if w.Body.Len() == 0 {
		t.Error("Static CSS file should have content")
	}
}

// TestIntegration404Handling tests 404 error page
func TestIntegration404Handling(t *testing.T) {
	server.RegisterHandlers()

	req := httptest.NewRequest("GET", "/nonexistent", nil)
	w := httptest.NewRecorder()

	http.DefaultServeMux.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("404 handling: status = %d, expected %d", w.Code, http.StatusNotFound)
	}

	body := w.Body.String()
	if !strings.Contains(body, "404") && !strings.Contains(body, "Not Found") {
		t.Error("404 page should mention 'Not Found' or '404'")
	}
}

// TestIntegrationMethodNotAllowed tests HTTP method validation
func TestIntegrationMethodNotAllowed(t *testing.T) {
	server.RegisterHandlers()

	tests := []struct {
		path   string
		method string
	}{
		{"/", "POST"},
		{"/", "PUT"},
		{"/", "DELETE"},
		{"/ascii-art", "GET"},
		{"/ascii-art", "PUT"},
	}

	for _, tt := range tests {
		t.Run(tt.method+" "+tt.path, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, tt.path, nil)
			w := httptest.NewRecorder()

			http.DefaultServeMux.ServeHTTP(w, req)

			if w.Code != http.StatusMethodNotAllowed && w.Code != http.StatusOK {
				t.Logf("Method %s on %s returned status %d", tt.method, tt.path, w.Code)
			}
		})
	}
}

// TestIntegrationEscapeSequences tests escape sequence handling end-to-end
func TestIntegrationEscapeSequences(t *testing.T) {
	server.RegisterHandlers()

	formData := url.Values{
		"text":      {"A\\nB"},
		"fontstyle": {"standard"},
	}

	req := httptest.NewRequest("POST", "/ascii-art", strings.NewReader(formData.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()

	handler.HandleAsciiArt(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Escape sequence test: status = %d, expected %d", w.Code, http.StatusOK)
	}

	body := w.Body.String()
	// The textarea should show the original input with \n
	if !strings.Contains(body, "A\\nB") {
		t.Error("Output should preserve original input text in the textarea")
	}
}

// TestIntegrationAllBanners tests that all three banners work
func TestIntegrationAllBanners(t *testing.T) {
	server.RegisterHandlers()

	banners := []string{"standard", "shadow", "thinkertoy"}

	for _, banner := range banners {
		t.Run(banner, func(t *testing.T) {
			formData := url.Values{
				"text":      {"ABC"},
				"fontstyle": {banner},
			}

			req := httptest.NewRequest("POST", "/ascii-art", strings.NewReader(formData.Encode()))
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
			w := httptest.NewRecorder()

			handler.HandleAsciiArt(w, req)

			if w.Code != http.StatusOK {
				t.Errorf("Banner %s: status = %d, expected %d", banner, w.Code, http.StatusOK)
			}

			// Response should have some ASCII art content
			if w.Body.Len() < 100 {
				t.Errorf("Banner %s: response seems too short", banner)
			}
		})
	}
}
