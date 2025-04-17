package config

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/lib/pq"
)

func InitDB() (*sql.DB, error) {
	dsn := "host=db user=postgres password=postgres dbname=url_shortener sslmode=disable"
	var db *sql.DB
	var err error

	maxRetries := 10
	for i := 1; i <= maxRetries; i++ {
		db, err = sql.Open("postgres", dsn)
		if err != nil {
			log.Printf("[Attempt %d] Failed to open DB: %v", i, err)
		} else {
			err = db.Ping()
			if err == nil {
				log.Println("✅ Connected to database.")
				return db, nil
			}
			log.Printf("[Attempt %d] Database not ready yet: %v", i, err)
		}

		time.Sleep(3 * time.Second)
	}

	return nil, fmt.Errorf("❌ failed to connect to database after %d attempts: %w", maxRetries, err)
}
