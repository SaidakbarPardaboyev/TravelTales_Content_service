package redis

import (
	"context"
	"testing"

	"github.com/redis/go-redis/v9"
)

func NewRedis() *DestinationRedisClient {
	return NewDestinationRedisClient()
}
func TestGetTopDestinations(t *testing.T) {
	_, err := NewRedis().GetTopDestinations(context.Background())
	if err == redis.Nil || err != nil {
		t.Error(err)
	}
}
