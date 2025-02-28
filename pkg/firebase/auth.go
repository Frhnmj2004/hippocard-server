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

// SetCustomClaims sets role-based claims for a user (Admin SDK)
func (ac *AuthClient) SetCustomClaims(uid, role string) error {
	ctx := context.Background()
	params := (&auth.UserToUpdate{}).CustomClaims(map[string]interface{}{
		"role": role,
	})
	_, err := ac.Client.UpdateUser(ctx, uid, params)
	if err != nil {
		log.Printf("Failed to set custom claims for user %s: %v", uid, err)
		return err
	}
	return nil
}

// GetUserByUID retrieves user info to verify or manage roles
func (ac *AuthClient) GetUserByUID(uid string) (*auth.UserRecord, error) {
	ctx := context.Background()
	user, err := ac.Client.GetUser(ctx, uid)
	if err != nil {
		log.Printf("Failed to get user %s: %v", uid, err)
		return nil, err
	}
	return user, nil
}
