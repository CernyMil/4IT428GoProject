package model

import (
	"subscriber-api/pkg/id"
	"time"
)

type Subscription struct {
	ID           id.Subscription `json:"id"`
	SubscriberId id.Subscriber   `json:"newsletter_id"`
	Email        string          `json:"subscriber"`
	CreatedAt    time.Time       `json:"created_at"`
}

type Newsletter struct {
	ID id.Newsletter `json:"id" validate:"required"`
}
