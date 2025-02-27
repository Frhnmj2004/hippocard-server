package main

import (
	"log"

	"github.com/Frhnmj2004/hippocard-server/api/routes"
	"github.com/Frhnmj2004/hippocard-server/configs"
	"github.com/Frhnmj2004/hippocard-server/pkg/blockchain"
	"github.com/Frhnmj2004/hippocard-server/pkg/firebase"
	"github.com/Frhnmj2004/hippocard-server/pkg/storage"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables from .env
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file: ", err)
	}

	// Load unified config from configs package
	config, err := configs.LoadConfig()
	if err != nil {
		log.Fatal("Failed to load config: ", err)
	}

	// Initialize Firebase Authentication
	authClient, err := firebase.NewAuthClient(&config.Firebase)
	if err != nil {
		log.Fatal("Could not initialize Firebase Auth: ", err)
	}

	// Initialize Firestore
	firestoreClient, err := firebase.NewFirestoreClient(&config.Firebase)
	if err != nil {
		log.Fatal("Could not initialize Firestore: ", err)
	}

	// Initialize Blockchain client (Polygon)
	blockchainClient, err := blockchain.NewClient(&config.Blockchain)
	if err != nil {
		log.Fatal("Could not initialize Blockchain client: ", err)
	}

	// Initialize IPFS client
	ipfsClient, err := storage.NewIPFSClient(&config.IPFS)
	if err != nil {
		log.Fatal("Could not initialize IPFS client: ", err)
	}

	// Create Repository with dependencies
	r := routes.Repository{
		Auth:       authClient,
		Firestore:  firestoreClient,
		Blockchain: blockchainClient,
		IPFS:       ipfsClient,
	}

	// Initialize Fiber app
	app := fiber.New()

	// Setup routes with the Repository
	r.SetupRoutes(app)

	// Start server
	log.Printf("Server starting on :%s", config.ServerPort)
	if err := app.Listen(":" + config.ServerPort); err != nil {
		log.Fatal("Failed to start server: ", err)
	}
}
