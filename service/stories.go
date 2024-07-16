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

func (s *Stories) GetStories(ctx context.Context, in *pb.RequestGetStories) (
	*pb.ResponseGetStories, error) {
	stories, err := s.StoriesRepo.GetStories(in)
	if err != nil {
		s.Logger.Error(fmt.Sprintf("error with getting stories: %s", err))
		return nil, err
	}

	resp := pb.ResponseGetStories{}
	for _, val := range *stories {
		auther, err := s.UserClient.GetAuthorInfo(ctx,
			&pbUser.RequestGetAuthorInfo{Id: val.AuthorId})
		if err != nil {
			return nil, err
		}
		story := pb.StoryForGet{
			Id:    val.Id,
			Title: val.Title,
			Author: &pb.Author{
				Id:       auther.Id,
				Username: auther.Username,
			},
			Location:      val.Location,
			LikesCount:    int64(val.LikesCount),
			CommentsCount: int64(val.CommentsCount),
			CreatedAt:     val.CreatedAt,
		}

		resp.Stories = append(resp.Stories, &story)
	}

	countOfStories, err := s.StoriesRepo.FindNumberOfStories()
	if err != nil {
		s.Logger.Error(fmt.Sprintf("error with getting total stories count: %s", err))
		return nil, err
	}
	resp.Total = int64(countOfStories)
	resp.Limit = in.Limit
	resp.Page = in.Page

	return &resp, nil
}

func (s *Stories) GetStoryFullInfo(ctx context.Context, in *pb.RequestGetStoryFullInfo) (
	*pb.ResponseGetStoryFullInfo, error) {
	story, err := s.StoriesRepo.GetStoryFullInfo(in.Id)
	if err != nil {
		return nil, err
	}

	auther, err := s.UserClient.GetAuthorInfo(ctx,
		&pbUser.RequestGetAuthorInfo{Id: story.AuthorId})
	if err != nil {
		return nil, err
	}

	resp := pb.ResponseGetStoryFullInfo{
		Id:       story.Id,
		Title:    story.Title,
		Content:  story.Content,
		Location: story.Location,
		Author: &pb.AuthorForGetStoryFullInfo{
			Id:       auther.Id,
			Username: auther.Username,
			FullName: auther.FullName,
		},
		LikesCount:    int64(story.LikesCount),
		CommentsCount: int64(story.CommentsCount),
		CreatedAt:     story.CreatedAt,
		UpdatedAt:     story.UpdatedAt,
	}

	tags, err := s.StoriesRepo.GetStoryTags(resp.Id)
	if err != nil {
		return nil, err
	}
	resp.Tags = *tags

	return &resp, nil
}

func (s *Stories) DeleteStory(ctx context.Context, in *pb.RequestDeleteStory) (
	*pb.ResponseDeleteStory, error) {
	err := s.StoriesRepo.DeleteStory(in.StoryId)
	if err != nil {
		return nil, err
	}

	return &pb.ResponseDeleteStory{Message: "Story was deleted Successfully"}, nil
}
