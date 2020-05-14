package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/dhyaniarun1993/foody-auth-service/cmd/auth-server/config"
	"github.com/dhyaniarun1993/foody-auth-service/controllers"
	repositories "github.com/dhyaniarun1993/foody-auth-service/repositories/mysql"
	"github.com/dhyaniarun1993/foody-auth-service/services"
	"github.com/dhyaniarun1993/foody-common/datastore/sql"
	"github.com/dhyaniarun1993/foody-common/logger"
	"github.com/dhyaniarun1993/foody-common/tracer"
	"github.com/gorilla/mux"
)

func main() {
	config := config.InitConfiguration()
	logger := logger.CreateLogger(config.Log)
	t, closer := tracer.InitJaeger(config.Jaeger)
	defer closer.Close()

	DB := sql.CreatePool(config.SQL, "mysql", t)

	healthRepository := repositories.NewHealthRepository(DB)

	healthService := services.NewHealthService(healthRepository, logger)

	router := mux.NewRouter()
	ignoredURLs := []string{"/health"}
	ignoredMethods := []string{"OPTION"}

	router.Use(tracer.TraceRequest(t, ignoredURLs, ignoredMethods))
	healthController := controllers.NewHealthController(healthService, logger)

	healthController.LoadRoutes(router)
	serverAddress := ":" + fmt.Sprint(config.Port)
	srv := &http.Server{
		Handler:      router,
		Addr:         serverAddress,
		WriteTimeout: 3 * time.Second,
		ReadTimeout:  3 * time.Second,
	}

	logger.Info("Starting Http server at " + serverAddress)
	serverError := srv.ListenAndServe()
	if serverError != http.ErrServerClosed {
		logger.Error("Http server stopped unexpected")
	} else {
		logger.Info("Http server stopped")
	}
}
