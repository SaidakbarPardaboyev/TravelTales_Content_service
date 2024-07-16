package service

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"time"
	pb "travel/genproto/stories"
	pbUser "travel/genproto/users"
	"travel/pkg/connections"
	"travel/pkg/logger"
	"travel/storage/postgres"
)

type Stories struct {
	pb.UnimplementedStoriesServer
	Logger      *slog.Logger
	StoriesRepo *postgres.StoriesRepo
	UserClient  pbUser.UsersClient
}

func NewContentService(db *sql.DB) *Stories {
	storiesRepo := postgres.NewStoriesRepo(db)
	Logger := logger.NewLogger()
	userClient := connections.NewUserClient()
	return &Stories{
		Logger:      Logger,
		StoriesRepo: storiesRepo,
		UserClient:  userClient,
	}
}

func (s *Stories) CreateStory(ctx context.Context, in *pb.RequestCreateStory) (
	*pb.ResponseCreateStory, error) {
	// checking user exists
	valid, err := s.UserClient.ValidateUser(ctx, &pbUser.RequestGetProfile{Id: in.AuthorId})
	if err != nil || !valid.Success {
		s.Logger.Error(fmt.Sprintf("error with validating user: %s", err))
		return nil, fmt.Errorf("error: invalid userID: %s", err)
	}

	id, err := s.StoriesRepo.CreateStory(in)
	if err != nil {
		s.Logger.Error(fmt.Sprintf("error with creating story: %s", err))
		return nil, err
	}

	resp := pb.ResponseCreateStory{
		Id:        id,
		Title:     in.Title,
		Content:   in.Content,
		Location:  in.Location,
		Tags:      in.Tags,
		AuthorId:  in.AuthorId,
		CreatedAt: time.Now().String(),
	}

	err = s.StoriesRepo.CreateStoryTags(id, &in.Tags)
	if err != nil {
		s.Logger.Error(fmt.Sprintf("error with creating story tags: %s", err))
		return nil, err
	}
	return &resp, nil
}

func (s *Stories) EditStory(ctx context.Context, in *pb.RequestEditStory) (
	*pb.ResponseEditStory, error) {

	authorId, err := s.StoriesRepo.EditStory(in)
	if err != nil {
		s.Logger.Error(fmt.Sprintf("error with editing story tags: %s", err))
		return nil, err
	}

	err = s.StoriesRepo.DeleteStoryTags(in.Id)
	if err != nil {
		s.Logger.Error(fmt.Sprintf("error with deleting story tags: %s", err))
		return nil, err
	}

	err = s.StoriesRepo.CreateStoryTags(in.Id, &in.Tags)
	if err != nil {
		s.Logger.Error(fmt.Sprintf("error with creating story tags: %s", err))
		return nil, err
	}

	resp := pb.ResponseEditStory{
		Id:        in.Id,
		Title:     in.Title,
		Content:   in.Content,
		Location:  in.Location,
		Tags:      in.Tags,
		AuthorId:  authorId,
		UpdatedAt: time.Now().String(),
	}
	return &resp, nil
}

// func (s *Stories) GetStories(ctx context.Context, in *pb.RequestGetStories) (*pb.ResponseGetStories, error)
// func (s *Stories) GetStoryFullInfo(ctx context.Context, in *pb.RequestGetStoryFullInfo) (*pb.ResponseGetStoryFullInfo, error)

// func (s *Stories) DeleteStory(ctx context.Context, in *pb.RequestDeleteStory) (*pb.ResponseDeleteStory, error)
