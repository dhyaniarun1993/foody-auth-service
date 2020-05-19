package dto

import (
	"fmt"
	"net/http"

	"github.com/dhyaniarun1993/foody-common/errors"
	"gopkg.in/go-playground/validator.v9"
)

// SendOtpRequestBody provides the schema definition for send otp api request body
type SendOtpRequestBody struct {
	PhoneNumber string `json:"phone_number" validate:"required,indiaPhoneNumber"`
}

// SendOtpRequest provides the schema definition for send otp api request
type SendOtpRequest struct {
	ClientID     string             `json:"-"`
	ClientSecret string             `json:"-"`
	Body         SendOtpRequestBody `json:"body" validate:"required,dive"`
}

// Validate validates SendOtpRequest
func (dto SendOtpRequest) Validate(validate *validator.Validate) errors.AppError {
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

// TokenRequestBody provides the schema definition for token api request body
type TokenRequestBody struct {
	PhoneNumber  string `json:"phone_number" validate:"indiaPhoneNumber"`
	Otp          int    `json:"otp"`
	RefreshToken string `json:"refresh_token"`
	GrantType    string `json:"grant_type" validate:"required"`
}

// TokenRequest provides the schema definition for token api request
type TokenRequest struct {
	ClientID     string           `json:"-"`
	ClientSecret string           `json:"-"`
	Body         TokenRequestBody `json:"body" validate:"required,dive"`
}

// Validate validates TokenRequest
func (dto TokenRequest) Validate(validate *validator.Validate) errors.AppError {
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

// TokenResponse provides the schema definition for token api response
type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	Type         string `json:"type"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
}
