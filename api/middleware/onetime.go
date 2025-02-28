// One-time access middleware for hospital
package middleware

import (
	"context"
	"log"
	"strings"
	"time"

	"github.com/Frhnmj2004/hippocard-server/internals/models"
	"github.com/Frhnmj2004/hippocard-server/pkg/firebase"

	"cloud.google.com/go/firestore"
	"github.com/gofiber/fiber/v2"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const accessTTL = 5 * time.Minute // Time-to-live for one-time access (5 minutes)

// OneTimeAccess ensures hospital access is limited to one-time use per token
func OneTimeAccess(authClient *firebase.AuthClient, firestore *firebase.FirestoreClient) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Step 1: Get the Authorization header (e.g., "Bearer <token>")
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "No Authorization header provided",
			})
		}

		// Step 2: Split into "Bearer" and token
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid Authorization header format (use 'Bearer <token>')",
			})
		}
		token := parts[1]

		// Step 3: Verify the token with Firebase Auth
		verifiedToken, err := authClient.VerifyIDToken(token)
		if err != nil {
			log.Printf("Token verification failed: %v", err)
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid or expired token",
			})
		}

		// Step 4: Check if the user has the hospital role
		role, ok := verifiedToken.Claims["role"].(string)
		if !ok || role != "hospital" {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "Insufficient permissions for hospital role",
			})
		}

		// Step 5: Generate a unique access key (endpoint + user ID + timestamp)
		userID := verifiedToken.UID
		nfcID := c.Params("nfc_id") // Assuming this is in the URL path
		accessKey := generateAccessKey(c.Path(), userID, nfcID)

		// Step 6: Check if this access has already been used
		used, err := checkAccess(firestore.Client, accessKey)
		if err != nil {
			log.Printf("Failed to check access: %v", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Internal server error",
			})
		}
		if used {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "This access has already been used",
			})
		}

		// Step 7: Mark this access as used and log it
		err = markAccess(firestore.Client, accessKey, userID, nfcID)
		if err != nil {
			log.Printf("Failed to mark access: %v", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Internal server error",
			})
		}

		// Step 8: Store access metadata in context for logging or cleanup
		c.Locals("accessKey", accessKey)

		// Step 9: Proceed to handler, but this token is now invalid for reuse
		return c.Next()
	}
}

// generateAccessKey creates a unique identifier for one-time access
func generateAccessKey(path, userID, nfcID string) string {
	return path + "_" + userID + "_" + nfcID + "_" + time.Now().UTC().Format(time.RFC3339Nano)
}

// checkAccess queries Firestore to see if access has been used
func checkAccess(client *firestore.Client, accessKey string) (bool, error) {
	ctx := context.Background()
	_, err := client.Collection("transactions").Doc(accessKey).Get(ctx)
	if err != nil {
		if status.Code(err) == codes.NotFound {
			return false, nil // Not found means not used
		}
		return false, err
	}
	return true, nil // Document exists, access used
}

// markAccess logs the access in Firestore and sets an expiration
func markAccess(client *firestore.Client, accessKey, userID, nfcID string) error {
	ctx := context.Background()
	transaction := &models.Transaction{
		ID:         accessKey,
		UserID:     userID,
		PatientID:  "", // Set in controller/service if needed
		NFCID:      nfcID,
		AccessTime: time.Now().UTC(),
	}

	_, err := client.Collection("transactions").Doc(accessKey).Set(ctx, transaction)
	if err != nil {
		return err
	}

	// Optionally, schedule cleanup (e.g., delete after accessTTL)
	go func() {
		time.Sleep(accessTTL)
		_, err := client.Collection("transactions").Doc(accessKey).Delete(context.Background())
		if err != nil {
			log.Printf("Failed to clean up access record: %v", err)
		}
	}()

	return nil
}
