package mysql

import (
	"context"
	"database/sql"
	"net/http"
	"time"

	"github.com/dhyaniarun1993/foody-auth-service/repositories"
	"github.com/dhyaniarun1993/foody-auth-service/schemas/models"
	instrumentedSQL "github.com/dhyaniarun1993/foody-common/datastore/sql"
	"github.com/dhyaniarun1993/foody-common/errors"
)

type refreshTokenRepository struct {
	*instrumentedSQL.DB
}

// NewRefreshTokenRepository creates and return mysql refresh token repository
func NewRefreshTokenRepository(db *instrumentedSQL.DB) repositories.RefreshTokenRepository {
	return &refreshTokenRepository{db}
}

func (db *refreshTokenRepository) Create(ctx context.Context, refreshToken models.RefreshToken) (int64, errors.AppError) {
	timedCtx, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()

	query := "INSERT INTO refresh_token SET token = ?, expiry_date = ?, client_id = ?, user_id = ?, user_role = ?"
	res, insertErr := db.ExecContext(timedCtx, query, refreshToken.Token, refreshToken.ExpiryDate, refreshToken.ClientID, refreshToken.UserID, refreshToken.UserRole)
	if insertErr != nil {
		return 0, errors.NewAppError("Something went wrong", http.StatusInternalServerError, insertErr)
	}

	tokenID, fetchIDError := res.LastInsertId()
	if fetchIDError != nil {
		return 0, errors.NewAppError("Unable to fetch Id", http.StatusInternalServerError, fetchIDError)
	}
	return tokenID, nil
}

func (db *refreshTokenRepository) GetByClientIDAndUserID(ctx context.Context,
	clientID string, userID string) (models.RefreshToken, errors.AppError) {

	var refreshToken models.RefreshToken
	timedCtx, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()

	query := "SELECT id, token, expiry_date, client_id, user_id, user_role, created_at, updated_at FROM refresh_token WHERE client_id = ? and user_id = ?"
	row := db.QueryRowContext(timedCtx, query, clientID, userID)
	err := row.Scan(&refreshToken.ID, &refreshToken.Token, &refreshToken.ExpiryDate, &refreshToken.ClientID, &refreshToken.UserID, &refreshToken.UserRole, &refreshToken.CreatedAt, &refreshToken.UpdatedAt)
	if err != nil && err != sql.ErrNoRows {
		return refreshToken, errors.NewAppError("Something went wrong", http.StatusInternalServerError, err)
	}
	return refreshToken, nil
}

func (db *refreshTokenRepository) GetByToken(ctx context.Context, token string) (models.RefreshToken, errors.AppError) {
	var refreshToken models.RefreshToken
	timedCtx, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()

	query := "SELECT id, token, expiry_date, client_id, user_id, user_role, created_at, updated_at FROM refresh_token WHERE token = ?"
	row := db.QueryRowContext(timedCtx, query, token)
	err := row.Scan(&refreshToken.ID, &refreshToken.Token, &refreshToken.ExpiryDate, &refreshToken.ClientID, &refreshToken.UserID, &refreshToken.UserRole, &refreshToken.CreatedAt, &refreshToken.UpdatedAt)
	if err != nil && err != sql.ErrNoRows {
		return refreshToken, errors.NewAppError("Something went wrong", http.StatusInternalServerError, err)
	}
	return refreshToken, nil
}
