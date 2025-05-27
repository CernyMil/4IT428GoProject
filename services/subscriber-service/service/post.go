package service

import (
	"context"
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

	for _, info := range subscriberInfo {

		html, err := mail.PrepareHTML("templates/newsletter_post.html", svcmodel.PostHTML{
			Email:           info.Email,
			Title:           post.Title,
			Content:         post.Body,
			UnsubscribeLink: baseUrl + "/api/v1/newsletters/" + post.NewsletterID.String() + "/unsubscribe?token=" + info.Token,
		})

		if err != nil {
			return err
		}

		if err := mail.SendMail([]string{info.Email}, post.Title, html); err != nil {
			return err
		}
	}

	return nil
}
