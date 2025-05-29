package service

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"
)

// Newsletter represents a newsletter entity.
type Newsletter struct {
	ID        string    `json:"id"`
	Subject   string    `json:"subject"`
	Body      string    `json:"body"`
	CreatedAt time.Time `json:"created_at"`
}

// Post represents a post entity associated with a newsletter.
type Post struct {
	ID           string    `json:"id"`
	NewsletterID string    `json:"newsletter_id"`
	Title        string    `json:"title"`
	Content      string    `json:"content"`
	CreatedAt    time.Time `json:"created_at"`
	Published    bool      `json:"published"`
}

// Input types for creating and updating newsletters and posts.
type CreateNewsletterInput struct {
	Subject string `json:"subject"`
	Body    string `json:"body"`
}

type UpdateNewsletterInput struct {
	Subject string `json:"subject"`
	Body    string `json:"body"`
}

type CreatePostInput struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

type UpdatePostInput struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

// Service defines the business logic for newsletters and posts.
type Service interface {
	CreateNewsletter(ctx context.Context, input CreateNewsletterInput) (*Newsletter, error)
	ListNewsletters(ctx context.Context) ([]Newsletter, error)
	UpdateNewsletter(ctx context.Context, id string, input UpdateNewsletterInput) (*Newsletter, error)

	// Post-related methods
	CreatePost(ctx context.Context, newsletterID string, input CreatePostInput) (*Post, error)
	ListPosts(ctx context.Context, newsletterID string) ([]Post, error)
	UpdatePost(ctx context.Context, newsletterID string, postID string, input UpdatePostInput) (*Post, error)
	DeletePost(ctx context.Context, newsletterID string, postID string) error
}

// Repository defines the data access layer for newsletters and posts.
type Repository interface {
	Save(ctx context.Context, newsletter *Newsletter) error
	FindAll(ctx context.Context) ([]Newsletter, error)
	Update(ctx context.Context, id string, input UpdateNewsletterInput) (*Newsletter, error)

	// Post-related methods
	CreatePost(ctx context.Context, post *Post) error
	FindPostsByNewsletterID(ctx context.Context, newsletterID string) ([]Post, error)
	UpdatePost(ctx context.Context, postID string, post *Post) error
	DeletePost(ctx context.Context, postID string) error
}

// newsletterService implements the Service interface.
type newsletterService struct {
	repo  Repository
	idGen func() string
}

// NewService creates a new instance of newsletterService.
func NewService(repo Repository, idGen func() string) Service {
	return &newsletterService{repo: repo, idGen: idGen}
}

// Newsletter methods
func (s *newsletterService) CreateNewsletter(ctx context.Context, input CreateNewsletterInput) (*Newsletter, error) {
	if input.Subject == "" || input.Body == "" {
		return nil, errors.New("subject and body are required")
	}

	n := &Newsletter{
		ID:        s.idGen(),
		Subject:   input.Subject,
		Body:      input.Body,
		CreatedAt: time.Now(),
	}
	if err := s.repo.Save(ctx, n); err != nil {
		return nil, err
	}

	url := "http://subscriber-service:8083/api/v1/nginx/newsletters"
	payload, err := json.Marshal(n.ID)
	if err != nil {
		return nil, err
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(payload))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("subscriber-service returned status %d", resp.StatusCode)
	}

	return n, nil
}

func (s *newsletterService) ListNewsletters(ctx context.Context) ([]Newsletter, error) {
	return s.repo.FindAll(ctx)
}

func (s *newsletterService) UpdateNewsletter(ctx context.Context, id string, input UpdateNewsletterInput) (*Newsletter, error) {
	if input.Subject == "" || input.Body == "" {
		return nil, errors.New("subject and body are required")
	}
	return s.repo.Update(ctx, id, input)
}

// Post methods
func (s *newsletterService) CreatePost(ctx context.Context, newsletterID string, input CreatePostInput) (*Post, error) {
	if input.Title == "" || input.Content == "" {
		return nil, errors.New("title and content are required")
	}

	p := &Post{
		ID:           s.idGen(),
		NewsletterID: newsletterID,
		Title:        input.Title,
		Content:      input.Content,
		CreatedAt:    time.Now(),
		Published:    false,
	}
	if err := s.repo.CreatePost(ctx, p); err != nil {
		return nil, err
	}
	return p, nil
}

func (s *newsletterService) ListPosts(ctx context.Context, newsletterID string) ([]Post, error) {
	return s.repo.FindPostsByNewsletterID(ctx, newsletterID)
}

func (s *newsletterService) UpdatePost(ctx context.Context, newsletterID string, postID string, input UpdatePostInput) (*Post, error) {
	if input.Title == "" || input.Content == "" {
		return nil, errors.New("title and content are required")
	}

	p := &Post{
		ID:           postID,
		NewsletterID: newsletterID,
		Title:        input.Title,
		Content:      input.Content,
	}
	if err := s.repo.UpdatePost(ctx, postID, p); err != nil {
		return nil, err
	}
	return p, nil
}

func (s *newsletterService) DeletePost(ctx context.Context, newsletterID string, postID string) error {
	return s.repo.DeletePost(ctx, postID)
}
