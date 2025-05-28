package model

import "context"

// Service defines the methods required for newsletter and post operations.
type Service interface {
	CreateNewsletter(ctx context.Context, input CreateNewsletterInput) (Newsletter, error)
	ListNewsletters(ctx context.Context) ([]Newsletter, error)
	UpdateNewsletter(ctx context.Context, id string, input UpdateNewsletterInput) (Newsletter, error)
	DeleteNewsletter(ctx context.Context, id string) error

	CreatePost(ctx context.Context, newsletterID string, input CreatePostInput) (Post, error)
	ListPosts(ctx context.Context, newsletterID string) ([]Post, error)
	UpdatePost(ctx context.Context, newsletterID, postID string, input UpdatePostInput) (Post, error)
	DeletePost(ctx context.Context, newsletterID, postID string) error
}

// Define the necessary structs for input and output.
type CreateNewsletterInput struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

type UpdateNewsletterInput struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

type Newsletter struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

type CreatePostInput struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

type UpdatePostInput struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

type Post struct {
	ID      string `json:"id"`
	Title   string `json:"title"`
	Content string `json:"content"`
}
