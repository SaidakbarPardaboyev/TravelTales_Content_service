package service

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"time"
	pb "travel/genproto/itineraries"
	pbUser "travel/genproto/users"
	"travel/pkg/connections"
	"travel/pkg/logger"
	"travel/storage/postgres"
	"travel/storage/redis"

	rd "github.com/redis/go-redis/v9"
)

type Itineraries struct {
	pb.UnimplementedItinerariesServer
	Logger          *slog.Logger
	ItinerariesRepo *postgres.ItinerariesRepo
	UserClient      pbUser.UsersClient
	Redis           redis.DestinationRedisClient
}

func NewItinerariesService(db *sql.DB) *Itineraries {
	ItinerariesRepo := postgres.NewItinerariesRepo(db)
	Logger := logger.NewLogger()
	userClient := connections.NewUserClient()
	redisClient := redis.NewDestinationRedisClient()
	return &Itineraries{
		Logger:          Logger,
		ItinerariesRepo: ItinerariesRepo,
		UserClient:      userClient,
		Redis:           *redisClient,
	}
}

func (i *Itineraries) CreateItineraries(ctx context.Context, in *pb.RequestCreateItineraries) (
	*pb.ResponseCreateItineraries, error) {

	// checking user exists
	valid, err := i.UserClient.ValidateUser(ctx, &pbUser.RequestGetProfile{Id: in.AutherId})
	if err != nil || !valid.Success {
		i.Logger.Error(fmt.Sprintf("error with validating user: %s", err))
		return nil, fmt.Errorf("error: invalid userID: %s", err)
	}

	id, err := i.ItinerariesRepo.CreateItineraries(in)
	if err != nil {
		i.Logger.Error(fmt.Sprintf("error with creating itineraries table: %s", err))
		return nil, err
	}

	err = i.ItinerariesRepo.CreateItinerariesDestinations(id, in.Destinations)
	if err != nil {
		i.Logger.Error(fmt.Sprintf("error with creating itineraries' destnation: %s", err))
		return nil, err
	}

	return &pb.ResponseCreateItineraries{
		Id:          id,
		Title:       in.Title,
		Description: in.Description,
		StartDate:   in.StartDate,
		EndDate:     in.EndDate,
		AuthorId:    in.AutherId,
		CreatedAt:   time.Now().String(),
	}, nil
}

func (i *Itineraries) EditItineraries(ctx context.Context, in *pb.RequestEditItineraries) (
	*pb.ResponseEditItineraries, error) {

	// checking user exists
	valid, err := i.UserClient.ValidateUser(ctx, &pbUser.RequestGetProfile{
		Id: in.AuthorId,
	})
	if err != nil || !valid.Success {
		i.Logger.Error(fmt.Sprintf("error with validating user: %s", err))
		return nil, fmt.Errorf("error: invalid userID: %s", err)
	}

	tx, err := i.ItinerariesRepo.DB.Begin()
	if err != nil {
		i.Logger.Error(fmt.Sprintf("error with creating transaction: %s", err))
		return nil, fmt.Errorf("error with creating transaction: %s", err)

	}
	defer tx.Commit()

	err = postgres.EditItineraries(tx, in)
	if err != nil {
		tx.Rollback()
		i.Logger.Error(fmt.Sprintf("error with editing itineraries table: %s", err))
		return nil, err
	}

	err = postgres.EditItinerariesDestinations(tx, in.Destinations)
	if err != nil {
		tx.Rollback()
		i.Logger.Error(fmt.Sprintf("error with editing itineraries' destnation: %s", err))
		return nil, err
	}

	return &pb.ResponseEditItineraries{
		Id:          in.Id,
		Title:       in.Title,
		Description: in.Description,
		StartDate:   in.StartDate,
		EndDate:     in.EndDate,
		AuthorId:    in.AuthorId,
		UpdatedAt:   time.Now().String(),
	}, nil
}

