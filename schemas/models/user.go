package models

import "github.com/dhyaniarun1993/foody-auth-service/constants"

// User provides the model definition for User
type User struct {
	ID     string `json:"id"`
	Status string `json:"status"`
}

// IsActive check if user status is active
func (user *User) IsActive() bool {
	if user.Status == constants.UserStatusActive {
		return true
	}
	return false
}
