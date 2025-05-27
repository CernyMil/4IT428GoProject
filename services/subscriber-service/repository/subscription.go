package repository

import (
	"context"
	"time"

	"google.golang.org/api/iterator"

	"subscriber-service/pkg/id"
	"subscriber-service/repository/model"
	svcmodel "subscriber-service/service/model"
)

func (r *Repository) AddSubscription(ctx context.Context, newsletterId id.Newsletter, subscriptionId id.Subscription, email string, token string) (*svcmodel.Subscription, error) {
	client := r.client

	storeSubscription := map[string]interface{}{
		"email":      email,
		"token":      token,
		"created_at": time.Now(),
	}

	if _, err := client.Collection("subscription_service_newsletters").Doc(newsletterId.String()).Collection("subscriptions").Doc(subscriptionId.String()).Set(ctx, storeSubscription); err != nil {
		return nil, err
	}

	subscription := svcmodel.Subscription{
		ID:           subscriptionId,
		NewsletterID: newsletterId,
		CreatedAt:    storeSubscription["created_at"].(time.Time),
		Email:        email,
	}

	return &subscription, nil
}

func (r *Repository) DeleteSubscription(ctx context.Context, newsletterId string, subscriptionId string) error {
	client := r.client
	_, err := client.Collection("subscription_service_newsletters").Doc(newsletterId).Collection("subscriptions").Doc(subscriptionId).Delete(ctx)
	if err != nil {
		return err
	}
	return nil
}

func (r *Repository) GetSubscribers(ctx context.Context, newsletterId id.Newsletter) ([]model.SubscriberInfo, error) {
	client := r.client
	iter := client.Collection("subscription_service_newsletters").Doc(newsletterId.String()).Collection("subscriptions").Documents(ctx)
	defer iter.Stop()

	var subscribers []model.SubscriberInfo
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}
		data := doc.Data()
		email, _ := data["email"].(string)
		token, _ := data["token"].(string)
		subscribers = append(subscribers, model.SubscriberInfo{
			Email: email,
			Token: token,
		})
	}
	return subscribers, nil
}