func (i *Itineraries) GetAllItineraries(ctx context.Context, in *pb.RequestGetAllItineraries) (
	*pb.ResponseGetAllItineraries, error) {
	itineraties, err := i.ItinerariesRepo.GetAllItineraries(in)
	if err != nil {
		i.Logger.Error(fmt.Sprintf("error with getting itineraries: %s", err))
		return nil, err
	}

	resp := pb.ResponseGetAllItineraries{}
	for _, val := range *itineraties {
		auther, err := i.UserClient.GetAuthorInfo(ctx,
			&pbUser.RequestGetAuthorInfo{Id: val.AutherId})
		if err != nil {
			i.Logger.Error(fmt.Sprintf("error with getting author info: %s", err))
			return nil, err
		}
		itiner := pb.Itinerary{
			Id:    val.Id,
			Title: val.Title,
			Auther: &pb.Author{
				Id:       auther.Id,
				Username: auther.Username,
			},
			StartDate:     val.StartDate,
			EndDate:       val.EndDate,
			LikesCount:    int64(val.LikesCount),
			CommentsCount: int64(val.CommentsCount),
			CreatedAt:     val.CreatedAt,
		}

		resp.Itineraries = append(resp.Itineraries, &itiner)
	}

	countOfItiner, err := i.ItinerariesRepo.FindNumberOfItineraries()
	if err != nil {
		i.Logger.Error(fmt.Sprintf("error with getting total itineraries count: %s", err))
		return nil, err
	}
	resp.Total = int64(countOfItiner)
	resp.Limit = in.Limit
	resp.Page = in.Page

	return &resp, nil
}

func (i *Itineraries) GetItineraryFullInfo(ctx context.Context, in *pb.RequestGetItineraryFullInfo) (
	*pb.ResponseGetItineraryFullInfo, error) {

	itineraries, err := i.ItinerariesRepo.GetItinerariesFullInfo(in.Id)
	if err != nil {
		i.Logger.Error(fmt.Sprintf("error with getting full info itineraries: %s", err))
		return nil, err
	}

	auther, err := i.UserClient.GetAuthorInfo(ctx,
		&pbUser.RequestGetAuthorInfo{Id: itineraries.AutherId})
	if err != nil {
		i.Logger.Error(fmt.Sprintf("error with getting author info: %s", err))
		return nil, err
	}

	resp := pb.ResponseGetItineraryFullInfo{
		Id:          itineraries.Id,
		Title:       itineraries.Title,
		Description: itineraries.Description,
		StartDate:   itineraries.StartDate,
		EndDate:     itineraries.EndDate,
		Author: &pb.Author{
			Id:       auther.Id,
			Username: auther.Username,
		},
		LikesCount:    int64(itineraries.LikesCount),
		CommentsCount: int64(itineraries.CommentsCount),
		CreatedAt:     itineraries.CreatedAt,
		UpdatedAt:     itineraries.UpdatedAt,
	}

	destinations, err := i.ItinerariesRepo.GetItinerariesDestinations(resp.Id)
	if err != nil {
		i.Logger.Error(fmt.Sprintf("error with getting itineraries's destination : %s", err))
		return nil, err
	}
	resp.Destinations = *destinations
	return &resp, nil
}

func (i *Itineraries) WriteCommentToItinerary(ctx context.Context, in *pb.RequestWriteCommentToItinerary) (
	*pb.ResponseWriteCommentToItinerary, error) {

	// checking user exists
	valid, err := i.UserClient.ValidateUser(ctx, &pbUser.RequestGetProfile{Id: in.AuthorId})
	if err != nil || !valid.Success {
		i.Logger.Error(fmt.Sprintf("error with validating user: %s", err))
		return nil, fmt.Errorf("error: invalid userID: %s", err)
	}

	id, err := i.ItinerariesRepo.WriteCommentToItinerary(in)
	if err != nil {
		i.Logger.Error(fmt.Sprintf("error with writing comment: %s", err))
		return nil, err
	}
	return &pb.ResponseWriteCommentToItinerary{
		Id:          id,
		Content:     in.Content,
		AuthorId:    in.AuthorId,
		ItineraryId: in.ItineraryId,
		CreatedAt:   time.Now().String(),
	}, nil
}

