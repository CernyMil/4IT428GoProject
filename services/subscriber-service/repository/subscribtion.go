package repository

import (
	"context"
	"time"

	"subscriber-api/pkg/id"
	"subscriber-api/repository/model"
	svcmodel "subscriber-api/service/model"
)

func (r *Repository) AddSubscription(ctx context.Context, newsletterId id.Newsletter, subscribtionId id.Subscription, email string) (*svcmodel.Subscription, error) {
	client := r.client

	storeSubscription := model.StoreSubscription{
		Email:     email,
		CreatedAt: time.Now(),
	}

	if _, err := client.Collection("subscription_service_newsletters").Doc(newsletterId.String()).Collection("subscriptions").Doc(subscribtionId.String()).Set(ctx, storeSubscription); err != nil {
		return nil, err
	}

	subscription := svcmodel.Subscription{
		ID:           subscribtionId,
		NewsletterID: newsletterId,
		CreatedAt:    storeSubscription.CreatedAt,
		Email:        email,
	}

	return &subscription, nil
}

func (r *Repository) DeleteSubscription(ctx context.Context, newsletterId id.Newsletter, subscribtionId id.Subscription) error {
	client := r.client
	_, err := client.Collection("subscription_service_newsletters").Doc(newsletterId.String()).Collection("subscriptions").Doc(subscribtionId.String()).Delete(ctx)
	if err != nil {
		return err
	}
	return nil
}
