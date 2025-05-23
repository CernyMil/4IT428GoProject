package model

import (
	"subscriber-api/pkg/id"
	"time"
)

type Subscription struct {
	ID           id.Subscription `json:"id" validate:"required"`
	NewsletterID id.Newsletter   `json:"newsletter_id" validate:"required"`
	Email        string          `json:"email" validate:"email" required:"true"`
	CreatedAt    time.Time       `json:"created_at"`
}

type PostHTML struct {
	Email           string `json:"Email" validate:"required"`
	Title           string `json:"title" validate:"required"`
	Content         string `json:"content" validate:"required"`
	UnsubscribeLink string `json:"unsubscribeLink" validate:"required"`
}

type SubscribeRequest struct {
	NewsletterID id.Newsletter `json:"newsletter_id" validate:"required,uuid"`
	Email        string        `json:"email" validate:"required,email"`
}
