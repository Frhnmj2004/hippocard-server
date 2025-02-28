package controllers

import (
	"github.com/Frhnmj2004/hippocard-server/api/routes"
	"github.com/Frhnmj2004/hippocard-server/internals/services"

	"github.com/gofiber/fiber/v2"
)

type HospitalController struct {
	Repo    *routes.Repository
	Service *services.HospitalService
}

func NewHospitalController(repo *routes.Repository) *HospitalController {
	service := services.NewHospitalService(repo.Firestore, repo.IPFS)
	return &HospitalController{Repo: repo, Service: service}
}

func (hc *HospitalController) PatientDataHandler(c *fiber.Ctx) error {
	nfcID := c.Params("nfc_id")
	// Temporary key—replace with real key management
	key := []byte("32-byte-key-here-1234567890123456")
	data, err := hc.Service.GetPatientData(nfcID, key)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(data)
}
