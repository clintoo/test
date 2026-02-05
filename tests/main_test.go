package tests

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"asciiartweb/server"
)

func TestIntegrationMainPage(t *testing.T) {
	// fix working directory (tests run inside /tests)
	os.Chdir("..")

	// reset mux to avoid "pattern conflicts" panic
	http.DefaultServeMux = http.NewServeMux()

	server.RegisterHandlers()

	req := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()

	http.DefaultServeMux.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200 got %d", w.Code)
	}
}
