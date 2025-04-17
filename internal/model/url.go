package model

type URL struct {
	ID        int    `db:"id"`
	Original  string `db:"original_url"`
	Short     string `db:"short_url"`
	CreatedAt string `db:"created_at"`
}
