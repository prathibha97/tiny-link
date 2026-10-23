package shortener

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestCreateURL(t *testing.T) {
	store := NewURLStore()
	handler := NewHandler(store, "http://localhost:8080")

	body := `{"url": "https://go.dev"}`
	req := httptest.NewRequest(http.MethodPost, "/api/urls", strings.NewReader(body))

	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	handler.CreateURL(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf(
			"expected status %d, got %d: %s",
			http.StatusCreated,
			rec.Code,
			rec.Body.String(),
		)
	}

	var response createURLResponse
	if err := json.NewDecoder(rec.Body).Decode(&response); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if response.Code == "" {
		t.Fatal("expected a generated short code")
	}

	if response.OriginalURL != "https://go.dev" {
		t.Errorf("unexpected original URL: %s", response.OriginalURL)
	}

	if response.ShortURL != "http://localhost:8080/"+response.Code {
		t.Errorf("unexpected short URL: %s", response.ShortURL)
	}
}

func TestRedirect(t *testing.T) {
	store := NewURLStore()
	handler := NewHandler(store, "http://localhost:8080")

	store.Save(context.Background(), URLMapping{
		Code:        "go123",
		OriginalURL: "https://go.dev",
	})

	req := httptest.NewRequest(http.MethodGet, "/go123", nil)
	rec := httptest.NewRecorder()
	handler.Redirect(rec, req)

	if rec.Code != http.StatusFound {
		t.Fatalf(
			"expected status %d, got %d: %s",
			http.StatusFound,
			rec.Code,
			rec.Body.String(),
		)
	}

	if rec.Header().Get("Location") != "https://go.dev" {
		t.Errorf("unexpected redirect location: %s", rec.Header().Get("Location"))
	}
}

func TestRedirectNotFound(t *testing.T) {
	store := NewURLStore()
	handler := NewHandler(store, "http://localhost:8080")

	req := httptest.NewRequest(http.MethodGet, "/missing123", nil)
	rec := httptest.NewRecorder()
	handler.Redirect(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf(
			"expected status %d, got %d: %s",
			http.StatusNotFound,
			rec.Code,
			rec.Body.String(),
		)
	}

	if rec.Header().Get("Location") != "" {
		t.Errorf("unexpected redirect location: %s", rec.Header().Get("Location"))
	}
}

func TestCreateURLInvalidScheme(t *testing.T) {
	store := NewURLStore()
	handler := NewHandler(store, "http://localhost:8080")

	body := `{"url": "ftp://example.com/file"}`
	req := httptest.NewRequest(http.MethodPost, "/api/urls", strings.NewReader(body))

	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	handler.CreateURL(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf(
			"expected status %d, got %d: %s",
			http.StatusBadRequest,
			rec.Code,
			rec.Body.String(),
		)
	}

	if len(store.urls) != 0 {
		t.Errorf("expected empty store, got %d mappings", len(store.urls))
	}

}

func TestCreateURLInvalidJSON(t *testing.T) {
	store := NewURLStore()
	handler := NewHandler(store, "http://localhost:8080")

	body := `{"url": }`
	req := httptest.NewRequest(http.MethodPost, "/api/urls", strings.NewReader(body))

	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	handler.CreateURL(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf(
			"expected status %d, got %d: %s",
			http.StatusBadRequest,
			rec.Code,
			rec.Body.String(),
		)
	}

	if len(store.urls) != 0 {
		t.Errorf("expected empty store, got %d mappings", len(store.urls))
	}

}
