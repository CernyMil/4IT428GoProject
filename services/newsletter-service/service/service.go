package service

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	//"time"
	"newsletter-service/pkg/id"
	svcmodel "newsletter-service/service/model"
)

// Repository defines the data access layer for newsletters and posts.
type Repository interface {
	Save(ctx context.Context, newsletter *svcmodel.Newsletter) error
	FindAll(ctx context.Context) ([]svcmodel.Newsletter, error)
	Update(ctx context.Context, id id.Newsletter, input svcmodel.UpdateNewsletterInput) (*svcmodel.Newsletter, error)
	Delete(ctx context.Context, id id.Newsletter) error

	CreatePost(ctx context.Context, post *svcmodel.Post) error
	FindPostsByNewsletterID(ctx context.Context, newsletterID id.Newsletter) ([]svcmodel.Post, error)
	UpdatePost(ctx context.Context, postID id.Post, post *svcmodel.Post) error
	DeletePost(ctx context.Context, postID id.Post) error
}

// newsletterService implements the Service interface.
type NewsletterService struct {
	repo Repository
}

// NewService creates a new instance of newsletterService.
func NewService(repo Repository) (NewsletterService, error) {
	return NewsletterService{
		repo: repo,
	}, nil
}

// CreateNewsletter creates a new newsletter and notifies subscriber service.
func (s *NewsletterService) CreateNewsletter(ctx context.Context, input svcmodel.CreateNewsletterInput) (*svcmodel.Newsletter, error) {
	if input.Title == "" || input.Description == "" {
		return nil, errors.New("subject and body are required")
	}

	n := &svcmodel.Newsletter{
		ID:          id.NewNewsletter(),
		Title:       input.Title,
		Description: input.Description,
		//CreatedAt: time.Now(),
	}
	if err := s.repo.Save(ctx, n); err != nil {
		return nil, err
	}

	// Notify subscriber-service about new newsletter (optional integration)
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

// ListNewsletters lists all newsletters.
func (s *NewsletterService) ListNewsletters(ctx context.Context) ([]svcmodel.Newsletter, error) {
	return s.repo.FindAll(ctx)
}

// UpdateNewsletter updates an existing newsletter.
func (s *NewsletterService) UpdateNewsletter(ctx context.Context, newsletterID id.Newsletter, input svcmodel.UpdateNewsletterInput) (*svcmodel.Newsletter, error) {
	if input.Title == "" || input.Description == "" {
		return nil, errors.New("subject and body are required")
	}
	return s.repo.Update(ctx, newsletterID, input)
}

// DeleteNewsletter deletes a newsletter by ID.
func (s *NewsletterService) DeleteNewsletter(ctx context.Context, newsletterID id.Newsletter) error {
	return s.repo.Delete(ctx, newsletterID)
}

// CreatePost creates a new post for a newsletter.
func (s *NewsletterService) CreatePost(ctx context.Context, newsletterID id.Newsletter, input svcmodel.CreatePostInput) (*svcmodel.Post, error) {
	if input.Title == "" || input.Content == "" {
		return nil, errors.New("title and content are required")
	}

	p := &svcmodel.Post{
		ID:           id.NewPost(),
		NewsletterID: newsletterID,
		Title:        input.Title,
		Content:      input.Content,
		//CreatedAt:    time.Now(),
		Published: false,
	}

	if err := s.repo.CreatePost(ctx, p); err != nil {
		return nil, err
	}
	return p, nil
}

// ListPosts lists all posts under a newsletter.
func (s *NewsletterService) ListPosts(ctx context.Context, newsletterID id.Newsletter) ([]svcmodel.Post, error) {
	return s.repo.FindPostsByNewsletterID(ctx, newsletterID)
}

// UpdatePost updates a post by ID.
func (s *NewsletterService) UpdatePost(ctx context.Context, newsletterID id.Newsletter, postID id.Post, input svcmodel.UpdatePostInput) (*svcmodel.Post, error) {
	if input.Title == "" || input.Content == "" {
		return nil, errors.New("title and content are required")
	}

	p := &svcmodel.Post{
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

// DeletePost deletes a post by ID.
func (s *NewsletterService) DeletePost(ctx context.Context, newsletterID id.Newsletter, postID id.Post) error {
	return s.repo.DeletePost(ctx, postID)
}
