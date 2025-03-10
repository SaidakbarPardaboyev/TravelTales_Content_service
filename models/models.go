package models

type Story struct {
	Id            string
	Title         string
	AuthorId      string
	Location      string
	LikesCount    int
	CommentsCount int
	CreatedAt     string
}

type StoryFullInfo struct {
	Id            string
	Title         string
	Content       string
	AuthorId      string
	Location      string
	LikesCount    int
	CommentsCount int
	CreatedAt     string
	UpdatedAt     string
}

type Comment struct {
	Id        string
	Content   string
	AuthorId  string
	CreatedAt string
}

type Itinerary struct {
	Id            string
	Title         string
	Description   string
	StartDate     string
	EndDate       string
	AutherId      string
	LikesCount    int
	CommentsCount int
	CreatedAt     string
}

type ItineraryFullInfo struct {
	Id            string
	Title         string
	Description   string
	StartDate     string
	EndDate       string
	AutherId      string
	LikesCount    int
	CommentsCount int
	CreatedAt     string
	UpdatedAt     string
}
