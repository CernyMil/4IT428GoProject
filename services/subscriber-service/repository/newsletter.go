package repository

import (
	"context"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/iterator"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"subscriber-service/pkg/id"
)

func (r *Repository) CreateNewsletter(ctx context.Context, newsletterId id.Newsletter) error {
	client := r.client

	newsletter := map[string]interface{}{
		"id": newsletterId.String(),
	}

	_, err := client.Collection("subscription_service_newsletters").Doc(newsletterId.String()).Create(ctx, newsletter)
	if err != nil {
		return err
	}
	return nil
}

func (r *Repository) GetNewsletterById(ctx context.Context, newsletterId id.Newsletter) (id.Newsletter, error) {
	client := r.client
	_, err := client.Collection("subscription_service_newsletters").Doc(newsletterId.String()).Get(ctx)
	if err != nil {
		if status.Code(err) == codes.NotFound {
			return id.Newsletter{}, nil // Newsletter does not exist
		}
		return id.Newsletter{}, err // Other error
	}

	return newsletterId, nil
}

func (r *Repository) DeleteNewsletter(ctx context.Context, newsletterId id.Newsletter) error {
	client := r.client
	err := deleteSubscriptionsForNewsletter(ctx, client, newsletterId.String(), 1000)
	if err != nil {
		return err
	}

	_, err = client.Collection("subscription_service_newsletters").Doc(newsletterId.String()).Delete(ctx)
	if err != nil {
		return err
	}
	return nil

}

func deleteSubscriptionsForNewsletter(ctx context.Context, client *firestore.Client, newsletterId string, batchSize int) error {
	col := client.Collection("subscription_service_newsletters").Doc(newsletterId).Collection("subscriptions")
	bulkwriter := client.BulkWriter(ctx)

	for {
		iter := col.Limit(batchSize).Documents(ctx)
		numDeleted := 0

		// Iterate through the documents, adding a delete operation for each one to the BulkWriter.
		for {
			doc, err := iter.Next()
			if err == iterator.Done {
				break
			}
			if err != nil {
				return err
			}
			bulkwriter.Delete(doc.Ref)
			numDeleted++
		}

		// If there are no documents to delete, the process is over.
		if numDeleted == 0 {
			bulkwriter.End()
			break
		}

		bulkwriter.Flush()
	}
	return nil
}
