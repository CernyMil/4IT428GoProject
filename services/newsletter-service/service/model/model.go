package model

import (
	id "newsletter-service/pkg/id"
	"time"
)

// CreateNewsletterInput is the input data required to create a newsletter.
type CreateNewsletterInput struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

// UpdateNewsletterInput is the input data required to update a newsletter.
type UpdateNewsletterInput struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

// Newsletter represents a newsletter with an ID, title, and description.
type Newsletter struct {
	ID          id.Newsletter `json:"id"`
	Title       string        `json:"title"`
	Description string        `json:"description"`
	CreatedAt   time.Time     `json:"created_at"` // ISO 8601 format
}

// CreatePostInput is the input data required to create a post.
type CreatePostInput struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

// UpdatePostInput is the input data required to update a post.
type UpdatePostInput struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

type Post struct {
	ID           id.Post       `json:"id"`
	NewsletterID id.Newsletter `json:"newsletter_id"`
	Title        string        `json:"title"`
	Content      string        `json:"content"`
	Published    bool          `json:"published"`
	CreatedAt    time.Time     `json:"created_at"` // ISO 8601 format
}

type PostToPublish struct {
	NewsletterID string `json:"newsletter_id"`
	Title        string `json:"title"`
	Content      string `json:"content"`
}
