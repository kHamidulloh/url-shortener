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
	_ "github.com/lib/pq"
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

	// Создание таблицы, если она не существует
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS urls (
			id SERIAL PRIMARY KEY,
			short_url VARCHAR(255) NOT NULL,
			original_url TEXT NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);
	`)
	if err != nil {
		log.Fatalf("❌ Error creating table: %v", err)
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
