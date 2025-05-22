package repository

import (
	"context"
	"subscriber-api/pkg/id"
	"subscriber-api/repository/model"
	svcmodel "subscriber-api/service/model"

	"cloud.google.com/go/firestore"
)

func AddSubscription(ctx context.Context, client *firestore.Client, subscription svcmodel.Subscription) error {
	newsletterID := subscription.NewsletterID.String()
	subscriptionID := subscription.ID.String()
	storeSubscription := model.Subscription{
		CreatedAt: subscription.CreatedAt,
		Email:     subscription.Subscriber.Email,
	}

	_, err := client.Collection("newsletters").Doc(newsletterID).Collection("subscriptions").Doc(subscriptionID).Set(ctx, storeSubscription)

	if err != nil {
		return err
	}
	return nil
}

func RemoveSubscription(ctx context.Context, client *firestore.Client, newsletter id.Newsletter, email string) error {

	_, err := client.Collection("newsletters").Doc(newsletter.String()).Collection("subscriptions").Doc(email).Delete(ctx)
	if err != nil {
		return err
	}
	return nil
}
