package http

import (
	"context"
	"net/http"

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

	user := models.User{
		ID:     string(customer.ID),
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