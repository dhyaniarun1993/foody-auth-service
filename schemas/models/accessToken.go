package models

import "time"

// AccessToken provides the model definition for Access Token
type AccessToken struct {
	Token      string    `json:"token"`
	ExpiryDate time.Time `json:"expiry_date"`
	ClientID   string    `json:"client_id"`
	UserID     string    `json:"user_id"`
	UserRole   string    `json:"user_role"`
}

// IsActive check if the refresh token is Active(isn't expired)
func (accessToken *AccessToken) IsActive(token string) bool {
	if accessToken.ExpiryDate.After(time.Now()) {
		return true
	}
	return false
}
