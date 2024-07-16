package postgres

import (
	"database/sql"
	"fmt"
	"log/slog"
	pb "travel/genproto/interactions"
	"travel/models"
	"travel/pkg/logger"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

type InterationsRepo struct {
	Logger *slog.Logger
	DB     *sql.DB
}

func NewInterationsRepo(db *sql.DB) *InterationsRepo {
	logger := logger.NewLogger()
	return &InterationsRepo{
		Logger: logger,
		DB:     db,
	}
}

func (i *InterationsRepo) CreateComment(req *pb.RequestCreateComment) (
	string, error) {

	query := `
		insert into comments(
			id, content, author_id, story_id
		) values (
			$1, $2, $3, $4 
		)
	`

	newId := uuid.NewString()
	_, err := i.DB.Exec(query, newId, req.Content, req.AuthorId, req.StoryId)
	return newId, err
}

func (i *InterationsRepo) GetComments(req *pb.RequestGetComments) (
	*[]models.Comment, error) {

	query := `
		select 
			id, content, author_id, created_at
		from
			comments
		where
			story_id = $1 and 
			deleted_at is null
		limit $2
		offset $3
	`

	rows, err := i.DB.Query(query, req.StoryId, req.Limit,
		req.Limit*req.Page)
	if err != nil {
		return nil, err
	}
	comments := []models.Comment{}
	for rows.Next() {
		comment := models.Comment{}
		err := rows.Scan(&comment.Id, &comment.Content, &comment.AuthorId,
			&comment.CreatedAt)
		if err != nil {
			return nil, err
		}

		comments = append(comments, comment)
	}
	return &comments, nil
}

func (i *InterationsRepo) CountComments(storyId string) (
	int, error) {

	query := `
		select 
			count(*)
		from
			comments
		where
			story_id = $1 and 
			deleted_at is null
	`

	res := 0
	err := i.DB.QueryRow(query, storyId).Scan(&res)
	return res, err
}

func (i *InterationsRepo) LikeStory(storyId string) error {

	query := `
		update
			stories
		set
			likes_count = likes_count + 1
		where
			id = $1 and 
			deleted_at is null
	`

	res, err := i.DB.Exec(query, storyId)
	if err != nil {
		return err
	}
	if num, _ := res.RowsAffected(); num <= 0 {
		return fmt.Errorf("story not found with the id")
	}
	return nil
}

func (i *InterationsRepo) CreateLike(req *pb.RequestLikeStory) error {

	query := `
		insert into likes (
			user_id, story_id
		) values (
			$1, $2 
		)
	`

	_, err := i.DB.Exec(query, req.UserId, req.StoryId)
	return err
}

// func (i *InterationsRepo) CreateComment(req *pb.RequestCreateComment) (
// 	string, error) {

// 	query := `

// 	`
// }

// func (i *InterationsRepo) CreateComment(req *pb.RequestCreateComment) (
// 	string, error) {

// 	query := `

// 	`
// }

// func (i *InterationsRepo) CreateComment(req *pb.RequestCreateComment) (
// 	string, error) {

// 	query := `

// 	`
// }
