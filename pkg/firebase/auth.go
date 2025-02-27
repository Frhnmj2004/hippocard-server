// Firebase Authentication SDK setup
package firebase

import (
	"context"
	"log"

	firebase "firebase.google.com/go" // Updated path
	"firebase.google.com/go/auth"
)

type AuthClient struct {
	Client *auth.Client
}

func NewAuthClient(app *firebase.App) (*AuthClient, error) {
	// Initialize the Auth client
	client, err := app.Auth(context.Background())
	if err != nil {
		log.Printf("Failed to initialize Firebase Auth client: %v", err)
		return nil, err
	}

	return &AuthClient{Client: client}, nil
}

// VerifyIDToken verifies a Firebase JWT token and returns the decoded token
func (ac *AuthClient) VerifyIDToken(token string) (*auth.Token, error) {
	verifiedToken, err := ac.Client.VerifyIDToken(context.Background(), token)
	if err != nil {
		log.Printf("Failed to verify ID token: %v", err)
		return nil, err
	}
	return verifiedToken, nil
}
