package main

import (
	"log"
	"net/http"

	"github.com/prathibha97/tiny-link/internal/shortener"
)

func main() {
	store := shortener.NewURLStore()
	handler := shortener.NewHandler(store, "http://localhost:8080")

	mux := http.NewServeMux()
	mux.HandleFunc("/api/urls", handler.CreateURL)
	mux.HandleFunc("/", handler.Redirect)

	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}
	log.Println("URL shortener listening on :8080")

	if err := server.ListenAndServe(); err != nil &&
		err != http.ErrServerClosed {
		log.Fatal(err)
	}
}
