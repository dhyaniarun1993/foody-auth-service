package mysql

import (
	"context"
	"net/http"
	"time"

	"github.com/dhyaniarun1993/foody-auth-service/repositories"
	"github.com/dhyaniarun1993/foody-common/datastore/sql"
	"github.com/dhyaniarun1993/foody-common/errors"
)

type healthRepository struct {
	*sql.DB
}

// NewHealthRepository creates and return mysql health repository
func NewHealthRepository(db *sql.DB) repositories.HealthRepository {
	return &healthRepository{db}
}

func (db *healthRepository) HealthCheck(ctx context.Context) errors.AppError {
	timedCtx, pingCancel := context.WithTimeout(ctx, 1*time.Second)
	defer pingCancel()
	pingError := db.PingContext(timedCtx)
	if pingError != nil {
		return errors.NewAppError("Unable to connect to MySQL", http.StatusServiceUnavailable, pingError)
	}
	return nil
}
