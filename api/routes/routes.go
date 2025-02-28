package routes

import (
	"github.com/Frhnmj2004/hippocard-server/api/middleware"
	"github.com/Frhnmj2004/hippocard-server/pkg/blockchain"
	"github.com/Frhnmj2004/hippocard-server/pkg/firebase"
	"github.com/Frhnmj2004/hippocard-server/pkg/storage"

	"github.com/gofiber/fiber/v2"
)

// Repository holds all clients and services for routing
type Repository struct {
	Auth       *firebase.AuthClient
	Firestore  *firebase.FirestoreClient
	Blockchain *blockchain.Client
	IPFS       *storage.IPFSClient
	App        *fiber.App
}

// NewRepository initializes a new Repository
func NewRepository(auth *firebase.AuthClient, firestore *firebase.FirestoreClient, blockchain *blockchain.Client, ipfs *storage.IPFSClient) *Repository {
	return &Repository{
		Auth:       auth,
		Firestore:  firestore,
		Blockchain: blockchain,
		IPFS:       ipfs,
	}
}

// SetupRoutes sets up all API routes using provided handlers
func (r *Repository) SetupRoutes(app *fiber.App, loginHandler func(*fiber.Ctx) error,
	patientProfileHandler func(*fiber.Ctx) error,
	patientPrescriptionsHandler func(*fiber.Ctx) error,
	patientMedicalHistoryHandler func(*fiber.Ctx) error,
	doctorPatientHandler func(*fiber.Ctx) error,
	doctorPrescriptionHandler func(*fiber.Ctx) error,
	doctorMedicalHistoryHandler func(*fiber.Ctx) error,
	doctorSearchPatientsHandler func(*fiber.Ctx) error,
	pharmacistActivePrescriptionsHandler func(*fiber.Ctx) error,
	pharmacistDispenseHandler func(*fiber.Ctx) error,
	hospitalPatientDataHandler func(*fiber.Ctx) error) {
	r.App = app

	// Public routes
	app.Post("/api/login", loginHandler)

	// Patient routes
	patient := app.Group("/api/patient", middleware.AuthMiddleware(r.Auth, "patient"))
	patient.Get("/profile", patientProfileHandler)
	patient.Get("/prescriptions", patientPrescriptionsHandler)
	patient.Get("/medical-history", patientMedicalHistoryHandler)

	// Doctor routes
	doctor := app.Group("/api/doctor", middleware.AuthMiddleware(r.Auth, "doctor"))
	doctor.Get("/patient/:nfc_id", doctorPatientHandler)
	doctor.Post("/prescription", doctorPrescriptionHandler)
	doctor.Post("/medical-history", doctorMedicalHistoryHandler)
	doctor.Get("/patients/search", doctorSearchPatientsHandler)

	// Pharmacist routes
	pharmacist := app.Group("/api/pharmacy", middleware.AuthMiddleware(r.Auth, "pharmacist"))
	pharmacist.Get("/prescriptions/active/:nfc_id", pharmacistActivePrescriptionsHandler)
	pharmacist.Post("/prescription/dispense", pharmacistDispenseHandler)

	// Hospital routes (with one-time access)
	hospital := app.Group("/api/hospital", middleware.AuthMiddleware(r.Auth, "hospital"), middleware.OneTimeAccess(r.Auth, r.Firestore))
	hospital.Get("/patient/:nfc_id", hospitalPatientDataHandler)
}
