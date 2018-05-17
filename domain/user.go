package domain

import (
	"time"

	"golang.org/x/crypto/bcrypt"
)

// User represents a single user in the tilas platform
type User struct {
	ID        string    `json:"id" bson:"_id"`
	Username  string    `json:"username"`
	Password  string    `json:"password,omitempty"`
	CreatedAt time.Time `json:"createdAt,omitempty"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func (u *User) ClearPassword() {
	u.Password = ""
}

func (u *User) EncryptPassword() error {
	b, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(b)
	return nil
}

func (u *User) ComparePassword(s string) error {
	return bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(s))
}
