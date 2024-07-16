package models

type Story struct {
	Id             string
	Title          string
	AuthorId       string
	Location       string
	LikesCount    int
	CommentsCount int
	CreatedAt     string
}

type StoryFullInfo struct {
	Id             string
	Title          string
	Content string
	AuthorId       string
	Location       string
	LikesCount    int
	CommentsCount int
	CreatedAt     string
	UpdatedAt string
}