package repository

import (
	"cloud.google.com/go/firestore"
)

type Repository struct {
	client *firestore.Client
}

func NewFirestoreRepository(client *firestore.Client) (*Repository, error) {
	return &Repository{
		client: client,
	}, nil
}
