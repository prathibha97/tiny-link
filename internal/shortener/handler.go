package shortener

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
	"strings"
)

type Handler struct {
	store    Store
	base_url string
}

type createURLRequest struct {
	URL string `json:"url"`
}

type createURLResponse struct {
	Code        string `json:"code"`
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}

func NewHandler(store Store, base_url string) *Handler {
	return &Handler{
		store:    store,
		base_url: strings.TrimRight(base_url, "/"),
	}
}

func (h *Handler) CreateURL(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)

	var req createURLRequest
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{
			"error": "invalid JSON request",
		})
		return
	}

	// Reject trailing JSON values or other non-whitespace data.
	if err := decoder.Decode(&struct{}{}); !errors.Is(err, io.EOF) {
		writeJSON(w, http.StatusBadRequest, map[string]string{
			"error": "request must contain one JSON object",
		})
		return
	}

	parsed, err := url.ParseRequestURI(req.URL)
	if err != nil ||
		parsed.Host == "" ||
		(parsed.Scheme != "http" && parsed.Scheme != "https") {
		writeJSON(w, http.StatusBadRequest, map[string]string{
			"error": "url must be an absolute HTTP or HTTPS URL",
		})
		return
	}

	code, err := generateCode()
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{
			"error": "could not generate short code",
		})
		return
	}

	mapping := URLMapping{
		Code:        code,
		OriginalURL: parsed.String(),
	}

	if err := h.store.Save(r.Context(), mapping); err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{
			"error": "could not save URL mapping",
		})
		return
	}

	writeJSON(w, http.StatusCreated, createURLResponse{
		Code:        code,
		ShortURL:    h.base_url + "/" + code,
		OriginalURL: mapping.OriginalURL,
	})
}

func (h *Handler) Redirect(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", http.MethodGet)
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	code := strings.TrimPrefix(r.URL.Path, "/")
	if code == "" || strings.Contains(code, "/") {
		http.NotFound(w, r)
		return
	}

	mapping, err := h.store.Get(r.Context(), code)

	if errors.Is(err, ErrNotFound) {
		http.NotFound(w, r)
		return
	}

	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, mapping.OriginalURL, http.StatusFound)
}

func writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(value)
}

func generateCode() (string, error) {
	bytes := make([]byte, 6)

	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(bytes), nil
}
