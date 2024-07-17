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

func (i *ItinerariesRepo) GetItinerariesFullInfo(id string) (
	*models.ItineraryFullInfo, error) {

	query := `
		select
			id, title, description, start_date, end_date, author_id, 
			likes_count, comments_count, created_at, updated_at
		from 
			itineraries
		where
			id = $1 and
			deleted_at is null`

	resp := models.ItineraryFullInfo{}
	err := i.DB.QueryRow(query, id).Scan(&resp.Id, &resp.Title,
		&resp.Description, &resp.StartDate, &resp.EndDate,
		&resp.AutherId, &resp.LikesCount, &resp.CommentsCount,
		&resp.CreatedAt, &resp.UpdatedAt)
	return &resp, err
}

func (i *ItinerariesRepo) GetItinerariesDestinations(id string) (*[]*pb.DestinationEdit,
	error) {

	query := `
		select
			id, name, start_date, end_date
		from
			itinerary_destinations
		where
			itinerary_id = $1 and
			deleted_at is null`

	resp := []*pb.DestinationEdit{}
	rows, err := i.DB.Query(query, id)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		des := pb.DestinationEdit{}
		err = rows.Scan(&des.Id, &des.Name, &des.StartDate, &des.EndDate)
		if err != nil {
			return nil, err
		}

		activities, err := i.GetDestinationActivities(des.Id)
		if err != nil {
			return nil, err
		}

		des.Activities = *activities
		resp = append(resp, &des)
	}
	return &resp, nil
}

func (i *ItinerariesRepo) GetDestinationActivities(desId string) (
	*[]*pb.Activity, error) {

	query := `
		select
			id, activity
		from
			itinerary_activities
		where
			destination_id = $1 and
			deleted_at is null`

	activities := []*pb.Activity{}
	rows, err := i.DB.Query(query, desId)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		activity := pb.Activity{}
		err = rows.Scan(&activity.Id, &activity.Activity)
		if err != nil {
			return nil, err
		}
		activities = append(activities, &activity)
	}
	return &activities, nil
}

func (i *ItinerariesRepo) WriteCommentToItinerary(req *pb.RequestWriteCommentToItinerary) (
	string, error) {

	query := `
		insert into commentsForItinerary(
			id, content, author_id, itinerary_id
		) values (
			$1, $2, $3, $4 
		)`

	newId := uuid.NewString()
	_, err := i.DB.Exec(query, newId, req.Content, req.AuthorId, req.ItineraryId)
	return newId, err
}

func (i *ItinerariesRepo) CreateDestination(req *pb.RequestCreateDestination) (
	string, error) {
	query := `
		insert into destinations(
			id, name, country, description, best_time_to_visit, 
			average_cost_per_day, currency, language, popularity_score
		) values (
			$1, $2, $3, $4, $5, $6, $7, $8, $9  
		)`

	newId := uuid.NewString()
	_, err := i.DB.Exec(query, newId, req.Name, req.Country, req.Description,
		req.BestTimeToVisit, req.AverageCostPerDay, req.Currency,
		req.Language, req.PopularityScore)
	return newId, err
}

func (i *ItinerariesRepo) GetTopDestinations(req *pb.RequestGetDestinations) (
	*pb.ResponseGetDestinations, error) {

	query := `
		select
			id, name, country, description
		from
			destinations
		order by
			popularity_score desc
		limit $1
		offset $2
	`

	rows, err := i.DB.Query(query, req.Limit, req.Limit*req.Page)
	if err != nil {
		return nil, err
	}
	res := pb.ResponseGetDestinations{}
	for rows.Next() {
		des := pb.DestionationInfo{}
		err := rows.Scan(&des.Id, &des.Name, &des.Country, &des.Description)
		if err != nil {
			return nil, err
		}
		res.Destinations = append(res.Destinations, &des)
	}
	return &res, nil
}

// func (i *ItinerariesRepo) CreateItineraries(req *pb.RequestCreateItineraries) () {}
// func (i *ItinerariesRepo) CreateItineraries(req *pb.RequestCreateItineraries) () {}
