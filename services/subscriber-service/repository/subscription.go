package repository

import (
	"context"

	"google.golang.org/api/iterator"

	"subscriber-service/pkg/id"
	"subscriber-service/repository/model"
	svcmodel "subscriber-service/service/model"
)

func (r *Repository) AddSubscription(ctx context.Context, subscription svcmodel.Subscription) error {
	client := r.client

	newsletterIdStr := subscription.NewsletterID.String()
	subscriptionIdStr := subscription.ID.String()

	storeSubscription := map[string]interface{}{
		"email":      subscription.Email,
		"token":      subscription.Token,
		"created_at": subscription.CreatedAt,
	}

	if _, err := client.Collection("subscription_service_newsletters").Doc(newsletterIdStr).Collection("subscriptions").Doc(subscriptionIdStr).Set(ctx, storeSubscription); err != nil {
		return err
	}

	return nil
}

func (r *Repository) DeleteSubscription(ctx context.Context, unsubReq svcmodel.UnsubscribeRequest) error {
	newsletterIdStr := unsubReq.NewsletterID.String()
	subscriptionIdStr := unsubReq.SubscriptionID.String()

	client := r.client
	_, err := client.Collection("subscription_service_newsletters").Doc(newsletterIdStr).Collection("subscriptions").Doc(subscriptionIdStr).Delete(ctx)
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
