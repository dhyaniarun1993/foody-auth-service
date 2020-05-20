package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"gopkg.in/go-playground/validator.v9"

	"github.com/dhyaniarun1993/foody-auth-service/schemas/dto"
	"github.com/dhyaniarun1993/foody-auth-service/services"
	"github.com/dhyaniarun1993/foody-common/logger"
)

type authController struct {
	authService services.AuthService
	validate    *validator.Validate
	logger      *logger.Logger
}

// NewAuthController initialize auth endpoint
func NewAuthController(authService services.AuthService, validate *validator.Validate,
	logger *logger.Logger) Controller {
	return &authController{
		authService: authService,
		validate:    validate,
		logger:      logger,
	}
}

func (controller *authController) LoadRoutes(router *mux.Router) {
	router.HandleFunc("/v1/auth/send-otp", controller.sendOtp).Methods("POST")
	router.HandleFunc("/v1/auth/token", controller.token).Methods("POST")
}

func (controller *authController) sendOtp(w http.ResponseWriter, r *http.Request) {
	clientID, clientSecret, ok := r.BasicAuth()
	if !ok {
		errMessage := "Basic Auth missing"
		controller.logger.Error(errMessage)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprintf(w, `{"message": %q}`, errMessage)
		return
	}

	var request dto.SendOtpRequest
	var requestBody dto.SendOtpRequestBody
	ctx := r.Context()
	request.ClientID = clientID
	request.ClientSecret = clientSecret

	logger := controller.logger.WithContext(ctx)
	decodingError := json.NewDecoder(r.Body).Decode(&requestBody)
	if decodingError != nil {
		errorMsg := "Invalid request"
		logger.WithError(decodingError).Error(errorMsg)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, `{"message": %q}`, errorMsg)
		return
	}

	request.Body = requestBody
	validationError := request.Validate(controller.validate)
	if validationError != nil {
		logger.WithError(validationError).Error("Invalid request body")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(validationError.StatusCode())
		fmt.Fprintf(w, `{"message": %q}`, validationError.Error())
		return
	}

	serviceError := controller.authService.SendOtp(ctx, request)
	if serviceError != nil {
		logger.WithError(serviceError).Error("Got Error from Service")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(serviceError.StatusCode())
		fmt.Fprintf(w, `{"message": %q}`, serviceError.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
}

func (controller *authController) token(w http.ResponseWriter, r *http.Request) {
	clientID, clientSecret, ok := r.BasicAuth()
	if !ok {
		errMessage := "Basic Auth missing"
		controller.logger.Error(errMessage)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprintf(w, `{"message": %q}`, errMessage)
		return
	}

	var request dto.TokenRequest
	var requestBody dto.TokenRequestBody
	ctx := r.Context()
	request.ClientID = clientID
	request.ClientSecret = clientSecret

	logger := controller.logger.WithContext(ctx)
	decodingError := json.NewDecoder(r.Body).Decode(&requestBody)
	if decodingError != nil {
		errorMsg := "Invalid request"
		logger.WithError(decodingError).Error(errorMsg)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, `{"message": %q}`, errorMsg)
		return
	}

	request.Body = requestBody
	validationError := request.Validate(controller.validate)
	if validationError != nil {
		logger.WithError(validationError).Error("Invalid request body")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(validationError.StatusCode())
		fmt.Fprintf(w, `{"message": %q}`, validationError.Error())
		return
	}

	result, serviceError := controller.authService.GetToken(ctx, request)
	if serviceError != nil {
		logger.WithError(serviceError).Error("Got Error from Service")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(serviceError.StatusCode())
		fmt.Fprintf(w, `{"message": %q}`, serviceError.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(result)
}
