package services

import (
	"context"
	"log"
	"time"

	"github.com/Frhnmj2004/hippocard-server/internals/models"
	"github.com/Frhnmj2004/hippocard-server/pkg/firebase"

	"cloud.google.com/go/firestore"
)

type PharmacistService struct {
	Firestore *firebase.FirestoreClient
}

func NewPharmacistService(firestore *firebase.FirestoreClient) *PharmacistService {
	return &PharmacistService{
		Firestore: firestore,
	}
}

func (ps *PharmacistService) GetActivePrescriptions(nfcID string) ([]*models.Prescription, error) {
	ctx := context.Background()

	// Step 1: Find patient by NFC ID
	userDocs, err := ps.Firestore.Client.Collection("users").
		Where("nfc_id", "==", nfcID).
		Where("role", "==", "patient").
		Documents(ctx).GetAll()
	if err != nil {
		log.Printf("Failed to query patient by NFC ID: %v", err)
		return nil, err
	}
	if len(userDocs) == 0 {
		return nil, logError("Patient not found with NFC ID: " + nfcID)
	}

	var patient models.User
	if err := userDocs[0].DataTo(&patient); err != nil {
		log.Printf("Failed to parse patient data: %v", err)
		return nil, err
	}
	patient.UID = userDocs[0].Ref.ID

	// Step 2: Query active prescriptions
	docs, err := ps.Firestore.Client.Collection("prescriptions").
		Where("user_id", "==", patient.UID).
		Where("is_active", "==", true).
		Documents(ctx).GetAll()
	if err != nil {
		log.Printf("Failed to query active prescriptions: %v", err)
		return nil, err
	}

	var prescriptions []*models.Prescription
	for _, doc := range docs {
		var p models.Prescription
		if err := doc.DataTo(&p); err != nil {
			log.Printf("Failed to parse prescription data: %v", err)
			continue
		}
		p.ID = doc.Ref.ID
		prescriptions = append(prescriptions, &p)
	}

	return prescriptions, nil
}

func (ps *PharmacistService) DispensePrescription(tokenID string) error {
	ctx := context.Background()

	// Update Firestore (temporary—blockchain burning TBD)
	_, err := ps.Firestore.Client.Collection("prescriptions").Doc(tokenID).Update(ctx, []firestore.Update{
		{Path: "is_active", Value: false},
		{Path: "dispensed_at", Value: time.Now().UTC()},
	})
	if err != nil {
		log.Printf("Failed to update prescription status: %v", err)
		return err
	}

	// TODO: Add blockchain NFT burning logic here
	log.Println("DispensePrescription blockchain logic not implemented yet")
	return nil
}
