package middleware

import (
	"context"
	"net/http"
	"strings"

	firebase "firebase.google.com/go"
)

// FirebaseAuthenticator implements the Authenticator interface using Firebase.
type FirebaseAuthenticator struct {
	app *firebase.App
}

// Authenticator defines an interface for authentication middleware.
type Authenticator interface {
	Authenticate(next http.Handler) http.Handler
}

// NewFirebaseAuthenticator creates a new FirebaseAuthenticator.
func NewFirebaseAuthenticator(app *firebase.App) *FirebaseAuthenticator {
	return &FirebaseAuthenticator{app: app}
}

// Authenticate verifies Firebase ID tokens for protected routes.
func (fa *FirebaseAuthenticator) Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		// Verify the token using Firebase
		client, err := fa.app.Auth(context.Background())
		if err != nil {
			http.Error(w, "Failed to initialize Firebase Auth client", http.StatusInternalServerError)
			return
		}

		token, err := client.VerifyIDToken(context.Background(), tokenString)
		if err != nil {
			http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
			return
		}

		// Add user information to the request context
		ctx := context.WithValue(r.Context(), "userID", token.UID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// FirebaseAuthMiddleware is a standalone middleware function for verifying Firebase ID tokens.
func FirebaseAuthMiddleware(firebaseApp *firebase.App) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			tokenString := strings.TrimPrefix(authHeader, "Bearer ")

			// Verify the token using Firebase
			client, err := firebaseApp.Auth(context.Background())
			if err != nil {
				http.Error(w, "Failed to initialize Firebase Auth client", http.StatusInternalServerError)
				return
			}

			token, err := client.VerifyIDToken(context.Background(), tokenString)
			if err != nil {
				http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
				return
			}

			// Add user information to the request context
			ctx := context.WithValue(r.Context(), "userID", token.UID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
