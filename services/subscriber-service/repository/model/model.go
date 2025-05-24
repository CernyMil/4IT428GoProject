package model

import (
	"subscriber-api/pkg/id"
	"time"
)

type StoreSubscription struct {
	Email     string    `json:"email"`
	Token     string    `json:"token"`
	CreatedAt time.Time `json:"created_at"`
}

type Newsletter struct {
	ID id.Newsletter `json:"id"`
}

type SubscriberInfo struct {
	Email string `json:"email"`
	Token string `json:"token"`
}
