package middleware

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	firebase "firebase.google.com/go"
	"google.golang.org/api/option"
)

func FirebaseAuthMiddleware(credentialsPath string) func(http.Handler) http.Handler {
	opt := option.WithCredentialsFile(credentialsPath)
	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		panic(fmt.Errorf("firebase init error: %w", err))
	}

	authClient, err := app.Auth(context.Background())
	if err != nil {
		panic(fmt.Errorf("firebase auth client error: %w", err))
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			header := r.Header.Get("Authorization")
			if header == "" || !strings.HasPrefix(header, "Bearer ") {
				http.Error(w, "Missing or invalid Authorization header", http.StatusUnauthorized)
				return
			}

			tokenStr := strings.TrimPrefix(header, "Bearer ")
			_, err := authClient.VerifyIDToken(context.Background(), tokenStr)
			if err != nil {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
