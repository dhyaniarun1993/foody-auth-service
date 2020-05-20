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

type clientRepository struct {
	*instrumentedSQL.DB
}

// NewClientRepository creates and return mysql client repository
func NewClientRepository(db *instrumentedSQL.DB) repositories.ClientRepository {
	return &clientRepository{db}
}

func (db *clientRepository) GetByID(ctx context.Context, clientID string) (models.Client, errors.AppError) {
	var client models.Client
	timedCtx, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()

	query := "SELECT id, name, secret, type, grant_type, user_role, access_token_lifetime, refresh_token_lifetime, status, created_at, updated_at FROM client WHERE id = ?"
	row := db.QueryRowContext(timedCtx, query, clientID)

	err := row.Scan(&client.ID, &client.Name, &client.Secret, &client.Type, &client.GrantType, &client.UserRole, &client.AccessTokenLifetime, &client.RefreshTokenLifetime, &client.Status, &client.CreatedAt, &client.UpdatedAt)
	if err != nil && err != sql.ErrNoRows {
		return client, errors.NewAppError("Something went wrong", http.StatusInternalServerError, err)
	}
	return client, nil
}
