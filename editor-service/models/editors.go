package models

import "time"

type Editor struct {
	ID             string    `json:"id"` // UUID
	Email          string    `json:"email"`
	FirstName      string    `json:"first_name"`
	LastName       string    `json:"last_name"`
	CreatedAt      time.Time `json:"created_at"`
	HashedPassword string    `db:"hashed_password"`
}
