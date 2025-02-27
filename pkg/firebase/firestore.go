// Firestore client setup and queries
package firebase

import (
	"context"
	"log"

	"cloud.google.com/go/firestore" // Updated path
	firebase "firebase.google.com/go"
)

type FirestoreClient struct {
	Client *firestore.Client
}

func NewFirestoreClient(app *firebase.App) (*FirestoreClient, error) {
	// Initialize the Firestore client
	client, err := app.Firestore(context.Background())
	if err != nil {
		log.Printf("Failed to initialize Firestore client: %v", err)
		return nil, err
	}

	return &FirestoreClient{Client: client}, nil
}

func (fc *FirestoreClient) Close() error {
	if fc.Client != nil {
		return fc.Client.Close()
	}
	return nil
}
