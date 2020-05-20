package redis

import (
	"context"
	"net/http"
	"strconv"
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

func (redis *otpRepository) Set(ctx context.Context, key string, value int) errors.AppError {

	ttl := 120 * time.Second
	redisWithContext := otredis.WrapRedisClient(ctx, redis.Client)

	err := redisWithContext.Set(key, value, ttl).Err()
	if err != nil {
		return errors.NewAppError("Something went wrong", http.StatusInternalServerError, err)
	}
	return nil
}

func (redis *otpRepository) Get(ctx context.Context, key string) (int, errors.AppError) {
	redisWithContext := otredis.WrapRedisClient(ctx, redis.Client)

	otpString, err := redisWithContext.Get(key).Result()
	if err != nil {
		return 0, errors.NewAppError("Something went wrong", http.StatusInternalServerError, err)
	}

	otp, conversionErr := strconv.Atoi(otpString)
	if conversionErr != nil {
		return 0, errors.NewAppError("Something went wrong", http.StatusInternalServerError, err)
	}

	return otp, nil
}

func (redis *otpRepository) Delete(ctx context.Context, key string) errors.AppError {
	redisWithContext := otredis.WrapRedisClient(ctx, redis.Client)

	err := redisWithContext.Del(key).Err()
	if err != nil {
		return errors.NewAppError("Something went wrong", http.StatusInternalServerError, err)
	}
	return nil
}
