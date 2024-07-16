package postgres

import (
	"database/sql"
	"fmt"
	"log/slog"
	pb "travel/genproto/stories"
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
			id = $5
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
