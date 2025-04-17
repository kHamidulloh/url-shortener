package main

import (
	"log"
	"net/http"
	"time"

	"url-shortener/internal/config"
	"url-shortener/internal/handler"
	"url-shortener/internal/repository"
	"url-shortener/internal/service"

	"github.com/go-chi/chi/v5"
	"github.com/jmoiron/sqlx"
)

func main() {
	cfg := config.Load()

	var db *sqlx.DB
	var err error
	maxRetries := 10

	for i := 1; i <= maxRetries; i++ {
		db, err = sqlx.Connect("postgres", cfg.DBUrl)
		if err == nil {
			log.Println("✅ Connected to database.")
			break
		}

		log.Printf("⏳ Attempt %d: waiting for database... (%v)", i, err)
		time.Sleep(3 * time.Second)
	}

	if err != nil {
		log.Fatalf("❌ Could not connect to database after %d attempts: %v", maxRetries, err)
	}

	repo := repository.NewURLRepository(db)
	svc := service.NewURLService(repo)
	h := handler.NewURLHandler(svc)

	r := chi.NewRouter()
	r.Post("/shorten", h.ShortenURL)
	r.Get("/{shortUrl}", h.ResolveURL)

	log.Println("🚀 Server started at :8000")
	http.ListenAndServe(":8000", r)
}
