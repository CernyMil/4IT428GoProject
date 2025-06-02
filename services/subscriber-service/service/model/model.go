package model

import (
	"subscriber-service/pkg/id"
	"time"
)

type Subscription struct {
	ID           id.Subscription `json:"id" validate:"required"`
	NewsletterID id.Newsletter   `json:"newsletter_id" validate:"required"`
	Email        string          `json:"email" validate:"email" required:"true"`
	CreatedAt    time.Time       `json:"created_at"`
	Token        string          `json:"token" validate:"required"`
}

type PostHTML struct {
	Email           string `json:"Email" validate:"required"`
	Title           string `json:"title" validate:"required"`
	Content         string `json:"content" validate:"required"`
	UnsubscribeLink string `json:"unsubscribeLink" validate:"required"`
}

type SubscribeRequest struct {
	NewsletterID id.Newsletter `json:"newsletter_id" validate:"required"`
	Email        string        `json:"email" validate:"required,email"`
}

type UnsubscribeRequest struct {
	NewsletterID   id.Newsletter   `json:"newsletter_id" validate:"required"`
	SubscriptionID id.Subscription `json:"subscription_id" validate:"required"`
}

type Newsletter struct {
	NewsletterID string `json:"newsletter_id"`
}
