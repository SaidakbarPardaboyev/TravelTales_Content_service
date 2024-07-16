package postgres

import (
	"database/sql"
	"fmt"
	"log/slog"
	"time"
	pb "travel/genproto/stories"
	"travel/models"
	"travel/pkg/logger"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

type StoriesRepo struct {
	Logger *slog.Logger
	DB     *sql.DB
}

func NewStoriesRepo(db *sql.DB) *StoriesRepo {
	logger := logger.NewLogger()
	return &StoriesRepo{
		Logger: logger,
		DB:     db,
	}
}

func (s *StoriesRepo) CreateStory(story *pb.RequestCreateStory) (
	string, error) {

	query := `
		insert into stories(
			id, title, content, location, author_id, images
		) values (
			$1, $2, $3, $4, $5, $6
		)
	`

	newId := uuid.NewString()
	_, err := s.DB.Exec(query, newId, story.Title, story.Content, story.Location,
		story.AuthorId, pq.Array(story.Images))

	return newId, err
}

func (s *StoriesRepo) CreateStoryTags(storyId string, tags *[]string) error {

	query := `
		insert into story_tags(
			story_id, tag
		) values (
			$1, $2 
		)
	`

	for _, tag := range *tags {
		res, err := s.DB.Exec(query, storyId, tag)
		if err != nil {
			return err
		}
		if num, _ := res.RowsAffected(); num <= 0 {
			return fmt.Errorf("duplicated values")
		}
	}
	return nil
}

func (s *StoriesRepo) EditStory(story *pb.RequestEditStory) (string, error) {

	query := `
		update
			stories
		set
			title = $1,
			content = $2,
			location = $3,
			images = $4
		where
			id = $5 and 
			deleted_at is null
		returning author_id
	`
	var AuthorId string
	err := s.DB.QueryRow(query, story.Title, story.Content, story.Location,
		pq.Array(story.Images), story.Id).Scan(&AuthorId)
	return AuthorId, err
}

func (s *StoriesRepo) DeleteStoryTags(storyId string) error {

	query := `
		delete from
			story_tags
		where
			story_id = $1
	`

	_, err := s.DB.Exec(query, storyId)
	return err
}

func (s *StoriesRepo) GetStories(filter *pb.RequestGetStories) (
	*[]models.Story, error) {

	query := `
		select
			id, title, author_id, location, likes_count, comments_count, 
			created_at
		from
			stories
		where
			deleted_at is null
		limit $1
		offset $2
	`

	rows, err := s.DB.Query(query, filter.Limit, filter.Limit*filter.Page)
	if err != nil {
		return nil, err
	}

	stories := []models.Story{}
	for rows.Next() {
		var story models.Story
		err := rows.Scan(&story.Id, &story.Title, &story.AuthorId,
			&story.Location, &story.LikesCount, &story.CommentsCount,
			&story.CreatedAt)
		if err != nil {
			return nil, err
		}
		stories = append(stories, story)
	}

	return &stories, nil
}

func (s *StoriesRepo) FindNumberOfStories() (int, error) {

	query := `
		select
			count(*)
		from
			stories
		where
			deleted_at is null
	`

	count := 0
	err := s.DB.QueryRow(query).Scan(&count)
	return count, err
}

func (s *StoriesRepo) GetStoryFullInfo(id string) (
	*models.StoryFullInfo, error) {

	query := `
		select
			id, title, content, author_id, location, likes_count, comments_count, 
			created_at, updated_at
		from
			stories
		where
			id = $1 and 
			deleted_at is null
	`

	res := models.StoryFullInfo{}
	err := s.DB.QueryRow(query, id).Scan(&res.Id, &res.Title, &res.Content,
		&res.AuthorId, &res.Location, &res.LikesCount, &res.CommentsCount,
		&res.CreatedAt, &res.UpdatedAt)
	return &res, err
}

func (s *StoriesRepo) GetStoryTags(storyId string) (*[]string, error) {

	query := `
		select
			array_agg(tag)
		from
			story_tags
		where
			story_id = $1
	`

	res := []string{}
	err := s.DB.QueryRow(query, storyId).Scan(pq.Array(&res))
	return &res, err
}

func (s *StoriesRepo) DeleteStory(id string) error {

	query := `
		update 
			stories
		set
			deleted_at = $1
		where
			id = $2 and 
			deleted_at is null
	`

	res, err := s.DB.Exec(query, time.Now(), id)
	if num, _ := res.RowsAffected(); num <= 0 {
		return fmt.Errorf("story is not found with the id")
	}
	return err
}
