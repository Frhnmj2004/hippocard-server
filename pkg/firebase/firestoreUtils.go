package firebase

import (
	"context"
	"log"

	"cloud.google.com/go/firestore"
	"github.com/Frhnmj2004/hippocard-server/internals/models"
)

// GetUserByUID fetches a user from Firestore by UID
func GetUserByUID(client *firestore.Client, uid string) (*models.User, error) {
	ctx := context.Background()
	doc, err := client.Collection("users").Doc(uid).Get(ctx)
	if err != nil {
		log.Printf("Failed to get user by UID: %v", err)
		return nil, err
	}

	var user models.User
	if err := doc.DataTo(&user); err != nil {
		log.Printf("Failed to parse user data: %v", err)
		return nil, err
	}
	user.UID = doc.Ref.ID
	return &user, nil
}

// BatchSavePrescriptions saves multiple prescriptions to Firestore
func BatchSavePrescriptions(client *firestore.Client, prescriptions []*models.Prescription) error {
	ctx := context.Background()
	batch := client.Batch()
	for _, p := range prescriptions {
		docRef := client.Collection("prescriptions").Doc(p.ID)
		batch.Set(docRef, p)
	}
	_, err := batch.Commit(ctx)
	if err != nil {
		log.Printf("Failed to batch save prescriptions: %v", err)
		return err
	}
	return nil
}

// GetPrescriptionsByUserID fetches all prescriptions for a user
func GetPrescriptionsByUserID(client *firestore.Client, userID string) ([]*models.Prescription, error) {
	ctx := context.Background()
	docs, err := client.Collection("prescriptions").
		Where("user_id", "==", userID).
		Documents(ctx).GetAll()
	if err != nil {
		log.Printf("Failed to query prescriptions for user %s: %v", userID, err)
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

// SaveMedicalHistory saves a medical history entry to Firestore
func SaveMedicalHistory(client *firestore.Client, history *models.MedicalHistory) error {
	ctx := context.Background()
	_, err := client.Collection("medical_history").Doc(history.ID).Set(ctx, history)
	if err != nil {
		log.Printf("Failed to save medical history for user %s: %v", history.UserID, err)
		return err
	}
	return nil
}

// LogTransaction logs a hospital one-time access event
func LogTransaction(client *firestore.Client, transaction *models.Transaction) error {
	ctx := context.Background()
	_, err := client.Collection("transactions").Doc(transaction.ID).Set(ctx, transaction)
	if err != nil {
		log.Printf("Failed to log transaction for user %s: %v", transaction.UserID, err)
		return err
	}
	return nil
}

// GetTransactionsByUserID fetches all transactions for a user (e.g., hospital)
func GetTransactionsByUserID(client *firestore.Client, userID string) ([]*models.Transaction, error) {
	ctx := context.Background()
	docs, err := client.Collection("transactions").
		Where("user_id", "==", userID).
		Documents(ctx).GetAll()
	if err != nil {
		log.Printf("Failed to query transactions for user %s: %v", userID, err)
		return nil, err
	}

	var transactions []*models.Transaction
	for _, doc := range docs {
		var t models.Transaction
		if err := doc.DataTo(&t); err != nil {
			log.Printf("Failed to parse transaction data: %v", err)
			continue
		}
		t.ID = doc.Ref.ID
		transactions = append(transactions, &t)
	}
	return transactions, nil
}
