package config

import (
	_ "database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq"
)

type Config struct {
	DBUrl string
}

func Load() *Config {
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_DBNAME")

	dbURL := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		user, password, host, port, dbname,
	)
	return &Config{
		DBUrl: dbURL,
	}
}
