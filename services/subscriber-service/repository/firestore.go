package repository

import (
	"context"
	"log"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
	"google.golang.org/api/option"
)

func InitializeFirebase() (*firestore.Client, error) {
	// Use a service account
	ctx := context.Background()
	sa := option.WithCredentialsFile("services/subscriber-service/goproject-4c949-firebase-adminsdk-fbsvc-564d6d7e63.json")
	app, err := firebase.NewApp(ctx, nil, sa)
	if err != nil {
		log.Fatalln(err)
	}

	client, err := app.Firestore(ctx)
	if err != nil {
		log.Fatalln(err)
	}
	defer client.Close()

	return client, nil
	/*

		ctx := context.Background()
		ProjectID := os.Getenv("FIREBASE_PROJECT_ID")
		conf := &firebase.Config{ProjectID: ProjectID}
		app, err := firebase.NewApp(ctx, conf)
		if err != nil {
			log.Fatalln(err)
		}

		client, err := app.Firestore(ctx)
		if err != nil {
			log.Fatalln(err)
		}
		return client, nil
	*/
}
