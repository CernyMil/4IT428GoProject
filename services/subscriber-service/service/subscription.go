package service

import (
	"context"
	"os"

	"github.com/google/uuid"

	"subscriber-api/pkg/id"
	token "subscriber-api/pkg/token"
	"subscriber-api/service/mail"
	svcmodel "subscriber-api/service/model"
)

func (s Service) SubscribeToNewsletter(ctx context.Context, subReq svcmodel.SubscribeRequest) error {

	newsletterId := subReq.NewsletterID.String()
	token, err := token.GenerateEncryptedToken(subReq.Email)
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

	token, err := token.GenerateEncryptedToken(subscriptionId.String())
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

func (s Service) UnsubscribeFromNewsletter(ctx context.Context, newsletterId id.Newsletter, subscriptionId id.Subscription) error {
	if err := s.repository.DeleteSubscription(ctx, newsletterId, subscriptionId); err != nil {
		return err
	}
	return nil
}

func sendConfirmationRequestMail(email string, newsletterId string, token string) error {
	baseUrl := os.Getenv("BASE_URL")
	confirmLink := baseUrl + "/api/v1/newsletters/" + newsletterId + "/confirm?token=" + token

	html, err := mail.PrepareHTML("templates/confirmation_request.html", struct {
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

	html, err := mail.PrepareHTML("templates/confirmation_request.html", struct {
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
