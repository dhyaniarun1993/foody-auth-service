package http

import (
	"context"
	"net/http"
	"strconv"

	"github.com/dhyaniarun1993/foody-auth-service/constants"
	"github.com/dhyaniarun1993/foody-auth-service/repositories"
	"github.com/dhyaniarun1993/foody-auth-service/schemas/models"
	"github.com/dhyaniarun1993/foody-common/errors"
	customer "github.com/dhyaniarun1993/foody-customer-service/client"
	customerDto "github.com/dhyaniarun1993/foody-customer-service/schemas/dto"
)

type userRepository struct {
	customerClient customer.Client
}

// NewUserRepository creates and return http user repository
func NewUserRepository(customerClient customer.Client) repositories.UserRepository {
	return &userRepository{customerClient}
}

func (repository *userRepository) getCustomerByPhoneNumber(ctx context.Context,
	phoneNumber string) (models.User, errors.AppError) {

	query := customerDto.GetCustomerRequestQuery{
		PhoneNumber: phoneNumber,
	}
	customer, err := repository.customerClient.InternalGetCustomer(ctx, query)
	if err != nil {
		return models.User{}, err
	}

	userID := strconv.Itoa(int(customer.ID))
	user := models.User{
		ID:     userID,
		Status: customer.Status,
	}
	return user, nil
}

func (repository *userRepository) GetByPhoneNumber(ctx context.Context, phoneNumber string,
	userType string) (models.User, errors.AppError) {

	switch userType {
	case constants.UserTypeCustomer:
		user, err := repository.getCustomerByPhoneNumber(ctx, phoneNumber)
		return user, err
	}
	return models.User{}, errors.NewAppError("Invalid user type", http.StatusInternalServerError, nil)
}

func (repository *userRepository) CreateCustomer(ctx context.Context, phoneNumber, email,
	firstName, lastName string) (int64, errors.AppError) {

	body := customerDto.CreateCustomerRequestBody{
		PhoneNumber: phoneNumber,
		Email:       email,
		FirstName:   firstName,
		LastName:    lastName,
	}
	response, err := repository.customerClient.InternalCreateCustomer(ctx, body)
	return response.ID, err
}
