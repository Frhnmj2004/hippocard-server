package controllers

import (
	"github.com/Frhnmj2004/hippocard-server/api/routes"
	"github.com/gofiber/fiber/v2"
)

// PharmacistController holds pharmacist-specific handlers
type PharmacistController struct {
	Repo *routes.Repository
}

// NewPharmacistController creates a new PharmacistController
func NewPharmacistController(repo *routes.Repository) *PharmacistController {
	return &PharmacistController{Repo: repo}
}

func (pc *PharmacistController) ActivePrescriptionsHandler(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{"message": "Pharmacist active prescriptions not implemented"})
}

func (pc *PharmacistController) DispensePrescriptionHandler(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{"message": "Pharmacist dispense prescription not implemented"})
}
