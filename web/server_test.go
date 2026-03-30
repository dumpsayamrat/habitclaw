package web

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestGetRoot(t *testing.T) {
	handler := NewHandler()
	srv := httptest.NewServer(handler)
	defer srv.Close()

	resp, err := http.Get(srv.URL + "/")
	if err != nil {
		t.Fatalf("GET / failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}

	ct := resp.Header.Get("Content-Type")
	if !strings.Contains(ct, "text/html") {
		t.Errorf("expected Content-Type to contain text/html, got %q", ct)
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("reading body: %v", err)
	}
	body := string(bodyBytes)

	if !strings.Contains(body, "HabitClaw") {
		t.Errorf("expected body to contain 'HabitClaw', got %q", body)
	}
}

func TestGetStaticCSS(t *testing.T) {
	handler := NewHandler()
	srv := httptest.NewServer(handler)
	defer srv.Close()

	resp, err := http.Get(srv.URL + "/static/style.css")
	if err != nil {
		t.Fatalf("GET /static/style.css failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}

	ct := resp.Header.Get("Content-Type")
	if !strings.Contains(ct, "text/css") {
		t.Errorf("expected Content-Type to contain text/css, got %q", ct)
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("reading body: %v", err)
	}
	if len(bodyBytes) == 0 {
		t.Error("expected non-empty CSS body")
	}
}

func TestGetAPIHealth(t *testing.T) {
	// Wire up the same way main.go does: health endpoint on a top-level mux,
	// web handler mounted at /.
	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"status":  "ok",
			"version": "0.1.0",
		})
	})
	mux.Handle("/", NewHandler())

	srv := httptest.NewServer(mux)
	defer srv.Close()

	resp, err := http.Get(srv.URL + "/api/health")
	if err != nil {
		t.Fatalf("GET /api/health failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}

	ct := resp.Header.Get("Content-Type")
	if !strings.Contains(ct, "application/json") {
		t.Errorf("expected Content-Type to contain application/json, got %q", ct)
	}

	var result map[string]string
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		t.Fatalf("failed to decode JSON: %v", err)
	}

	if result["status"] != "ok" {
		t.Errorf("expected status 'ok', got %q", result["status"])
	}
	if result["version"] != "0.1.0" {
		t.Errorf("expected version '0.1.0', got %q", result["version"])
	}
}

func TestEmbeddedFiles(t *testing.T) {
	tests := []struct {
		name string
		path string
	}{
		{"index.html", "static/index.html"},
		{"style.css", "static/style.css"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := StaticFiles.ReadFile(tt.path)
			if err != nil {
				t.Fatalf("ReadFile(%q) error: %v", tt.path, err)
			}
			if len(data) == 0 {
				t.Errorf("ReadFile(%q) returned empty content", tt.path)
			}
		})
	}
}
