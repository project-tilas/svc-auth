package domain

import "time"

type Token struct {
	ID        string    `json:"id,omitempty" bson:"_id"`
	UserID    string    `json:"userId" bson:"userId"`
	Token     string    `json:"token"`
	ExpiresAt time.Time `json:"expires"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}
