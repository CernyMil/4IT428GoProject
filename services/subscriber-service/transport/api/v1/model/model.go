package model

import (
	"subscriber-api/pkg/id"
	"time"
)

type Subscriber struct {
	SubcriberID id.Subscriber `json:"subscriber_id" validate:"required"`
	Email       string        `json:"email" validate:"email" required:"true"`
}

type Subscription struct {
	ID           id.Subscription `json:"id" validate:"required"`
	NewsletterID id.Newsletter   `json:"newsletter_id" validate:"required"`
	Subscriber   Subscriber      `json:"subscriber" validate:"required"`
	CreatedAt    time.Time       `json:"created_at"`
}

type SubscribeRequest struct {
	NewsletterID id.Newsletter `json:"newsletter_id" validate:"required,uuid"`
	Email        string        `json:"email" validate:"required,email"`
}
