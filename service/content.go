package service

import (
	"database/sql"
	"log/slog"
	"travel/pkg/logger"
	"travel/storage/postgres"

	pb "travel/genproto/users"
)

type ContentService struct {
	pb.UnimplementedUsersServer
	Logger   *slog.Logger
	UserRepo *postgres.UserRepo
}

func NewContentService(db *sql.DB) *ContentService {
	userRepo := postgres.NewUserRepo(db)
	Logger := logger.NewLogger()
	return &ContentService{
		Logger:   Logger,
		UserRepo: userRepo,
	}
}
