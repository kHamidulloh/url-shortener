package repository

import (
	"github.com/jmoiron/sqlx"
)

type URLRepository struct {
	db *sqlx.DB
}

func NewURLRepository(db *sqlx.DB) *URLRepository {
	return &URLRepository{db: db}
}

func (r *URLRepository) Save(original, short string) error {
	query := `INSERT INTO urls (original_url,short_url) VALUES ($1, $2)`
	_, err := r.db.Exec(query, original, short)
	return err
}

func (r *URLRepository) GetOriginal(short string) (string, error) {
	var original string
	query := `SELECT original_url FROM urls WHERE short_url = $1`
	err := r.db.Get(&original, query, short)
	return original, err
}
