package postgres

import (
	"log"
	"testing"
	pb "travel/genproto/stories"
)

func NewRepo() *StoriesRepo {
	db, err := ConnectDB()
	if err != nil {
		log.Panic(err)
	}
	return NewStoriesRepo(db)
}

func TestCreateStory(t *testing.T) {

	req := pb.RequestCreateStory{
		AuthorId: "ff9ae172-18f9-4f81-98ff-5db600ce05a7",
		Title:    "Go Home",
		Content:  "About going home",
		Location: "Uzbekistan",
		Tags:     []string{},
		Images:   []string{"go", "home"},
	}

	_, err := NewRepo().CreateStory(&req)
	if err != nil {
		t.Error(err)
	}
}

func TestCreateStoryTags(t *testing.T) {

	storyId := "cefbcf04-172e-4b01-88ed-763ab5848d45"
	tags := []string{"gfc", "dfg"}

	err := NewRepo().CreateStoryTags(storyId, &tags)
	if err != nil {
		t.Error(err)
	}
}

func TestEditStory(t *testing.T) {
	req := pb.RequestEditStory{
		Id:       "24c22836-26fa-486d-b660-262e123a1a5c",
		Title:    "Sleeping well",
		Content:  "About Sleeping well",
		Location: "in mindset",
		Tags:     []string{"go home"},
		Images:   []string{},
	}

	_, err := NewRepo().EditStory(&req)
	if err != nil {
		t.Error(err)
	}
}

func TestDeleteStoryTags(t *testing.T) {
	storyId := "24c22836-26fa-486d-b660-262e123a1a5c"

	err := NewRepo().DeleteStoryTags(storyId)
	if err != nil {
		t.Error(err)
	}
}

func TestGetStories(t *testing.T) {
	req := pb.RequestGetStories{
		Page:  0,
		Limit: 10,
	}
	_, err := NewRepo().GetStories(&req)
	if err != nil {
		t.Error(err)
	}
}

func TestFindNumberOfStories(t *testing.T) {
	_, err := NewRepo().FindNumberOfStories()
	if err != nil {
		t.Error(err)
	}
}

func TestGetStoryFullInfo(t *testing.T) {
	_, err := NewRepo().GetStoryFullInfo("cefbcf04-172e-4b01-88ed-763ab5848d45")
	if err != nil {
		t.Error(err)
	}
}

func TestGetStoryTags(t *testing.T) {
	_, err := NewRepo().GetStoryTags("cefbcf04-172e-4b01-88ed-763ab5848d45")
	if err != nil {
		t.Error(err)
	}
}

func TestDeleteStory(t *testing.T) {
	err := NewRepo().DeleteStory("cefbcf04-172e-4b01-88ed-763ab5848d45")
	if err != nil {
		t.Error(err)
	}
}
