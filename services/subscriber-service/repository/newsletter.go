package repository

import (
	"context"
	"subscriber-api/pkg/id"

	"cloud.google.com/go/firestore"
)

func DeleteSubscription(ctx context.Context, client *firestore.Client, newsletter id.Newsletter, email string) error {

	_, err := client.Collection("newsletters").Doc(newsletter.String()).Collection("subscriptions").Doc(email).Delete(ctx)
	if err != nil {
		return err
	}
	return nil
}
