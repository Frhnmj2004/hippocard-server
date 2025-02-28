package main

import (
	"context"
	"log"

	"github.com/joho/godotenv"

	//"os"

	"github.com/Frhnmj2004/hippocard-server/api/controllers"
	"github.com/Frhnmj2004/hippocard-server/api/routes"
	"github.com/Frhnmj2004/hippocard-server/configs"
	"github.com/Frhnmj2004/hippocard-server/pkg/blockchain"
	"github.com/Frhnmj2004/hippocard-server/pkg/firebase"
	"github.com/Frhnmj2004/hippocard-server/pkg/storage"

	firebaseLib "firebase.google.com/go"
	"github.com/gofiber/fiber/v2"
	"google.golang.org/api/option"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file: ", err)
	}

	config, err := configs.LoadConfig()
	if err != nil {
		log.Fatal("Failed to load config: ", err)
	}

	// Initialize Firebase App
	firebaseApp, err := firebaseLib.NewApp(context.Background(), nil, option.WithCredentialsFile(config.Firebase.CredentialsPath))
	if err != nil {
		log.Fatal("Failed to initialize Firebase app: ", err)
	}

	// Initialize Firebase Auth and Firestore clients
	authClient, err := firebase.NewAuthClient(firebaseApp)
	if err != nil {
		log.Fatal("Could not initialize Firebase Auth: ", err)
	}

	firestoreClient, err := firebase.NewFirestoreClient(firebaseApp)
	if err != nil {
		log.Fatal("Could not initialize Firestore: ", err)
	}

	// Initialize blockchain and IPFS clients
	blockchainClient, err := blockchain.NewClient(config)
	if err != nil {
		log.Fatal("Could not initialize Blockchain client: ", err)
	}

	ipfsClient, err := storage.NewIPFSClient(config)
	if err != nil {
		log.Fatal("Could not initialize IPFS client: ", err)
	}

	// Set up routes with repository and custom handlers
	r := routes.NewRepository(authClient, firestoreClient, blockchainClient, ipfsClient)
	app := fiber.New()

	// Create controllers and get handlers
	authController := controllers.NewAuthController(authClient)
	patientController := controllers.NewPatientController(r)
	doctorController := controllers.NewDoctorController(r)
	pharmacistController := controllers.NewPharmacistController(r)
	hospitalController := controllers.NewHospitalController(r)

	// Define handlers
	loginHandler := authController.LoginHandler
	patientProfileHandler := patientController.ProfileHandler
	patientPrescriptionsHandler := patientController.PrescriptionsHandler
	patientMedicalHistoryHandler := patientController.MedicalHistoryHandler
	doctorPatientHandler := doctorController.GetPatientHandler
	doctorPrescriptionHandler := doctorController.CreatePrescriptionHandler
	doctorMedicalHistoryHandler := doctorController.AddMedicalHistoryHandler
	doctorSearchPatientsHandler := doctorController.SearchPatientsHandler
	pharmacistActivePrescriptionsHandler := pharmacistController.ActivePrescriptionsHandler
	pharmacistDispenseHandler := pharmacistController.DispensePrescriptionHandler
	hospitalPatientDataHandler := hospitalController.PatientDataHandler

	// Set up routes with all handlers
	r.SetupRoutes(app,
		loginHandler,
		patientProfileHandler,
		patientPrescriptionsHandler,
		patientMedicalHistoryHandler,
		doctorPatientHandler,
		doctorPrescriptionHandler,
		doctorMedicalHistoryHandler,
		doctorSearchPatientsHandler,
		pharmacistActivePrescriptionsHandler,
		pharmacistDispenseHandler,
		hospitalPatientDataHandler,
	)

	log.Printf("Server starting on :%s", config.ServerPort)
	if err := app.Listen(":" + config.ServerPort); err != nil {
		log.Fatal("Failed to start server: ", err)
	}
}
