package model

import (
	"subscriber-service/pkg/id"
)

type Post struct {
	NewsletterID id.Newsletter `json:"newsletter_id" validate:"required"`
	Title        string        `json:"title" validate:"required"`
	Body         string        `json:"body" validate:"required"`
}
