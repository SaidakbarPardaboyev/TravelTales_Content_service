package redis

import (
	"context"
	"encoding/json"
	"time"

	pb "travel/genproto/itineraries"

	"github.com/redis/go-redis/v9"
)

type DestinationRedisClient struct {
	Redis redis.Client
}

func NewDestinationRedisClient() *DestinationRedisClient {
	return &DestinationRedisClient{
		Redis: NewRedicClient(),
	}
}

func (r *DestinationRedisClient) GetTopDestinations(ctx context.Context) (
	*pb.ResponseGetDestinations, error) {
	des, err := r.Redis.Get(ctx, "TopDestinations").Bytes()
	if err == redis.Nil || err != nil {
		return nil, err
	}
	resp := pb.ResponseGetDestinations{}
	err = json.Unmarshal(des, &resp)
	return &resp, err
}

func (r *DestinationRedisClient) SetTopDestinations(ctx context.Context, req *pb.ResponseGetDestinations) error {

	reqMar, err := json.Marshal(req)
	if err != nil {
		return err
	}
	des := r.Redis.Set(ctx, "TopDestinations", string(reqMar), time.Hour)
	if des.Err() != nil {
		return des.Err()
	}
	return nil
}
