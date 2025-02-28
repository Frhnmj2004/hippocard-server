package controllers

import (
	"github.com/Frhnmj2004/hippocard-server/api/routes"
	"github.com/Frhnmj2004/hippocard-server/internals/services"
	"github.com/gofiber/fiber/v2"
)

// DoctorController holds doctor-specific handlers
type DoctorController struct {
	Repo    *routes.Repository
	Service *services.DoctorService
}

// NewDoctorController creates a new DoctorController
func NewDoctorController(repo *routes.Repository) *DoctorController {
	service := services.NewDoctorService(repo.Firestore, repo.IPFS)
	return &DoctorController{Repo: repo, Service: service}
}

func (dc *DoctorController) GetPatientHandler(c *fiber.Ctx) error {
	nfcID := c.Params("nfc_id")
	patient, err := dc.Service.GetPatientByNFC(nfcID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(patient)
}

func (dc *DoctorController) CreatePrescriptionHandler(c *fiber.Ctx) error {
	type Request struct {
		PatientID  string `json:"patient_id"`
		Medication string `json:"medication"`
		Dosage     uint64 `json:"dosage"`
	}
	var req Request
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}
	prescriptionID, err := dc.Service.CreatePrescription(req.PatientID, req.Medication, req.Dosage)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	if prescriptionID == "" {
		return c.JSON(fiber.Map{"message": "Prescription creation not implemented yet—waiting for blockchain"})
	}
	return c.JSON(fiber.Map{"prescription_id": prescriptionID})
}

func (dc *DoctorController) AddMedicalHistoryHandler(c *fiber.Ctx) error {
	type Request struct {
		PatientID string `json:"patient_id"`
		History   string `json:"history"`
	}
	var req Request
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}
	// Temporary key—replace with real key management
	key := []byte("32-byte-key-here-1234567890123456")
	docID, err := dc.Service.AddMedicalHistory(req.PatientID, req.History, key)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"doc_id": docID})
}

func (dc *DoctorController) SearchPatientsHandler(c *fiber.Ctx) error {
	name := c.Query("name")
	patients, err := dc.Service.SearchPatients(name)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(patients)
}
