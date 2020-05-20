package services

import (
	"context"
	"net/http"
	"reflect"

	"github.com/dhyaniarun1993/foody-auth-service/repositories"
	"github.com/dhyaniarun1993/foody-auth-service/schemas/models"
	"github.com/dhyaniarun1993/foody-common/errors"
)

// constants for Grant type
const (
	OtpGrant          = "otp"
	RefreshTokenGrant = "refresh_token"
)

type grantService struct {
	tokenService   TokenService
	otpService     OtpService
	userRepository repositories.UserRepository
}

// NewGrantService creates and return grant service object.
func NewGrantService(tokenService TokenService, otpService OtpService,
	userRepository repositories.UserRepository) GrantService {
	return &grantService{tokenService, otpService, userRepository}
}

func (service *grantService) HandleOTPGrant(ctx context.Context, phoneNumber string, otp string,
	client models.Client) (models.AccessToken, models.RefreshToken, errors.AppError) {

	isValidOtp, otpValidateErr := service.otpService.Validate(ctx, phoneNumber, otp, client)
	if otpValidateErr != nil {
		return models.AccessToken{}, models.RefreshToken{}, otpValidateErr
	}

	if !isValidOtp {
		return models.AccessToken{}, models.RefreshToken{}, errors.NewAppError("Invalid Otp", http.StatusUnauthorized, nil)
	}

	user, userErr := service.userRepository.GetByPhoneNumber(ctx, phoneNumber, client.UserRole)
	if reflect.DeepEqual(user, models.User{}) {
		return models.AccessToken{}, models.RefreshToken{}, errors.NewAppError("Unable to find user", http.StatusUnauthorized, nil)
	}
	if userErr != nil {
		return models.AccessToken{}, models.RefreshToken{}, userErr
	}

	accessToken, accessTokenErr := service.tokenService.GenerateAccessToken(client, user)
	if accessTokenErr != nil {
		return models.AccessToken{}, models.RefreshToken{}, accessTokenErr
	}

	refreshToken, refreshTokenErr := service.tokenService.GenerateRefreshToken(ctx, client, user)
	if refreshTokenErr != nil {
		return models.AccessToken{}, models.RefreshToken{}, refreshTokenErr
	}

	return accessToken, refreshToken, nil
}

func (service *grantService) HandleRefreshTokenGrant(ctx context.Context, tokenString string,
	client models.Client) (models.AccessToken, models.RefreshToken, errors.AppError) {

	refreshToken, verifyTokenError := service.tokenService.VerifyRefreshToken(ctx, tokenString)
	if verifyTokenError != nil {
		return models.AccessToken{}, models.RefreshToken{}, verifyTokenError
	}

	if refreshToken.ClientID != client.ID {
		return models.AccessToken{}, models.RefreshToken{}, errors.NewAppError("Invalid token", http.StatusUnauthorized, nil)
	}

	user := models.User{
		ID: refreshToken.UserID,
	}

	accessToken, accessTokenErr := service.tokenService.GenerateAccessToken(client, user)
	if accessTokenErr != nil {
		return models.AccessToken{}, models.RefreshToken{}, accessTokenErr
	}

	return accessToken, refreshToken, nil
}
