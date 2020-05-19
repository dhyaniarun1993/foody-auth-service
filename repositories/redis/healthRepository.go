package redis

import (
	"context"
	"net/http"

	"github.com/go-redis/redis"

	"github.com/dhyaniarun1993/foody-auth-service/repositories"
	"github.com/dhyaniarun1993/foody-common/errors"
)

type healthRepository struct {
	*redis.Client
}

// NewHealthRepository creates and return mysql health repository
func NewHealthRepository(client *redis.Client) repositories.HealthRepository {
	return &healthRepository{client}
}

func (redis *healthRepository) HealthCheck(ctx context.Context) errors.AppError {

	_, pingError := redis.Ping().Result()
	if pingError != nil {
		return errors.NewAppError("Unable to connect to Redis", http.StatusServiceUnavailable, pingError)
	}
	return nil
}
