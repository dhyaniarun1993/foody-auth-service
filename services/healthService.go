package services

import (
	"context"

	"github.com/dhyaniarun1993/foody-common/errors"
	"github.com/dhyaniarun1993/foody-common/logger"
	"github.com/dhyaniarun1993/foody-customer-service/repositories"
)

type healthService struct {
	healthRepository repositories.HealthRepository
	logger           *logger.Logger
}

// NewHealthService creates and return health service object
func NewHealthService(healthRepository repositories.HealthRepository,
	logger *logger.Logger) HealthService {
	return &healthService{healthRepository, logger}
}

func (service *healthService) HealthCheck(ctx context.Context) errors.AppError {
	repositoryError := service.healthRepository.HealthCheck(ctx)
	return repositoryError
}
