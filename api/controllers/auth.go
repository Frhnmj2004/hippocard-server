package controllers

import (
	"github.com/Frhnmj2004/hippocard-server/internals/services"
	"github.com/Frhnmj2004/hippocard-server/pkg/firebase"

	"github.com/gofiber/fiber/v2"
)

// AuthController handles authentication-related endpoints
type AuthController struct {
	AuthService *services.AuthService
	AuthClient  *firebase.AuthClient // Expects *firebase.AuthClient
}

// NewAuthController initializes a new AuthController with the AuthClient
func NewAuthController(authClient *firebase.AuthClient) *AuthController {
	authService := services.NewAuthService(authClient) // Initialize AuthService with Firebase Auth client
	return &AuthController{
		AuthService: authService,
		AuthClient:  authClient,
	}
}

// LoginHandler handles user login and returns a Firebase JWT token
func (ac *AuthController) LoginHandler(c *fiber.Ctx) error {
	type Request struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	var req Request
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	// Call AuthService to authenticate and get a token
	token, err := ac.AuthService.Login(c, req.Email, req.Password)
	if err != nil {
		// Use fiber.Error.Status to get the status code
		status := fiber.StatusInternalServerError
		if e, ok := err.(*fiber.Error); ok {
			status = e.Code
		}
		return c.Status(status).JSON(fiber.Map{"error": err.Error()})
	}

	// Return the token (mocked for now, update for real Firebase ID token in production)
	return c.JSON(fiber.Map{
		"token":   token,
		"message": "Login successful",
	})
}
