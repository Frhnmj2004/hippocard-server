package services

import (
	"context"
	"log"

	"github.com/Frhnmj2004/hippocard-server/pkg/firebase"
	"github.com/gofiber/fiber/v2"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	//"google.golang.org/api/iterator"
)

// AuthService handles authentication operations with Firebase
type AuthService struct {
	AuthClient *firebase.AuthClient
}

// NewAuthService initializes a new AuthService with Firebase Auth client
func NewAuthService(authClient *firebase.AuthClient) *AuthService {
	return &AuthService{
		AuthClient: authClient,
	}
}

// Login authenticates a user and returns a Firebase JWT
func (as *AuthService) Login(c *fiber.Ctx, email, password string) (string, error) {
	if email == "" || password == "" {
		return "", fiber.NewError(fiber.StatusBadRequest, "Email and password are required")
	}

	ctx := context.Background()

	// Attempt to find the user by email
	user, err := as.AuthClient.Client.GetUserByEmail(ctx, email)
	if err != nil {
		if status.Code(err) == codes.NotFound {
			return "", fiber.NewError(fiber.StatusUnauthorized, "Invalid email or password")
		}
		log.Printf("Failed to find user with email %s: %v", email, err)
		return "", fiber.NewError(fiber.StatusInternalServerError, "Login failed")
	}

	// Generate a custom token for the user
	customToken, err := as.AuthClient.Client.CustomToken(ctx, user.UID)
	if err != nil {
		log.Printf("Failed to generate custom token for user %s: %v", user.UID, err)
		return "", fiber.NewError(fiber.StatusInternalServerError, "Login failed")
	}

	// Simulate exchanging the custom token for an ID token (placeholder)
	idToken, err := as.exchangeCustomTokenForIDToken(ctx, customToken)
	if err != nil {
		log.Printf("Failed to exchange custom token for ID token: %v", err)
		return "", fiber.NewError(fiber.StatusInternalServerError, "Login failed")
	}

	// Optionally, set or verify custom claims (e.g., role)
	userRecord, err := as.AuthClient.GetUserByUID(user.UID)
	if err != nil {
		log.Printf("Failed to get user record for %s: %v", user.UID, err)
		return "", fiber.NewError(fiber.StatusInternalServerError, "Login failed")
	}
	role, ok := userRecord.CustomClaims["role"].(string)
	if !ok || role == "" {
		log.Printf("No role found for user %s, setting default to 'patient'", user.UID)
		if err := as.AuthClient.SetCustomClaims(user.UID, "patient"); err != nil {
			log.Printf("Failed to set default role for user %s: %v", user.UID, err)
			return "", fiber.NewError(fiber.StatusInternalServerError, "Login failed")
		}
	}

	return idToken, nil
}

// exchangeCustomTokenForIDToken simulates exchanging a custom token for an ID token
// In production, this would be handled by the frontend Firebase SDK or REST API
func (as *AuthService) exchangeCustomTokenForIDToken(ctx context.Context, customToken string) (string, error) {
	// This is a placeholder—real implementation requires frontend or Firebase REST API
	// For testing, we'll return a mock ID token (replace with actual Firebase API call)
	log.Printf("Exchanging custom token for ID token (placeholder implementation)")
	return "mock-id-token-" + customToken, nil
}

// implemented in middleware
//Firebase Authentication logic
