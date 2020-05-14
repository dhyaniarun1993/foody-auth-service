package controllers

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/dhyaniarun1993/foody-common/logger"
	"github.com/dhyaniarun1993/foody-customer-service/services"
)

type healthController struct {
	healthService services.HealthService
	logger        *logger.Logger
}

// NewHealthController initialize health endpoint
func NewHealthController(healthService services.HealthService,
	logger *logger.Logger) Controller {
	return &healthController{
		healthService: healthService,
		logger:        logger,
	}
}

func (controller *healthController) LoadRoutes(router *mux.Router) {
	router.HandleFunc("/health", controller.healthCheck).Methods("GET")
}

func (controller *healthController) healthCheck(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	serviceError := controller.healthService.HealthCheck(ctx)
	if serviceError != nil {
		controller.logger.WithError(serviceError).Error("Error occurred")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(serviceError.StatusCode())
		fmt.Fprintf(w, `{"message": %q}`, serviceError.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}
