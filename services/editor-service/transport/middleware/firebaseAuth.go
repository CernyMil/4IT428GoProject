package middleware

import (
	"context"
	"fmt"

	fb "firebase.google.com/go"
	"firebase.google.com/go/auth"
	"google.golang.org/api/option"
)

type FirebaseAuth struct {
	Client *auth.Client
}

func NewFirebaseAuth(credPath string) (*FirebaseAuth, error) {
	app, err := fb.NewApp(context.Background(), nil, option.WithCredentialsFile(credPath))
	if err != nil {
		return nil, fmt.Errorf("firebase init error: %w", err)
	}

	client, err := app.Auth(context.Background())
	if err != nil {
		return nil, fmt.Errorf("firebase auth client error: %w", err)
	}

	return &FirebaseAuth{Client: client}, nil
}

func (f *FirebaseAuth) VerifyIDToken(ctx context.Context, idToken string) (*auth.Token, error) {
	return f.Client.VerifyIDToken(ctx, idToken)
}
