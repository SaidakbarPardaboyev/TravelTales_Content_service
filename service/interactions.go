package service

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"time"
	pb "travel/genproto/interactions"
	pbUser "travel/genproto/users"
	"travel/pkg/connections"
	"travel/pkg/logger"
	"travel/storage/postgres"
)

type Interations struct {
	pb.UnimplementedInteractionsServer
	Logger          *slog.Logger
	InterationsRepo *postgres.InterationsRepo
	UserClient      pbUser.UsersClient
}

func NewInterationsService(db *sql.DB) *Interations {
	InterationsRepo := postgres.NewInterationsRepo(db)
	Logger := logger.NewLogger()
	userClient := connections.NewUserClient()
	return &Interations{
		Logger:          Logger,
		InterationsRepo: InterationsRepo,
		UserClient:      userClient,
	}
}

func (i *Interations) CreateComment(ctx context.Context, in *pb.RequestCreateComment) (
	*pb.ResponseCreateComment, error) {

	// checking user exists
	valid, err := i.UserClient.ValidateUser(ctx, &pbUser.RequestGetProfile{Id: in.AuthorId})
	if err != nil || !valid.Success {
		i.Logger.Error(fmt.Sprintf("error with validating user: %s", err))
		return nil, fmt.Errorf("error: invalid userID: %s", err)
	}

	id, err := i.InterationsRepo.CreateComment(in)
	if err != nil {
		i.Logger.Error(fmt.Sprintf("error with creating comment: %s", err))
		return nil, err
	}

	resp := pb.ResponseCreateComment{
		Id:        id,
		StoryId:   in.StoryId,
		Content:   in.Content,
		AuthorId:  in.AuthorId,
		CreatedAt: time.Now().String(),
	}
	return &resp, nil
}

func (i *Interations) GetComments(ctx context.Context, in *pb.RequestGetComments) (
	*pb.ResponseGetComments, error) {
	comments, err := i.InterationsRepo.GetComments(in)
	if err != nil {
		i.Logger.Error(fmt.Sprintf("error with getting comments by story Id: %s", err))
		return nil, err
	}

	resp := pb.ResponseGetComments{}

	for _, com := range *comments {
		comment := pb.Comment{
			Id:        com.Id,
			Content:   com.Content,
			CreatedAt: com.CreatedAt,
		}

		author, err := i.UserClient.GetAuthorInfo(ctx, &pbUser.RequestGetAuthorInfo{
			Id: com.AuthorId,
		})
		if err != nil {
			if err.Error() == "rpc error: code = Unknown desc = sql: no rows in result set" {
				continue
			}
			i.Logger.Error(fmt.Sprintf("error with getting Author info: %s", err))
			return nil, err
		}

		comment.Author = &pb.Author{
			Id:       author.Id,
			Username: author.Username,
		}
		resp.Comments = append(resp.Comments, &comment)
	}

	resp.Limit = in.Limit
	resp.Page = in.Page

	count, err := i.InterationsRepo.CountComments(in.StoryId)
	if err != nil {
		return nil, err
	}
	resp.Total = int64(count)
	return &resp, err
}

// func (i *Interations) LikeStory(ctx context.Context, in *pb.RequestLikeStory) (*pb.ResponseLikeStory, error)
