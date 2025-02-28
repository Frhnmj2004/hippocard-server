package controllers

import (
	"github.com/Frhnmj2004/hippocard-server/api/routes"

	"github.com/gofiber/fiber/v2"
)

type AuthController struct {
	Repo *routes.Repository
}

func NewAuthController(repo *routes.Repository) *AuthController {
	return &AuthController{Repo: repo}
}

func (ac *AuthController) LoginHandler(c *fiber.Ctx) error {
	type Request struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	var req Request
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}
	// TODO: Implement Firebase Auth login to get JWT token
	// For now, return a placeholder response
	return c.JSON(fiber.Map{
		"message": "Login not implemented yet—use Firebase Auth to get JWT",
		"email":   req.Email,
	})
}
