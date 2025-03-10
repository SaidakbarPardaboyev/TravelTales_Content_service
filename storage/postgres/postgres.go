package postgres

import (
	"database/sql"
	"fmt"
	"travel/config"

	_ "github.com/lib/pq"
)

func ConnectDB() (*sql.DB, error) {
	cfg := config.Load()

	conn := fmt.Sprintf(`host=%s port=%s user=%s dbname=%s password=%s 
	sslmode=disable`, cfg.DB_HOST, cfg.DB_PORT, cfg.DB_USER, cfg.DB_NAME,
		cfg.DB_PASSWORD)

	db, err := sql.Open("postgres", conn)
	return db, err
}
