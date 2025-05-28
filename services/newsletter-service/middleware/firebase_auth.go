package middleware

import (
	"context"
	"net/http"
	"strings"

	"firebase.google.com/go"
)

// FirebaseAuthMiddleware verifies Firebase ID tokens for protected routes.
func FirebaseAuthMiddleware(app *firebase.App, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		// Verify the token using Firebase
		client, err := app.Auth(context.Background())
		if err != nil {
			http.Error(w, "Failed to initialize Firebase Auth client", http.StatusInternalServerError)
			return
		}

		token, err := client.VerifyIDToken(context.Background(), tokenString)
		if err != nil {
			http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
			return
		}

		// Add editor information to the request context
		ctx := context.WithValue(r.Context(), "editorID", token.UID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
