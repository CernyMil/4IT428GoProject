package middleware

import (
	"net/http"
)

const (
	authHeader = "Authorization"
	authSchema = "Bearer "
)

func InternalOnlyMiddleware(expectedToken string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get(authHeader)
			const prefix = authSchema
			if len(authHeader) < len(prefix) || authHeader[:len(prefix)] != prefix {
				http.Error(w, "Forbidden", http.StatusForbidden)
				return
			}
			token := authHeader[len(prefix):]
			if token != expectedToken {
				http.Error(w, "Forbidden", http.StatusForbidden)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
