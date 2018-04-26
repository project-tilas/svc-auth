package domain

import (
	"time"
)

// User represents a single user in the tilas platform
type User struct {
	ID        string    `json:"id"`
	Username  string    `json:"username"`
	Password  string    `json:"password,omitempty"`
	CreatedAt time.Time `json:"createdAt,omitempty"`
}

// Validate tries to validate the struct, if it does not validate correctly it give back an error.
func (u *User) Validate() error {
	if u.Password == "" {
		return ErrPasswordRequired
	}

	if u.Username == "" {
		return ErrUsernameRequired
	}

	return nil
}
