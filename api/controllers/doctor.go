package controllers

import (
	"github.com/Frhnmj2004/hippocard-server/api/routes"
	"github.com/gofiber/fiber/v2"
)

// DoctorController holds doctor-specific handlers
type DoctorController struct {
	Repo *routes.Repository
}

// NewDoctorController creates a new DoctorController
func NewDoctorController(repo *routes.Repository) *DoctorController {
	return &DoctorController{Repo: repo}
}

func (dc *DoctorController) GetPatientHandler(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{"message": "Doctor get patient not implemented"})
}

func (dc *DoctorController) CreatePrescriptionHandler(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{"message": "Doctor create prescription not implemented"})
}

func (dc *DoctorController) AddMedicalHistoryHandler(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{"message": "Doctor add medical history not implemented"})
}

func (dc *DoctorController) SearchPatientsHandler(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{"message": "Doctor search patients not implemented"})
}
