package models

import (
	"strings"
	"time"

	"github.com/dhyaniarun1993/foody-auth-service/constants"
)

// Client provides the model definition for Client
type Client struct {
	ID                   string    `json:"id"`
	Name                 string    `json:"name"`
	Secret               string    `json:"secret"`
	Type                 string    `json:"type"`
	GrantType            string    `json:"grant_type"`
	UserRole             string    `json:"user_role"`
	AccessTokenLifetime  int       `json:"access_token_lifetime"`
	RefreshTokenLifetime int       `json:"refresh_token_lifetime"`
	Status               string    `json:"status"`
	CreatedAt            time.Time `json:"created_at"`
	UpdatedAt            time.Time `json:"updated_at"`
}

// IsActive check if client status is active
func (client *Client) IsActive() bool {
	if client.Status == constants.ClientStatusActive {
		return true
	}
	return false
}

// IsValidSecret validates the secret provided
func (client *Client) IsValidSecret(secret string) bool {
	if client.Secret == secret {
		return true
	}
	return false
}

// IsValidGrantType checks if grant type provided is valid for the client
func (client *Client) IsValidGrantType(grantType string) bool {
	grantTypes := strings.Split(client.GrantType, " ")
	for _, item := range grantTypes {
		if item == grantType {
			return true
		}
	}
	return false
}
