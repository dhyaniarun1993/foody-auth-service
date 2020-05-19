package models

import "time"

// RefreshToken provides the model definition for Refresh Token
type RefreshToken struct {
	ID         int64     `json:"id"`
	Token      string    `json:"token"`
	ExpiryDate time.Time `json:"expiry_date"`
	ClientID   string    `json:"client_id"`
	UserID     string    `json:"user_id"`
	UserRole   string    `json:"user_role"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

// IsActive check if the refresh token is Active(isn't expired)
func (refreshToken *RefreshToken) IsActive(token string) bool {
	if refreshToken.ExpiryDate.After(time.Now()) {
		return true
	}
	return false
}
