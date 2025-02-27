package controllers

import (
	"github.com/Frhnmj2004/hippocard-server/api/routes"
	"github.com/gofiber/fiber/v2"
)

// HospitalController holds hospital-specific handlers
type HospitalController struct {
	Repo *routes.Repository
}

// NewHospitalController creates a new HospitalController
func NewHospitalController(repo *routes.Repository) *HospitalController {
	return &HospitalController{Repo: repo}
}

func (hc *HospitalController) PatientDataHandler(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{"message": "Hospital patient data not implemented"})
}
