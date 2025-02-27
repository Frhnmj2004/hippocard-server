package controllers

import (
	"github.com/Frhnmj2004/hippocard-server/api/routes"
	"github.com/gofiber/fiber/v2"
)

type PatientController struct {
	Repo *routes.Repository
}

func NewPatientController(repo *routes.Repository) *PatientController {
	return &PatientController{Repo: repo}
}

func (pc *PatientController) ProfileHandler(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{"message": "Patient profile not implemented"})
}

func (pc *PatientController) PrescriptionsHandler(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{"message": "Patient prescriptions not implemented"})
}

func (pc *PatientController) MedicalHistoryHandler(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{"message": "Patient medical history not implemented"})
}
