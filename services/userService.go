package services

import (
	"context"
	"net/http"

	"github.com/dhyaniarun1993/foody-auth-service/constants"
	"github.com/dhyaniarun1993/foody-auth-service/repositories"
	"github.com/dhyaniarun1993/foody-auth-service/schemas/dto"
	"github.com/dhyaniarun1993/foody-auth-service/schemas/models"
	"github.com/dhyaniarun1993/foody-common/errors"
)

type userService struct {
	clientRepository repositories.ClientRepository
	userRepository   repositories.UserRepository
	otpService       OtpService
}

// NewUserService creates and return user service object
func NewUserService(clientRepository repositories.ClientRepository,
	userRepository repositories.UserRepository, otpService OtpService) UserService {
	return &userService{clientRepository, userRepository, otpService}
}

func (service *userService) validateAndGetClient(ctx context.Context, clientID string,
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

func (service *userService) Register(ctx context.Context, request dto.RegisterRequest) (dto.RegisterResponse, errors.AppError) {

	client, validateClientErr := service.validateAndGetClient(ctx, request.ClientID, request.ClientSecret)
	if validateClientErr != nil {
		return dto.RegisterResponse{}, validateClientErr
	}

	if client.UserRole != constants.UserTypeCustomer {
		return dto.RegisterResponse{}, errors.NewAppError("client not authorized to register user", http.StatusForbidden, nil)
	}

	// Todo: Use Otp verification before creating the customer
	userID, createCustomerErr := service.userRepository.CreateCustomer(ctx, request.Body.PhoneNumber, request.Body.Email,
		request.Body.FirstName, request.Body.LastName)
	if createCustomerErr != nil {
		return dto.RegisterResponse{}, createCustomerErr
	}

	generateOtpErr := service.otpService.Generate(ctx, request.Body.PhoneNumber, client)
	if generateOtpErr != nil {
		return dto.RegisterResponse{}, generateOtpErr
	}
	response := dto.RegisterResponse{
		ID: userID,
	}
	return response, nil
}
