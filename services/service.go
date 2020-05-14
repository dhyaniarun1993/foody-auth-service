package services

import (
	"context"

	"github.com/dhyaniarun1993/foody-common/errors"
)

// HealthService provides interface for health service
type HealthService interface {
	HealthCheck(context.Context) errors.AppError
}
