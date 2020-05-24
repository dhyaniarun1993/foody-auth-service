package services

import (
	"context"
	"net/http"
	"reflect"

	"github.com/dhyaniarun1993/foody-auth-service/constants"
	"github.com/dhyaniarun1993/foody-auth-service/repositories"
	"github.com/dhyaniarun1993/foody-auth-service/schemas/dto"
	"github.com/dhyaniarun1993/foody-auth-service/schemas/models"
	"github.com/dhyaniarun1993/foody-common/errors"
)

type authService struct {
	grantService     GrantService
	otpService       OtpService
	clientRepository repositories.ClientRepository
	userRepository   repositories.UserRepository
}

// NewAuthService creates and return auth service object
func NewAuthService(grantService GrantService, otpService OtpService,
	clientRepository repositories.ClientRepository, userRepository repositories.UserRepository) AuthService {
	return &authService{grantService, otpService, clientRepository, userRepository}
}

func (service *authService) validateAndGetClient(ctx context.Context, clientID string,
	clientSecret string) (models.Client, errors.AppError) {
	client, getClientErr := service.clientRepository.GetByID(ctx, clientID)
	if getClientErr != nil {
		return models.Client{}, getClientErr
	}

	if !client.IsActive() {
		return models.Client{}, errors.NewAppError("Invalid client", http.StatusUnauthorized, nil)
	}

	// Todo: Store hash secret in the DB
	if !client.IsValidSecret(clientSecret) {
		return models.Client{}, errors.NewAppError("Unauthorized client", http.StatusUnauthorized, nil)
	}

	return client, nil
}

func (service *authService) SendOtp(ctx context.Context, request dto.SendOtpRequest) errors.AppError {

	client, validateClientErr := service.validateAndGetClient(ctx, request.ClientID, request.ClientSecret)
	if validateClientErr != nil {
		return validateClientErr
	}

	if !client.IsValidGrantType(constants.GrantTypeOtp) {
		return errors.NewAppError("Client not authorized to send otp", http.StatusForbidden, nil)
	}

	user, getUserErr := service.userRepository.GetByPhoneNumber(ctx, request.Body.PhoneNumber, client.UserRole)
	if getUserErr != nil {
		return getUserErr
	}
	if reflect.DeepEqual(user, models.User{}) {
		return errors.NewAppError("Phone Number not registered", http.StatusUnprocessableEntity, nil)
	}
	if !user.IsActive() {
		return errors.NewAppError("Account is inactive", http.StatusUnprocessableEntity, nil)
	}
	generateOtpErr := service.otpService.Generate(ctx, request.Body.PhoneNumber, client)
	if generateOtpErr != nil {
		return generateOtpErr
	}
	return nil
}

func (service *authService) GetToken(ctx context.Context, request dto.TokenRequest) (dto.TokenResponse, errors.AppError) {

	var accessToken models.AccessToken
	var refreshToken models.RefreshToken
	var handleGrantErr errors.AppError
	client, validateClientErr := service.validateAndGetClient(ctx, request.ClientID, request.ClientSecret)
	if validateClientErr != nil {
		return dto.TokenResponse{}, validateClientErr
	}

	if !client.IsValidGrantType(request.Body.GrantType) {
		return dto.TokenResponse{}, errors.NewAppError("Invalid grant type for client", http.StatusBadRequest, nil)
	}

	switch request.Body.GrantType {
	case "otp":
		accessToken, refreshToken, handleGrantErr = service.grantService.HandleOTPGrant(ctx,
			request.Body.PhoneNumber, request.Body.Otp, client)
	case "refresh_token":
		accessToken, refreshToken, handleGrantErr = service.grantService.HandleRefreshTokenGrant(ctx,
			request.Body.RefreshToken, client)
	}

	if handleGrantErr != nil {
		return dto.TokenResponse{}, handleGrantErr
	}

	token := dto.TokenResponse{
		AccessToken:  accessToken.Token,
		Type:         "Bearer",
		ExpiresIn:    client.AccessTokenLifetime,
		RefreshToken: refreshToken.Token,
	}

	return token, nil
}
