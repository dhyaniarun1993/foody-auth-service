package redis

import (
	"context"
	"net/http"
	"time"

	"github.com/go-redis/redis"
	otredis "github.com/smacker/opentracing-go-redis"

	"github.com/dhyaniarun1993/foody-auth-service/repositories"
	"github.com/dhyaniarun1993/foody-common/errors"
)

type otpRepository struct {
	*redis.Client
}

// NewOtpRepository creates and return redis otp repository
func NewOtpRepository(redisClient *redis.Client) repositories.OtpRepository {
	return &otpRepository{redisClient}
}

func (db *otpRepository) Set(ctx context.Context, key string, value string) errors.AppError {

	ttl := 120 * time.Second
	redisWithContext := otredis.WrapRedisClient(ctx, db.Client)

	err := redisWithContext.Set(key, value, ttl).Err()
	if err != nil {
		return errors.NewAppError("Something went wrong", http.StatusInternalServerError, err)
	}
	return nil
}

func (db *otpRepository) Get(ctx context.Context, key string) (string, errors.AppError) {
	redisWithContext := otredis.WrapRedisClient(ctx, db.Client)

	otp, err := redisWithContext.Get(key).Result()
	if err != nil && err != redis.Nil {
		return "", errors.NewAppError("Something went wrong", http.StatusInternalServerError, err)
	}

	return otp, nil
}

func (db *otpRepository) Delete(ctx context.Context, key string) errors.AppError {
	redisWithContext := otredis.WrapRedisClient(ctx, db.Client)

	err := redisWithContext.Del(key).Err()
	if err != nil {
		return errors.NewAppError("Something went wrong", http.StatusInternalServerError, err)
	}
	return nil
}
