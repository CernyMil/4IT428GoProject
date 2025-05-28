package model

import (
	"subscriber-service/pkg/id"
)

type Post struct {
	ID           id.Post       `json:"id" validate:"required"`
	NewsletterID id.Newsletter `json:"newsletter_id" validate:"required"`
	Title        string        `json:"title" validate:"required"`
	Body         string        `json:"body" validate:"required"`
}
