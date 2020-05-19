package services

import (
	"context"

	"github.com/dhyaniarun1993/foody-auth-service/schemas/dto"
	"github.com/dhyaniarun1993/foody-auth-service/schemas/models"
	"github.com/dhyaniarun1993/foody-common/errors"
)

// HealthService provides interface for health service
type HealthService interface {
	HealthCheck(context.Context) errors.AppError
}

// AuthService provides interface for auth service
type AuthService interface {
	SendOtp(ctx context.Context, request dto.SendOtpRequest) errors.AppError
	GetToken(ctx context.Context, request dto.TokenRequest) (dto.TokenResponse, errors.AppError)
}

// TokenService provides interface for token service
type TokenService interface {
	GenerateAccessToken(client models.Client, user models.User) (models.AccessToken, errors.AppError)
	VerifyAccessToken(tokenString string) (models.AccessToken, errors.AppError)
	GenerateRefreshToken(ctx context.Context, client models.Client, user models.User) (models.RefreshToken, errors.AppError)
	VerifyRefreshToken(ctx context.Context, tokenString string) (models.RefreshToken, errors.AppError)
}

// GrantService provides interface for grant service
type GrantService interface {
	HandleOTPGrant(ctx context.Context, phoneNumber string, otp int,
		client models.Client) (models.AccessToken, models.RefreshToken, errors.AppError)
	HandleRefreshTokenGrant(ctx context.Context, tokenString string,
		client models.Client) (models.AccessToken, models.RefreshToken, errors.AppError)
}

// OtpService provides interface for otp srevice
type OtpService interface {
	Validate(ctx context.Context, phoneNumber string, otp int, client models.Client) (bool, errors.AppError)
	Generate(ctx context.Context, phoneNumber string, client models.Client) errors.AppError
}
