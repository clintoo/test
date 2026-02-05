package server

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestRegisterHandlers(t *testing.T) {
	// Reset DefaultServeMux so tests don't conflict
	http.DefaultServeMux = http.NewServeMux()

	RegisterHandlers()

	tests := []string{"/", "/ascii-art", "/static/"}

	for _, path := range tests {
		h, pattern := http.DefaultServeMux.Handler(&http.Request{Method: "GET", URL: mustParseURL(path)})
		if h == nil {
			t.Fatalf("expected handler for %s, got nil", path)
		}
		if pattern == "" {
			t.Fatalf("expected pattern for %s, got empty", path)
		}
	}
}

func mustParseURL(path string) *url.URL {
	req := httptest.NewRequest("GET", "http://localhost"+path, nil)
	return req.URL
}
