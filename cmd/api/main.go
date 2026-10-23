package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/prathibha97/tiny-link/internal/shortener"
)

func main() {
	databaseUrl := os.Getenv("DATABASE_URL")
	if databaseUrl == "" {
		log.Fatal("DATABASE_URL is required")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	store, err := shortener.NewPostgresPool(ctx, databaseUrl)

	if err != nil {
		log.Fatalf("initialize database: %v", err)
	}
	defer store.Close()
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
