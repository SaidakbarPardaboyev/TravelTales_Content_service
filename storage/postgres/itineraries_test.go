package postgres

import (
	"log"
	"testing"
	pb "travel/genproto/itineraries"
)

func NewItinarRepo() *ItinerariesRepo {
	db, err := ConnectDB()
	if err != nil {
		log.Panic(err)
	}
	return NewItinerariesRepo(db)
}

func TestCreateItineraries(t *testing.T) {
	req := pb.RequestCreateItineraries{
		AutherId:    "030c9cdc-c410-4e94-a5f6-4152fd4eafcb",
		Title:       "dfghjk",
		Description: "dsfhgfs dgh fghf hgh f",
		StartDate:   "2024-07-16",
		EndDate:     "2024-07-16",
		Destinations: []*pb.Destination{{
			Name:       "Tashkent",
			StartDate:  "2024-07-16",
			EndDate:    "2024-07-16",
			Activities: []string{"swimming"},
		}}}
	_, err := NewItinarRepo().CreateItineraries(&req)
	if err != nil {
		t.Error(err)
	}
}

func TestCreateItinerariesDestinations(t *testing.T) {
	des := []*pb.Destination{{
		Name:       "Tashkent",
		StartDate:  "2024-07-16",
		EndDate:    "2024-07-16",
		Activities: []string{"swimming"}},}
	
	err := NewItinarRepo().CreateItinerariesDestinations("b041bc66-3857-4720-a811-3d8a080a6343",
	des)
	if err != nil {
		t.Error(err)
	}
}

func TestCreateActivities(t *testing.T) {
	desId := "1e9c7187-229e-42cc-957a-dd343d0c51f3"
	ac := []string{"swimming", "doing sport"}
	err := NewItinarRepo().CreateActivities(desId, &ac)
	if err != nil {
		t.Error(err)
	}
}
