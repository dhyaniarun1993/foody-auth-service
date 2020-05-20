package repositories

import (
	"context"

	"github.com/dhyaniarun1993/foody-auth-service/schemas/models"
	"github.com/dhyaniarun1993/foody-common/errors"
)

// HealthRepository provides interface for Health repository
type HealthRepository interface {
	HealthCheck(context.Context) errors.AppError
}

// ClientRepository provides interface for Client repository
type ClientRepository interface {
	GetByID(ctx context.Context, clientID string) (models.Client, errors.AppError)
}

// RefreshTokenRepository provides interface from Refresh Token repository
type RefreshTokenRepository interface {
	Create(ctx context.Context, refreshToken models.RefreshToken) (int64, errors.AppError)
	GetByClientIDAndUserID(ctx context.Context, clientID string, userID string) (models.RefreshToken, errors.AppError)
	GetByToken(ctx context.Context, token string) (models.RefreshToken, errors.AppError)
}

// OtpRepository provides interface for Otp repository
type OtpRepository interface {
	Set(ctx context.Context, key string, value string) errors.AppError
	Get(ctx context.Context, key string) (string, errors.AppError)
	Delete(ctx context.Context, key string) errors.AppError
}

// UserRepository provides interface for user repository
type UserRepository interface {
	GetByPhoneNumber(ctx context.Context, phoneNumber string, userType string) (models.User, errors.AppError)
}
