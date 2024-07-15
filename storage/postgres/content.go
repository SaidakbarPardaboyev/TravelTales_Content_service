package postgres

import (
	"database/sql"
	"log/slog"
	"travel/pkg/logger"
)

type ContentRepo struct {
	Logger *slog.Logger
	DB     *sql.DB
}

func NewUserRepo(db *sql.DB) *ContentRepo {
	logger := logger.NewLogger()
	return &ContentRepo{
		Logger: logger,
		DB:     db,
	}
}
