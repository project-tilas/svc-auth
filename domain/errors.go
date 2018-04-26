package domain

import (
	"errors"
)

// Validation errors
var (
	ErrUsernameRequired  = errors.New("Username is required")
	ErrUserNotFound      = errors.New("User not found")
	ErrUserAlreadyExists = errors.New("User already exists")
	ErrPasswordRequired  = errors.New("Password is required")
)
