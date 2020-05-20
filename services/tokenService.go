package services

import (
	"context"
	"encoding/base64"
	"fmt"
	"net/http"
	"reflect"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/dhyaniarun1993/foody-auth-service/repositories"
	"github.com/dhyaniarun1993/foody-auth-service/schemas/models"
	"github.com/dhyaniarun1993/foody-common/errors"
	"golang.org/x/crypto/bcrypt"
)

type tokenService struct {
	accessTokenSecret      string
	accessTokenIssuer      string
	refreshTokenRepository repositories.RefreshTokenRepository
}

// AccessTokenClaims Struct that will encoded to jwt
type AccessTokenClaims struct {
	UserID   string `json:"user_id"`
	UserRole string `json:"user_role"`
	ClientID string `json:"client_id"`
	jwt.StandardClaims
}

// NewTokenservice creates and return token service object
func NewTokenservice(accessTokenSecret string, accessTokenIssuer string,
	refreshTokenRepository repositories.RefreshTokenRepository) TokenService {
	return &tokenService{accessTokenSecret, accessTokenIssuer, refreshTokenRepository}
}

// GenerateAccessToken generates jwt Access token
func (service *tokenService) GenerateAccessToken(client models.Client, user models.User) (models.AccessToken, errors.AppError) {

	expirationTime := time.Now().Add(time.Duration(client.AccessTokenLifetime) * time.Second)
	claims := &AccessTokenClaims{
		UserID:   user.ID,
		UserRole: client.UserRole,
		ClientID: client.ID,
		StandardClaims: jwt.StandardClaims{
			Subject:   client.UserRole + "|" + user.ID,
			Issuer:    service.accessTokenIssuer,
			IssuedAt:  time.Now().Unix(),
			ExpiresAt: expirationTime.Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(service.accessTokenSecret))
	accessToken := models.AccessToken{
		Token:      tokenString,
		ClientID:   client.ID,
		UserID:     user.ID,
		UserRole:   client.UserRole,
		ExpiryDate: expirationTime,
	}
	if err != nil {
		return models.AccessToken{}, errors.NewAppError("Something went wrong", http.StatusInternalServerError, err)
	}
	return accessToken, nil
}

func (service *tokenService) VerifyAccessToken(tokenString string) (models.AccessToken, errors.AppError) {

	claims := AccessTokenClaims{}

	token, err := jwt.ParseWithClaims(tokenString, &claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return service.accessTokenSecret, nil
	})

	if err != nil {
		return models.AccessToken{}, errors.NewAppError("Unauthorized", http.StatusUnauthorized, err)
	}

	if !token.Valid {
		return models.AccessToken{}, errors.NewAppError("Unauthorized", http.StatusUnauthorized, nil)
	}
	accessToken := models.AccessToken{
		Token:    tokenString,
		ClientID: claims.ClientID,
		UserID:   claims.UserID,
		UserRole: claims.UserRole,
	}
	return accessToken, nil
}

func (service *tokenService) GenerateRefreshToken(ctx context.Context,
	client models.Client, user models.User) (models.RefreshToken, errors.AppError) {

	var refreshToken models.RefreshToken

	// check if get exist in DB for provided client ID and User ID
	refreshToken, getTokenError := service.refreshTokenRepository.GetByClientIDAndUserID(ctx, client.ID, user.ID)
	if getTokenError != nil {
		return refreshToken, errors.NewAppError("Something went worng", http.StatusInternalServerError, getTokenError)
	}

	// check if token is active
	if !reflect.DeepEqual(refreshToken, models.RefreshToken{}) &&
		refreshToken.ExpiryDate.After(time.Now()) {

		return refreshToken, nil
	}

	// use role and user id to generate hash
	hashInput := client.UserRole + "-" + user.ID
	hash, err := bcrypt.GenerateFromPassword([]byte(hashInput), bcrypt.DefaultCost)
	if err != nil {
		return refreshToken, errors.NewAppError("Something went worng", http.StatusInternalServerError, err)
	}
	token := base64.StdEncoding.EncodeToString(hash)
	refreshToken = models.RefreshToken{
		Token:      token,
		ExpiryDate: time.Now().Add(time.Duration(client.RefreshTokenLifetime) * time.Second),
		ClientID:   client.ID,
		UserID:     user.ID,
		UserRole:   client.UserRole,
	}

	// save the token the in the database
	refreshTokenID, insertTokenErr := service.refreshTokenRepository.Create(ctx, refreshToken)
	if insertTokenErr != nil {
		return refreshToken, errors.NewAppError("Something went wrong", http.StatusInternalServerError, insertTokenErr)
	}

	refreshToken.ID = refreshTokenID
	return refreshToken, nil
}

func (service *tokenService) VerifyRefreshToken(ctx context.Context,
	tokenString string) (models.RefreshToken, errors.AppError) {

	refreshToken, getTokenErr := service.refreshTokenRepository.GetByToken(ctx, tokenString)
	if getTokenErr != nil {
		return models.RefreshToken{}, getTokenErr
	}

	if reflect.DeepEqual(refreshToken, models.RefreshToken{}) {
		return models.RefreshToken{}, errors.NewAppError("Invalid token", http.StatusUnauthorized, nil)
	}

	// check if token is expired
	if !refreshToken.ExpiryDate.After(time.Now()) {
		return models.RefreshToken{}, errors.NewAppError("Token expired", http.StatusUnauthorized, nil)
	}

	return refreshToken, nil
}
