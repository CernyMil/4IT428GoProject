package repository

import (
	"context"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/iterator"

	"subscriber-api/pkg/id"
)

func (r *Repository) DeleteNewsletter(ctx context.Context, newsletter id.Newsletter) error {
	client := r.client
	err := deleteCollection(ctx, client, "subscriptions", 1000)
	if err != nil {
		return err
	}

	_, err = client.Collection("subscription_service_newsletters").Doc(newsletter.String()).Delete(ctx)
	if err != nil {
		return err
	}
	return nil

}

func deleteCollection(ctx context.Context, client *firestore.Client, collectionName string, batchSize int) error {
	// Instantiate a client

	col := client.Collection(collectionName)
	bulkwriter := client.BulkWriter(ctx)

	for {
		// Get a batch of documents
		iter := col.Limit(batchSize).Documents(ctx)
		numDeleted := 0

		// Iterate through the documents, adding
		// a delete operation for each one to the BulkWriter.
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

		// If there are no documents to delete,
		// the process is over.
		if numDeleted == 0 {
			bulkwriter.End()
			break
		}

		bulkwriter.Flush()
	}
	return nil
}
