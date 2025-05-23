package model

import (
	"subscriber-api/pkg/id"
	"time"
)

type StoreSubscription struct {
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}

type Newsletter struct {
	ID id.Newsletter `json:"id" validate:"required"`
}
