package postgres

import (
	"database/sql"
	"fmt"
	"log/slog"
	"time"
	pb "travel/genproto/itineraries"
	"travel/models"
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

func EditItineraries(tx *sql.Tx, req *pb.RequestEditItineraries) error {

	query := `
		update
			itineraries
		set
			title = $1,
			description = $2,
			start_date = $3,
			end_date = $4,
			updated_at = $5
		where
			id = $6 and
			deleted_at is null`

	res, err := tx.Exec(query, req.Title, req.Description, req.StartDate,
		req.EndDate, time.Now(), req.Id)
	if err != nil {
		return err
	}
	if num, _ := res.RowsAffected(); num <= 0 {
		return fmt.Errorf("itinerary not found with the id: %s", req.Id)
	}
	return nil
}

func EditItinerariesDestinations(tx *sql.Tx,
	destinations []*pb.DestinationEdit) error {

	query := `
		update
			itinerary_destinations
		set
			name = $1,
			start_date = $2, 
			end_date = $3
		where
			id = $4 and
			deleted_at is null`

	for _, des := range destinations {
		res, err := tx.Exec(query, des.Name, des.StartDate,
			des.EndDate, des.Id)
		if err != nil {
			return err
		}
		if num, _ := res.RowsAffected(); num <= 0 {
			return fmt.Errorf("destination not found with the id: %s", des.Id)
		}
		err = EditActivities(tx, &des.Activities)
		if err != nil {
			return err
		}
	}
	return nil
}

func EditActivities(tx *sql.Tx, activities *[]*pb.Activity) error {

	query := `
		update
			itinerary_activities
		set
			activity = $1
		where
			id = $2 and
			deleted_at is null`

	for _, act := range *activities {
		res, err := tx.Exec(query, act.Activity, act.Id)
		if err != nil {
			return err
		}
		if num, _ := res.RowsAffected(); num <= 0 {
			return fmt.Errorf("activity not found with the id: %s", act.Id)
		}
	}
	return nil
}

func (i *ItinerariesRepo) GetAllItineraries(req *pb.RequestGetAllItineraries) (
	*[]models.Itinerary, error) {

	query := `
		select
			id, title, description, start_date, end_date, author_id, 
			likes_count, comments_count, created_at
		from 
			itineraries
		where
			deleted_at is null
		limit $1
		offset $2
	`

	rows, err := i.DB.Query(query, req.Limit, req.Limit*req.Page)
	if err != nil {
		return nil, err
	}

	itineraries := []models.Itinerary{}
	for rows.Next() {
		itinerary := models.Itinerary{}
		err := rows.Scan(&itinerary.Id, &itinerary.Title, &itinerary.Description,
			&itinerary.StartDate, &itinerary.EndDate, &itinerary.AutherId,
			&itinerary.LikesCount, &itinerary.CommentsCount, &itinerary.CreatedAt)
		if err != nil {
			return nil, err
		}
		itineraries = append(itineraries, itinerary)
	}
	return &itineraries, nil
}

func (i *ItinerariesRepo) FindNumberOfItineraries() (int, error) {

	query := `
		select
			count(*)
		from
			itineraries
		where
			deleted_at is null
	`

	count := 0
	err := i.DB.QueryRow(query).Scan(&count)
	return count, err
}

// func (i *ItinerariesRepo) CreateItineraries(req *pb.RequestCreateItineraries) () {}
// func (i *ItinerariesRepo) CreateItineraries(req *pb.RequestCreateItineraries) () {}
// func (i *ItinerariesRepo) CreateItineraries(req *pb.RequestCreateItineraries) () {}
// func (i *ItinerariesRepo) CreateItineraries(req *pb.RequestCreateItineraries) () {}
// func (i *ItinerariesRepo) CreateItineraries(req *pb.RequestCreateItineraries) () {}
// func (i *ItinerariesRepo) CreateItineraries(req *pb.RequestCreateItineraries) () {}
// func (i *ItinerariesRepo) CreateItineraries(req *pb.RequestCreateItineraries) () {}
// func (i *ItinerariesRepo) CreateItineraries(req *pb.RequestCreateItineraries) () {}
