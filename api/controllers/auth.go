// Authentication end points firebase.
package controllers

import (
	//controller called in routes and routes is called in controllers causing cycle. need to fix this.
	"github.com/Frhnmj2004/hippocard-server/api/routes"
	"github.com/gofiber/fiber/v2"
)

// AuthController holds authentication-related handlers
type AuthController struct {
	Repo *routes.Repository
}

// NewAuthController creates a new AuthController
func NewAuthController(repo *routes.Repository) *AuthController {
	return &AuthController{Repo: repo}
}

func (ac *AuthController) LoginHandler(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{"message": "Login not implemented yet"})
}
