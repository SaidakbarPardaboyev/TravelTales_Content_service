package postgres

import (
	"log"
	"testing"
	pb "travel/genproto/interactions"
)

func NewIntRepo() *InterationsRepo {
	db, err := ConnectDB()
	if err != nil {
		log.Panic(err)
	}
	return NewInterationsRepo(db)
}

func TestCreateComment(t *testing.T) {
	req := pb.RequestCreateComment{
		StoryId:  "24c22836-26fa-486d-b660-262e123a1a5c",
		AuthorId: "5960639a-383e-4437-9f1a-9657f9f99964",
		Content:  "I have never seen story like this",
	}
	_, err := NewIntRepo().CreateComment(&req)
	if err != nil {
		t.Error(err)
	}
}

func TestGetComments(t *testing.T) {
	req := pb.RequestGetComments{
		StoryId: "24c22836-26fa-486d-b660-262e123a1a5c",
		Page:    0,
		Limit:   10,
	}

	_, err := NewIntRepo().GetComments(&req)
	if err != nil {
		t.Error(err)
	}
}

func TestCountComments(t *testing.T) {
	_, err := NewIntRepo().CountComments("24c22836-26fa-486d-b660-262e123a1a5c")
	if err != nil {
		t.Error(err)
	}
}

func TestLikeStory(t *testing.T) {
	err := NewIntRepo().LikeStory("24c22836-26fa-486d-b660-262e123a1a5c")
	if err != nil {
		t.Error(err)
	}
}

func TestCreateLike(t *testing.T) {
	req := pb.RequestLikeStory{
		UserId:  "9446b610-2ee7-46b4-98a1-ff905b016d2b",
		StoryId: "a9a0e266-ef7b-459e-aa28-ee091fc3eafb",
	}

	err := NewIntRepo().CreateLike(&req)
	if err != nil {
		t.Error(err)
	}
}
