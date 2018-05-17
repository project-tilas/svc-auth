package domain

import (
	"fmt"
)

type NotFoundError struct {
	ID       string
	Resource string
}

func (e *NotFoundError) Error() string {
	if e.ID != "" {
		return fmt.Sprintf("%s %s not found", e.Resource, e.ID)
	}
	return fmt.Sprintf("%s not found", e.ID)
}

type AlreadyExistsError struct {
	Resource string
}

func (e *AlreadyExistsError) Error() string {
	return fmt.Sprintf("%s already exists", e.Resource)
}

type ValidationError struct {
	Reason string
}

func (e *ValidationError) Error() string {
	return e.Reason
}
