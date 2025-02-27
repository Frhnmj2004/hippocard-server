package routes

import (
	"log"

	"github.com/Frhnmj2004/hippocard-server/api/controllers"
	"github.com/Frhnmj2004/hippocard-server/pkg/blockchain"

	"github.com/Frhnmj2004/hippocard-server/pkg/firebase"
	"github.com/Frhnmj2004/hippocard-server/pkg/storage"

	"github.com/gofiber/fiber/v2"
)

// Repository holds all service clients for routing
type Repository struct {
	Auth       *firebase.AuthClient
	Firestore  *firebase.FirestoreClient
	Blockchain *blockchain.Client
	IPFS       *storage.IPFSClient
}

// NewRepository creates a new Repository instance
func NewRepository(auth *firebase.AuthClient, firestore *firebase.FirestoreClient, blockchain *blockchain.Client, ipfs *storage.IPFSClient) *Repository {
	return &Repository{
		Auth:       auth,
		Firestore:  firestore,
		Blockchain: blockchain,
		IPFS:       ipfs,
	}
}

// SetupRoutes configures all API routes with Fiber
func (r *Repository) SetupRoutes(app *fiber.App) {
	// Initialize controllers
	authCtrl := controllers.NewAuthController(r)
	patientCtrl := controllers.NewPatientController(r)
	doctorCtrl := controllers.NewDoctorController(r)
	pharmacistCtrl := controllers.NewPharmacistController(r)
	hospitalCtrl := controllers.NewHospitalController(r)

	// Public routes (e.g., login)
	app.Post("/api/login", authCtrl.LoginHandler) // Placeholder for auth logic

	// Patient routes (authenticated)
	patient := app.Group("/api/patient", r.AuthMiddleware("patient"))
	patient.Get("/profile", patientCtrl.ProfileHandler) // Placeholder for
	patient.Get("/prescriptions", patientCtrl.PrescriptionsHandler)
	patient.Get("/medical-history", patientCtrl.MedicalHistoryHandler)

	// Doctor routes (authenticated)
	doctor := app.Group("/api/doctor", r.AuthMiddleware("doctor"))
	doctor.Get("/patient/:nfc_id", doctorCtrl.GetPatientHandler)
	doctor.Post("/prescription", doctorCtrl.CreatePrescriptionHandler)
	doctor.Post("/medical-history", doctorCtrl.AddMedicalHistoryHandler)
	doctor.Get("/patients/search", doctorCtrl.SearchPatientsHandler)

	// Pharmacist routes (authenticated)
	pharmacist := app.Group("/api/pharmacy", r.AuthMiddleware("pharmacist"))
	pharmacist.Get("/prescriptions/active/:nfc_id", pharmacistCtrl.ActivePrescriptionsHandler)
	pharmacist.Post("/prescription/dispense", pharmacistCtrl.DispensePrescriptionHandler)

	// Hospital routes (authenticated, one-time access)
	hospital := app.Group("/api/hospital", r.AuthMiddleware("hospital"))
	hospital.Get("/patient/:nfc_id", hospitalCtrl.PatientDataHandler) // One-time access logic TBD
}

func (r *Repository) AuthMiddleware(role string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		log.Println("Auth middleware not implemented yet for role:", role)
		return c.Next()
	}
}
