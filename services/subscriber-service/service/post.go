package service

import (
	"context"
	"fmt"
	"os"

	"subscriber-service/service/mail"
	svcmodel "subscriber-service/service/model"
	"subscriber-service/transport/api/v1/model"
)

func (s Service) SendPublishedPost(ctx context.Context, post model.Post) error {
	subscriberInfo, err := s.repository.GetSubscribers(ctx, post.NewsletterID)
	if err != nil {
		return err
	}
	baseUrl := os.Getenv("BASE_URL")

	// Read the embedded template once (more efficient)
	templateContent, err := templateFS_NewsPost.ReadFile("templates/newsletter_post.html")
	if err != nil {
		return fmt.Errorf("failed to read newsletter post template: %w", err)
	}

	for _, info := range subscriberInfo {
		html, err := mail.PrepareHTMLFromBytes(templateContent, svcmodel.PostHTML{
			Email:           info.Email,
			Title:           post.Title,
			Content:         post.Body,
			UnsubscribeLink: baseUrl + "/api/v1/newsletters/" + post.NewsletterID.String() + "/unsubscribe?token=" + info.Token,
		})

		if err != nil {
			return fmt.Errorf("failed to prepare email for %s: %w", info.Email, err)
		}

		if err := mail.SendMail([]string{info.Email}, post.Title, html); err != nil {
			return fmt.Errorf("failed to send email to %s: %w", info.Email, err)
		}
	}

	return nil
}
