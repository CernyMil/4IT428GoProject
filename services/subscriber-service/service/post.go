package service

import (
	"context"
	"embed"
	"fmt"
	"os"

	token "subscriber-service/pkg/token"
	"subscriber-service/service/mail"
	svcmodel "subscriber-service/service/model"
	"subscriber-service/transport/api/v1/model"
)

//go:embed templates/newsletter_post.html
var templateFS_NewsPost embed.FS

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
			Content:         post.Content,
			UnsubscribeLink: baseUrl + "/subscriber-service/api/v1/subscriptions" + "/unsubscribe?token=" + info.Token,
		})
		if err != nil {
			return fmt.Errorf("failed to prepare HTML for email %s: %w", info.Email, err)
		}

		claims, err := token.ParseJWT(info.Token)
		if err != nil {
		}

		subcscriptionIdStr, ok := claims["subscriptionId"].(string)
		if !ok {
			return fmt.Errorf("invalid subscriptionId %s in token claims for %s", subcscriptionIdStr, info.Email)
		}

		if err != nil {
			return fmt.Errorf("failed to prepare email for %s: %w", info.Email, err)
		}

		if err := mail.SendMail([]string{info.Email}, post.Title, html); err != nil {
			return fmt.Errorf("failed to send email to %s: %w", info.Email, err)
		}
	}

	return nil
}
