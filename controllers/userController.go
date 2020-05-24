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

type userController struct {
	userService services.UserService
	validate    *validator.Validate
	logger      *logger.Logger
}

// NewUserController initialize auth endpoint
func NewUserController(userService services.UserService, validate *validator.Validate,
	logger *logger.Logger) Controller {
	return &userController{
		userService: userService,
		validate:    validate,
		logger:      logger,
	}
}

func (controller *userController) LoadRoutes(router *mux.Router) {
	router.HandleFunc("/v1/register", controller.register).Methods("POST")
}

func (controller *userController) register(w http.ResponseWriter, r *http.Request) {
	clientID, clientSecret, ok := r.BasicAuth()
	if !ok {
		errMessage := "Basic Auth missing"
		controller.logger.Error(errMessage)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprintf(w, `{"message": %q}`, errMessage)
		return
	}

	var request dto.RegisterRequest
	var requestBody dto.RegisterRequestBody
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

	result, serviceError := controller.userService.Register(ctx, request)
	if serviceError != nil {
		logger.WithError(serviceError).Error("Got Error from Service")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(serviceError.StatusCode())
		fmt.Fprintf(w, `{"message": %q}`, serviceError.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(result)
}
