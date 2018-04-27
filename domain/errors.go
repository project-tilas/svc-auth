package domain

import (
	"errors"
)

// Validation errors
var (
	ErrUsernameRequired = errors.New("Username is required")
	ErrPasswordRequired = errors.New("Password is required")
	ErrInvalidPassword  = errors.New("Incorrect Password")
)

// User errors
var (
	ErrUserNotFound      = errors.New("User not found")
	ErrUserAlreadyExists = errors.New("User already exists")
)

// Token errors
var (
	ErrTokenNotFound = errors.New("Token not found")
)
