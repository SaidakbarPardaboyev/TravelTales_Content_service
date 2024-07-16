package models

type Story struct {
	Id             string
	Title          string
	AuthorId       string
	Location       string
	Likes_count    int
	Comments_count int
	Created_at     string
}
