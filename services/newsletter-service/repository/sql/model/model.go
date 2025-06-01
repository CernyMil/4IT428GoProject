package model

import (
	//"newsletter-service/repository"
	"newsletter-service/pkg/id"
	"time"
)

// Newsletter represents a newsletter entity in the database.
type Newsletter struct {
	ID        id.Newsletter `json:"id"`
	Subject   string        `json:"subject"`
	Body      string        `json:"body"`
	CreatedAt time.Time     `json:"created_at"`
	EditorID  string        `json:"editor_id"`
}

// UpdateNewsletterInput represents the input for updating a newsletter.
type UpdateNewsletterInput struct {
	Subject string `json:"subject"`
	Body    string `json:"body"`
}

// Post represents a post entity in the database.
type Post struct {
	ID           id.Post       `json:"id"`
	NewsletterID id.Newsletter `json:"newsletter_id"`
	Title        string        `json:"title"`
	Content      string        `json:"content"`
	CreatedAt    time.Time     `json:"created_at"`
	Published    bool          `json:"published"`
}
