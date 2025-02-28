package controllers

import (
	"github.com/Frhnmj2004/hippocard-server/api/routes"
	"github.com/Frhnmj2004/hippocard-server/internals/services"
	"github.com/Frhnmj2004/hippocard-server/pkg/firebase"
	"github.com/Frhnmj2004/hippocard-server/pkg/storage"

	"github.com/gofiber/fiber/v2"
)

// PatientController handles patient-related endpoints
type PatientController struct {
	PatientService *services.PatientService
	AuthClient     *firebase.AuthClient
	Firestore      *firebase.FirestoreClient
	IPFS           *storage.IPFSClient
}

// NewPatientController initializes a new PatientController with the repository
func NewPatientController(repo *routes.Repository) *PatientController {
	return &PatientController{
		PatientService: services.NewPatientService(repo.Firestore, repo.IPFS),
		AuthClient:     repo.Auth,
		Firestore:      repo.Firestore,
		IPFS:           repo.IPFS,
	}
}

// ProfileHandler returns the patient's profile
func (pc *PatientController) ProfileHandler(c *fiber.Ctx) error {
	userID, ok := c.Locals("userID").(string)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}

	user, err := firebase.GetUserByUID(pc.Firestore.Client, userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(user)
}

// PrescriptionsHandler returns the patient's prescriptions (placeholder until blockchain)
func (pc *PatientController) PrescriptionsHandler(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{"message": "Prescriptions not implemented yet—waiting for blockchain"})
}

// MedicalHistoryHandler returns the patient's medical history
func (pc *PatientController) MedicalHistoryHandler(c *fiber.Ctx) error {
	userID, ok := c.Locals("userID").(string)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}

	// Simulate a key for decryption (replace with actual key management)
	key := []byte("your-32-byte-encryption-key-here") // Use a secure key from crypto package
	history, err := pc.PatientService.GetMedicalHistory(userID, key)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(history)
}
