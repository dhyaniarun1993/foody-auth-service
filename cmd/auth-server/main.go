package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"

	"github.com/dhyaniarun1993/foody-auth-service/cmd/auth-server/config"
	"github.com/dhyaniarun1993/foody-auth-service/controllers"
	httpRepositories "github.com/dhyaniarun1993/foody-auth-service/repositories/http"
	mysqlRepositories "github.com/dhyaniarun1993/foody-auth-service/repositories/mysql"
	redisRepositories "github.com/dhyaniarun1993/foody-auth-service/repositories/redis"
	"github.com/dhyaniarun1993/foody-auth-service/services"
	"github.com/dhyaniarun1993/foody-common/datastore/redis"
	"github.com/dhyaniarun1993/foody-common/datastore/sql"
	"github.com/dhyaniarun1993/foody-common/logger"
	"github.com/dhyaniarun1993/foody-common/tracer"
	"github.com/dhyaniarun1993/foody-common/validator"
	httpCustomerClient "github.com/dhyaniarun1993/foody-customer-service/client/http"
)

func main() {
	config := config.InitConfiguration()
	logger := logger.CreateLogger(config.Log)
	validate := validator.New()
	t, closer := tracer.InitJaeger(config.Jaeger)
	defer closer.Close()

	DB := sql.CreatePool(config.SQL, "mysql", t)
	redisClient := redis.CreateRedisCLient(config.Redis, t)
	customerClient := httpCustomerClient.NewCustomerClient(config.Customer, t)

	mysqlhealthRepository := mysqlRepositories.NewHealthRepository(DB)
	redisHealthRepository := redisRepositories.NewHealthRepository(redisClient)
	clientRepository := mysqlRepositories.NewClientRepository(DB)
	refreshTokenRepository := mysqlRepositories.NewRefreshTokenRepository(DB)
	otpRepository := redisRepositories.NewOtpRepository(redisClient)
	userRepository := httpRepositories.NewUserRepository(customerClient)

	healthService := services.NewHealthService(mysqlhealthRepository, redisHealthRepository, logger)
	otpService := services.NewOtpService(otpRepository)
	tokenService := services.NewTokenservice(config.AccessTokenSecret, config.AccessTokenIssuer, refreshTokenRepository)
	grantService := services.NewGrantService(tokenService, otpService, userRepository)
	authService := services.NewAuthService(grantService, otpService, clientRepository, userRepository)
	userService := services.NewUserService(clientRepository, userRepository, otpService)

	router := mux.NewRouter()
	ignoredURLs := []string{"/health"}
	ignoredMethods := []string{"OPTION"}

	router.Use(tracer.TraceRequest(t, ignoredURLs, ignoredMethods))
	healthController := controllers.NewHealthController(healthService, logger)
	authController := controllers.NewAuthController(authService, validate, logger)
	userController := controllers.NewUserController(userService, validate, logger)

	healthController.LoadRoutes(router)
	authController.LoadRoutes(router)
	userController.LoadRoutes(router)
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
