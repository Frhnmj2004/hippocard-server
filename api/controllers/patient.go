package controllers

import (
	"github.com/Frhnmj2004/hippocard-server/api/routes"
	"github.com/Frhnmj2004/hippocard-server/internals/services"
	"github.com/gofiber/fiber/v2"
)

type PatientController struct {
	Repo    *routes.Repository
	Service *services.PatientService
}

func NewPatientController(repo *routes.Repository) *PatientController {
	service := services.NewPatientService(repo.Firestore, repo.IPFS)
	return &PatientController{Repo: repo, Service: service}
}

func (pc *PatientController) ProfileHandler(c *fiber.Ctx) error {
	userID := c.Locals("userID").(string)
	profile, err := pc.Service.GetProfile(userID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(profile)
}

func (pc *PatientController) PrescriptionsHandler(c *fiber.Ctx) error {
	userID := c.Locals("userID").(string)
	prescriptions, err := pc.Service.GetPrescriptions(userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	if prescriptions == nil {
		return c.JSON(fiber.Map{"message": "Prescriptions not implemented yet—waiting for blockchain"})
	}
	return c.JSON(prescriptions)
}
func (pc *PatientController) MedicalHistoryHandler(c *fiber.Ctx) error {
	userID := c.Locals("userID").(string)
	// Temporary key—replace with real key management
	key := []byte("32-byte-key-here-1234567890123456")
	history, err := pc.Service.GetMedicalHistory(userID, key)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(history)
}
