package middleware

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	fb "firebase.google.com/go"
	"firebase.google.com/go/auth"
	"google.golang.org/api/option"
)

type FirebaseAuth struct {
	Client *auth.Client
	APIKey string
}

func NewFirebaseAuth(credPath, apiKey string) (*FirebaseAuth, error) {
	app, err := fb.NewApp(context.Background(), nil, option.WithCredentialsFile(credPath))
	if err != nil {
		return nil, fmt.Errorf("firebase init error: %w", err)
	}
	client, err := app.Auth(context.Background())
	if err != nil {
		return nil, fmt.Errorf("firebase auth client error: %w", err)
	}
	return &FirebaseAuth{Client: client, APIKey: apiKey}, nil
}

func (f *FirebaseAuth) VerifyIDToken(ctx context.Context, idToken string) (*auth.Token, error) {
	return f.Client.VerifyIDToken(ctx, idToken)
}
func (f *FirebaseAuth) VerifyPasswordWithREST(ctx context.Context, email, password string) (string, error) {
	url := fmt.Sprintf("https://identitytoolkit.googleapis.com/v1/accounts:signInWithPassword?key=%s", f.APIKey)
	payload := map[string]interface{}{
		"email":             email,
		"password":          password,
		"returnSecureToken": true,
	}
	body, _ := json.Marshal(payload)
	req, _ := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	respBody, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != 200 {
		return "", fmt.Errorf("firebase REST error: %s", respBody)
	}
	var result struct {
		IDToken string `json:"idToken"`
	}
	if err := json.Unmarshal(respBody, &result); err != nil {
		return "", err
	}
	return result.IDToken, nil
}
