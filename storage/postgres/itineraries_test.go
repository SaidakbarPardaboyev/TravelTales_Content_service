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
		Activities: []string{"swimming"}}}

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

func TestEditItineraries(t *testing.T) {
	req := pb.RequestEditItineraries{
		Id:           "00d47248-2563-4494-9561-d8c10749b8b6",
		Title:        "Saidakbar",
		Description:  "Saidakbar Pradaboyev",
		StartDate:    "2024-07-12",
		EndDate:      "2024-07-16",
		Destinations: []*pb.DestinationEdit{},
	}

	tx, err := NewItinarRepo().DB.Begin()
	if err != nil {
		t.Error(err)
	}
	defer tx.Commit()

	err = EditItineraries(tx, &req)
	if err != nil {
		tx.Rollback()
		t.Error(err)
	}
}

func TestEditItinerariesDestinations(t *testing.T) {
	req := []*pb.DestinationEdit{{
		Id:        "ef63014f-094f-4331-8a7d-9003c0d5d26e",
		Name:      "Jizzax",
		StartDate: "2024-07-05",
		EndDate:   "2024-07-30",
		Activities: []*pb.Activity{{
			Id:       "c9694101-61aa-414c-9834-2277cc654a7d",
			Activity: "Learing engling",
		}, {
			Id:       "3b6f1610-4370-48af-898e-f472506f6a5f",
			Activity: "Lering arab language",
		}},
	}}

	tx, err := NewItinarRepo().DB.Begin()
	if err != nil {
		t.Error(err)
	}
	defer tx.Commit()
	err = EditItinerariesDestinations(tx, req)
	if err != nil {
		tx.Rollback()
		t.Error(err)
	}
}

func TestGetAllItineraries(t *testing.T) {
	req := pb.RequestGetAllItineraries{
		Page:  0,
		Limit: 10,
	}

	_, err := NewItinarRepo().GetAllItineraries(&req)
	if err != nil {
		t.Error(err)
	}
}

func TestFindNumberOfItineraries(t *testing.T) {
	_, err := NewItinarRepo().FindNumberOfItineraries()
	if err != nil {
		t.Error(err)
	}
}

func TestGetItinerariesFullInfo(t *testing.T) {
	_, err := NewItinarRepo().GetItinerariesFullInfo(
		"35d06507-51e2-41cc-87eb-a0107f4f8217")
	if err != nil {
		t.Error(err)
	}
}

func TestGetItinerariesDestinations(t *testing.T) {
	_, err := NewItinarRepo().GetItinerariesDestinations(
		"35d06507-51e2-41cc-87eb-a0107f4f8217")
	if err != nil {
		t.Error(err)
	}
}

func TestWriteCommentToItinerary(t *testing.T) {
	req := pb.RequestWriteCommentToItinerary{
		Content:     "It is so good journay",
		AuthorId:    "030c9cdc-c410-4e94-a5f6-4152fd4eafcb",
		ItineraryId: "00d47248-2563-4494-9561-d8c10749b8b6",
	}
	_, err := NewItinarRepo().WriteCommentToItinerary(&req)
	if err != nil {
		t.Error(err)
	}
}

func TestCreateDestination(t *testing.T) {
	res := pb.RequestCreateDestination{
		Name:              "Makka",
		Country:           "Misr",
		Description:       "It is so so gooood",
		BestTimeToVisit:   "Winter and Spring",
		AverageCostPerDay: 1500,
		Currency:          "UDT",
		Language:          "English",
		PopularityScore:   100,
	}

	_, err := NewItinarRepo().CreateDestination(&res)
	if err != nil {
		t.Error(err)
	}
}

func TestGetTopDestinations(t *testing.T) {
	req := pb.RequestGetDestinations{
		Limit: 5,
		Page:  0,
	}
	_, err := NewItinarRepo().GetTopDestinations(&req)
	if err != nil {
		t.Error(err)
	}
}

// func Test(t *testing.T) {

// }

// func Test(t *testing.T) {

// }

// func Test(t *testing.T) {

// }

// func Test(t *testing.T) {

// }

// func Test(t *testing.T) {

// }

// func Test(t *testing.T) {

// }

// func Test(t *testing.T) {

// }

// func Test(t *testing.T) {

// }

// func Test(t *testing.T) {

// }