func (i *Itineraries) CreateDestination(ctx context.Context, in *pb.RequestCreateDestination) (
	*pb.ResponseCreateDestination, error) {
	id, err := i.ItinerariesRepo.CreateDestination(in)
	if err != nil {
		i.Logger.Error(fmt.Sprintf("error with creating destinations: %s", err))
		return nil, err
	}
	return &pb.ResponseCreateDestination{
		Id:                id,
		Name:              in.Name,
		Country:           in.Country,
		Description:       in.Description,
		BestTimeToVisit:   in.BestTimeToVisit,
		AverageCostPerDay: in.AverageCostPerDay,
		Currency:          in.Currency,
		Language:          in.Language,
		PopularityScore:   in.PopularityScore,
		CreatedAt:         time.Now().String(),
	}, nil
}

func (i *Itineraries) UpdateTopDestinations(ctx context.Context, in *pb.RequestGetDestinations) (
	*pb.ResponseGetDestinations, error) {

	destinations, err := i.ItinerariesRepo.GetTopDestinations(in)
	if err != nil {
		i.Logger.Error(fmt.Sprintf("error with getting top destinations: %s", err))
		return nil, err
	}
	err = i.Redis.SetTopDestinations(ctx, destinations)
	return destinations, err
}

func (i *Itineraries) GetDestinations(ctx context.Context, in *pb.RequestGetDestinations) (
	*pb.ResponseGetDestinations, error) {

	destinations, err := i.Redis.GetTopDestinations(ctx)
	if err == rd.Nil {
		destinations, err := i.UpdateTopDestinations(ctx, in)
		if err != nil {
			i.Logger.Error(fmt.Sprintf("error with updating destinations in redis: %s", err))
			return nil, err
		}
		if len(destinations.Destinations) >= int(in.Limit) {
			destinations.Destinations = destinations.Destinations[:in.Limit]
		}
		return destinations, nil
	}
	if err != nil {
		i.Logger.Error(fmt.Sprintf("error with getting destinations: %s", err))
		return nil, err
	}

	if len(destinations.Destinations) < int(in.Limit) {
		destinations, err = i.UpdateTopDestinations(ctx, in)
		if err != nil {
			i.Logger.Error(fmt.Sprintf("error with updating destinations in redis: %s", err))
			return nil, err
		}
	}
	return destinations, nil
}

// func (i *Itineraries) GetDestinationsAllInfo(ctx context.Context, in *pb.RequestGetDestinationsAllInfo) (*pb.ResponseGetDestinationsAllInfo, error)
// func (i *Itineraries) GetMessages(ctx context.Context, in *pb.RequestGetMessages) (*pb.ResponseGetMessages, error)
// func (i *Itineraries) GetUserStatistic(ctx context.Context, in *pb.RequestGetUserStatistic) (*pb.ResponseGetUserStatistic, error)

// func (i *Itineraries) DeleteItineraries(ctx context.Context, in *pb.RequestDeleteItineraries) (
// 	*pb.ResponseDeleteItineraries, error) {
// 		tx, err := i.ItinerariesRepo.DB.Begin()
// 		if err != nil {
// 			i.Logger.Error(fmt.Sprintf("error with creating transaction: %s", err))
// 			return nil, fmt.Errorf("error with creating transaction: %s", err)

// 		}
// 		defer tx.Commit()

// 		err = postgres.EditItineraries(tx, in)
// 		if err != nil {
// 			tx.Rollback()
// 			i.Logger.Error(fmt.Sprintf("error with editing itineraries table: %s", err))
// 			return nil, err
// 		}

// 		err = postgres.EditItinerariesDestinations(tx, in.Destinations)
// 		if err != nil {
// 			tx.Rollback()
// 			i.Logger.Error(fmt.Sprintf("error with editing itineraries' destnation: %s", err))
// 			return nil, err
// 		}

// 		return &pb.ResponseDeleteItineraries{
// 			Message: "Itinerary was deleted successfully",
// 		}, nil
// }
