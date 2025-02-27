// Firebase Authentication SDK setup
package firebase

import (
	"context"
	"log"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/auth"
	"github.com/Frhnmj2004/hippocard-server/configs"
	"google.golang.org/api/option"
)

type AuthClient struct {
	Client *auth.Client
}

func NewAuthClient(config *configs.Config) (*AuthClient, error) {
	opt := option.WithCredentialsFile(config.Firebase.CredentialsPath)
	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		log.Printf("Failed to initialize Firebase app: %v", err)
		return nil, err
	}

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
