package services

import (
	"context"

	"github.com/dhyaniarun1993/foody-auth-service/repositories"
	"github.com/dhyaniarun1993/foody-common/errors"
	"github.com/dhyaniarun1993/foody-common/logger"
)

type healthService struct {
	mysqlHealthRepository repositories.HealthRepository
	redisHealthRepository repositories.HealthRepository
	logger                *logger.Logger
}

// NewHealthService creates and return health service object
func NewHealthService(mysqlHealthRepository repositories.HealthRepository,
	redisHealthRepository repositories.HealthRepository, logger *logger.Logger) HealthService {
	return &healthService{mysqlHealthRepository, redisHealthRepository, logger}
}

func (service *healthService) HealthCheck(ctx context.Context) errors.AppError {
	mysqlRepositoryError := service.mysqlHealthRepository.HealthCheck(ctx)
	if mysqlRepositoryError != nil {
		return mysqlRepositoryError
	}
	redisRepositoryError := service.redisHealthRepository.HealthCheck(ctx)
	if redisRepositoryError != nil {
		return redisRepositoryError
	}
	return nil
}
