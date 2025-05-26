package model

import (
	"time"
)

type StoreSubscription struct {
	Email     string    `json:"email"`
	Token     string    `json:"token"`
	CreatedAt time.Time `json:"created_at"`
}

type SubscriberInfo struct {
	Email string `json:"email"`
	Token string `json:"token"`
}
