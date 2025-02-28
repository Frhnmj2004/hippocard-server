package controllers

import (
	"github.com/Frhnmj2004/hippocard-server/api/routes"
	"github.com/Frhnmj2004/hippocard-server/internals/services"

	"github.com/gofiber/fiber/v2"
)

type PharmacistController struct {
	Repo    *routes.Repository
	Service *services.PharmacistService
}

func NewPharmacistController(repo *routes.Repository) *PharmacistController {
	service := services.NewPharmacistService(repo.Firestore)
	return &PharmacistController{Repo: repo, Service: service}
}

func (pc *PharmacistController) ActivePrescriptionsHandler(c *fiber.Ctx) error {
	nfcID := c.Params("nfc_id")
	prescriptions, err := pc.Service.GetActivePrescriptions(nfcID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": err.Error()})
	}
	if len(prescriptions) == 0 {
		return c.JSON(fiber.Map{"message": "No active prescriptions found"})
	}
	return c.JSON(prescriptions)
}

func (pc *PharmacistController) DispensePrescriptionHandler(c *fiber.Ctx) error {
	type Request struct {
		TokenID string `json:"token_id"`
	}
	var req Request
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}
	if err := pc.Service.DispensePrescription(req.TokenID); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"message": "Prescription dispensed"})
}
