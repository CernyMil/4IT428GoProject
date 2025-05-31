package service

import (
	"context"
	"embed"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"

	"subscriber-service/pkg/id"
	token "subscriber-service/pkg/token"
	"subscriber-service/service/mail"
	svcmodel "subscriber-service/service/model"
)

var validate = validator.New()

func (s Service) SubscribeToNewsletter(ctx context.Context, subReq svcmodel.SubscribeRequest) error {
	claims := map[string]interface{}{
		"email":        subReq.Email,
		"newsletterId": subReq.NewsletterID,
	}
	token, err := token.GenerateJWT(claims, 24*time.Hour)
	if err != nil {
		return err
	}

	if err := validate.Var(subReq.Email, "required,email"); err != nil {
		return fmt.Errorf("invalid email: %w", err)
	}

	if err := validate.Var(subReq.NewsletterID, "required"); err != nil {
		return fmt.Errorf("invalid newsletter ID: %w", err)
	}

	if err := sendConfirmationRequestMail(subReq.Email, subReq.NewsletterID.String(), token); err != nil {
		return err
	}
	return nil
}

func (s Service) ConfirmSubscription(ctx context.Context, tokenString string) (svcmodel.Subscription, error) {
	claims, err := token.ParseJWT(tokenString)
	if err != nil {
		return svcmodel.Subscription{}, fmt.Errorf("invalid or expired token: %w", err)
	}

	email, ok := claims["email"].(string)
	if !ok {
		return svcmodel.Subscription{}, errors.New("invalid email in token claims")
	}

	newsletterIdStr, ok := claims["newsletterId"].(string)
	if !ok {
		return svcmodel.Subscription{}, errors.New("invalid newsletterId in token claims")
	}

	var newsletterID id.Newsletter
	if err := newsletterID.FromString(newsletterIdStr); err != nil {
		return svcmodel.Subscription{}, fmt.Errorf("invalid newsletterId: %w", err)
	}

	subscription := svcmodel.Subscription{
		ID:           id.Subscription(uuid.New()),
		NewsletterID: newsletterID,
		Email:        email,
		CreatedAt:    time.Now(),
		Token:        tokenString,
	}

	if err := validate.Struct(subscription); err != nil {
		return svcmodel.Subscription{}, fmt.Errorf("invalid subscription request: %w", err)
	}

	claimsSub := map[string]interface{}{
		"email":          subscription.Email,
		"newsletterId":   subscription.NewsletterID,
		"subscriptionId": subscription.ID,
	}

	token, err := token.GenerateJWT(claimsSub, -1)
	if err != nil {
		return svcmodel.Subscription{}, err
	}

	if err := s.repository.AddSubscription(ctx, subscription); err != nil {
		return svcmodel.Subscription{}, err
	}

	if err := sendConfirmationMail(subscription.Email, subscription.NewsletterID.String(), token); err != nil {
		return svcmodel.Subscription{}, err
	}

	return subscription, nil
}

func (s Service) UnsubscribeFromNewsletter(ctx context.Context, tokenString string) error {
	claims, err := token.ParseJWT(tokenString)
	if err != nil {
	}

	subcscriptionIdStr, ok := claims["subscriptionId"].(string)
	if !ok {
		return errors.New("invalid subscriptionId in token claims")
	}
	newsletterIdStr, ok := claims["newsletterId"].(string)
	if !ok {
		return errors.New("invalid newsletterId in token claims")
	}

	var unsubReq svcmodel.UnsubscribeRequest
	if err := unsubReq.SubscriptionID.FromString(subcscriptionIdStr); err != nil {
		return fmt.Errorf("invalid subscriptionId: %w", err)
	}
	if err := unsubReq.NewsletterID.FromString(newsletterIdStr); err != nil {
		return fmt.Errorf("invalid newsletterId: %w", err)
	}

	if err := s.repository.DeleteSubscription(ctx, unsubReq); err != nil {
		return err
	}
	return nil
}

//go:embed templates/confirmation_request.html
var templateFS_ConfReq embed.FS

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

//go:embed templates/confirmation_request.html
var templateFS_Conf embed.FS

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
