package dto

import (
	"fmt"
	"net/http"

	"github.com/dhyaniarun1993/foody-common/errors"
	"gopkg.in/go-playground/validator.v9"
)

// RegisterRequestBody provides the schema definition for register api request body
type RegisterRequestBody struct {
	PhoneNumber string `json:"phone_number" validate:"required,indiaPhoneNumber"`
	Email       string `json:"email" validate:"required,email"`
	FirstName   string `json:"first_name" validate:"required"`
	LastName    string `json:"last_name" validate:"required"`
}

// RegisterRequest provides the schema definition for register api request
type RegisterRequest struct {
	ClientID     string              `json:"-"`
	ClientSecret string              `json:"-"`
	Body         RegisterRequestBody `json:"body" validate:"required,dive"`
}

// Validate validates SendOtpRequest
func (dto RegisterRequest) Validate(validate *validator.Validate) errors.AppError {
	var errMsg string
	err := validate.Struct(dto)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			errMsg = fmt.Sprintf("Invalid value for field '%s'", err.Field())
			break
		}
		return errors.NewAppError(errMsg, http.StatusBadRequest, err)
	}

	return nil
}

// RegisterResponse provides the schema definition for register api response
type RegisterResponse struct {
	ID int64 `json:"id"`
}
