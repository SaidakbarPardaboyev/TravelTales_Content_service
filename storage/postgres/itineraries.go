package postgres

import (
	"database/sql"
	"log/slog"
	pb "travel/genproto/itineraries"
	"travel/pkg/logger"

	"github.com/google/uuid"
)

type ItinerariesRepo struct {
	Logger *slog.Logger
	DB     *sql.DB
}

func NewItinerariesRepo(db *sql.DB) *ItinerariesRepo {
	logger := logger.NewLogger()
	return &ItinerariesRepo{
		Logger: logger,
		DB:     db,
	}
}

func (i *ItinerariesRepo) CreateItineraries(req *pb.RequestCreateItineraries) (
	string, error) {

	query := `
		insert into itineraries(
			id, title, description, start_date, end_date, author_id
		) values (
			$1, $2, $3, $4, $5, $6 
		)
	`

	newId := uuid.NewString()
	_, err := i.DB.Exec(query, newId, req.Title, req.Description,
		req.StartDate, req.EndDate, req.AutherId)
	return newId, err
}

func (i *ItinerariesRepo) CreateItinerariesDestinations(itineraryId string,
	destinations []*pb.Destination) error {

	query := `
		insert into itinerary_destinations(
			id, itinerary_id, name, start_date, end_date
		) values (
			$1, $2, $3, $4, $5
		)`

	for _, des := range destinations {
		newId := uuid.NewString()
		_, err := i.DB.Exec(query, newId, itineraryId, des.Name, des.StartDate,
			des.EndDate)
		if err != nil {
			return err
		}
		err = i.CreateActivities(newId, &des.Activities)
		if err != nil {
			return err
		}
	}
	return nil
}

func (i *ItinerariesRepo) CreateActivities(desId string, req *[]string) error {

	query := `
		insert into itinerary_activities(
			id, destination_id, activity
		) values (
			$1, $2, $3
		)`

	for _, activity := range *req {
		newId := uuid.NewString()
		_, err := i.DB.Exec(query, newId, desId, activity)
		if err != nil {
			return err
		}
	}
	return nil
}

// func (i *ItinerariesRepo) CreateItineraries(req *pb.RequestCreateItineraries) () {}
// func (i *ItinerariesRepo) CreateItineraries(req *pb.RequestCreateItineraries) () {}
// func (i *ItinerariesRepo) CreateItineraries(req *pb.RequestCreateItineraries) () {}
// func (i *ItinerariesRepo) CreateItineraries(req *pb.RequestCreateItineraries) () {}
// func (i *ItinerariesRepo) CreateItineraries(req *pb.RequestCreateItineraries) () {}
// func (i *ItinerariesRepo) CreateItineraries(req *pb.RequestCreateItineraries) () {}
