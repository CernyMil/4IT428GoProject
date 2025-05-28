package newsletter

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"
)

type Newsletter struct {
	ID        string    `json:"id"`
	Subject   string    `json:"subject"`
	Body      string    `json:"body"`
	CreatedAt time.Time `json:"created_at"`
}

type CreateNewsletterInput struct {
	Subject string `json:"subject"`
	Body    string `json:"body"`
}

type UpdateNewsletterInput struct {
	Subject string `json:"subject"`
	Body    string `json:"body"`
}

type Service interface {
	CreateNewsletter(ctx context.Context, input CreateNewsletterInput) (*Newsletter, error)
	ListNewsletters(ctx context.Context) ([]Newsletter, error)
	UpdateNewsletter(ctx context.Context, id string, input UpdateNewsletterInput) (*Newsletter, error)
}

type Repository interface {
	Save(ctx context.Context, newsletter *Newsletter) error
	FindAll(ctx context.Context) ([]Newsletter, error)
	Update(ctx context.Context, id string, input UpdateNewsletterInput) (*Newsletter, error)
}

type newsletterService struct {
	repo  Repository
	idGen func() string
}

func NewService(repo Repository, idGen func() string) Service {
	return &newsletterService{repo: repo, idGen: idGen}
}

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
