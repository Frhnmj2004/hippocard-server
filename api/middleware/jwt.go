package middleware

import (
	"log"
	"strings"

	"github.com/Frhnmj2004/hippocard-server/pkg/firebase"

	"github.com/gofiber/fiber/v2"
)

// AuthMiddleware verifies Firebase JWT tokens and enforces role-based access
func AuthMiddleware(authClient *firebase.AuthClient, requiredRole string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Get the Authorization header (e.g., "Bearer <token>")
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "No Authorization header provided",
			})
		}

		// Split into "Bearer" and token
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid Authorization header format (use 'Bearer <token>')",
			})
		}
		token := parts[1]

		// Verify token with Firebase Auth
		verifiedToken, err := authClient.VerifyIDToken(token)
		if err != nil {
			log.Printf("Token verification failed: %v", err)
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid or expired token",
			})
		}

		// Check role from custom claims (assumes "role" is set in Firebase token)
		role, ok := verifiedToken.Claims["role"].(string)
		if !ok || role != requiredRole {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "Insufficient permissions for this role",
			})
		}

		// Store user ID for handlers
		c.Locals("userID", verifiedToken.UID)

		// Continue to the handler
		return c.Next()
	}
}
