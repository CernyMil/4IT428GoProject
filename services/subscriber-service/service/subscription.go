package service

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/google/uuid"

	"subscriber-service/pkg/id"
	token "subscriber-service/pkg/token"
	"subscriber-service/service/mail"
	svcmodel "subscriber-service/service/model"
)

func (s Service) SubscribeToNewsletter(ctx context.Context, subReq svcmodel.SubscribeRequest) error {
	newsletterId := subReq.NewsletterID.String()

	claims := map[string]interface{}{
		"email":        subReq.Email,
		"newsletterId": subReq.NewsletterID,
	}

	token, err := token.GenerateJWT(claims, 24*time.Hour)
	if err != nil {
		return err
	}

	if err := sendConfirmationRequestMail(subReq.Email, newsletterId, token); err != nil {
		return err
	}

	return nil
}

func (s Service) ConfirmSubscription(ctx context.Context, subReq svcmodel.SubscribeRequest) (svcmodel.Subscription, error) {
	subscriptionId := id.Subscription(uuid.New())

	claims := map[string]interface{}{
		"email":          subReq.Email,
		"newsletterId":   subReq.NewsletterID,
		"subscriptionId": subscriptionId,
	}

	token, err := token.GenerateJWT(claims, 0)
	if err != nil {
		return svcmodel.Subscription{}, err
	}

	subscription, err := s.repository.AddSubscription(ctx, subReq.NewsletterID, subscriptionId, subReq.Email, token)
	if err != nil {
		return svcmodel.Subscription{}, err
	}

	if err := sendConfirmationMail(subReq.Email, subReq.NewsletterID.String(), token); err != nil {
		return svcmodel.Subscription{}, err
	}

	return *subscription, nil
}

func (s Service) UnsubscribeFromNewsletter(ctx context.Context, unsubReq svcmodel.UnsubscribeRequest) error {
	newsletterId := unsubReq.NewsletterID.String()
	subscriptionId := unsubReq.SubscriptionID.String()

	if err := s.repository.DeleteSubscription(ctx, newsletterId, subscriptionId); err != nil {
		return err
	}
	return nil
}

func sendConfirmationRequestMail(email string, newsletterId string, token string) error {
	baseUrl := os.Getenv("BASE_URL")
	confirmLink := baseUrl + "/api/v1/newsletters/" + newsletterId + "/confirm?token=" + token

	templateContent, err := templateFS_ConfReq.ReadFile("templates/confirmation_request.html")
	if err != nil {
		return fmt.Errorf("failed to read embedded template: %w", err)
	}

	html, err := mail.PrepareHTMLFromBytes(templateContent, struct {
		Email       string
		ConfirmLink string
	}{
		Email:       email,
		ConfirmLink: confirmLink,
	})

	if err != nil {
		return err
	}

	if err := mail.SendMail([]string{email}, "Subscription confirmation request", html); err != nil {
		return err
	}

	return nil
}

func sendConfirmationMail(email string, newsletterId string, token string) error {

	baseUrl := os.Getenv("BASE_URL")
	unsubscribeLink := baseUrl + "/api/v1/newsletters/" + newsletterId + "/unsubscribe?token=" + token

	templateContent, err := templateFS_Conf.ReadFile("templates/confirmation_request.html")
	if err != nil {
		return fmt.Errorf("failed to read embedded template: %w", err)
	}

	html, err := mail.PrepareHTMLFromBytes(templateContent, struct {
		Email           string
		UnsubscribeLink string
	}{
		Email:           email,
		UnsubscribeLink: unsubscribeLink,
	})

	if err != nil {
		return err
	}

	if err := mail.SendMail([]string{email}, "Subscription confirmation", html); err != nil {
		return err
	}

	return nil
}
