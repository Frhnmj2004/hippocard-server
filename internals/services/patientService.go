package services

import (
	"context"
	"log"

	//"strings"

	"github.com/Frhnmj2004/hippocard-server/internals/models"
	"github.com/Frhnmj2004/hippocard-server/pkg/crypto"
	"github.com/Frhnmj2004/hippocard-server/pkg/firebase"
	"github.com/Frhnmj2004/hippocard-server/pkg/storage"
	//"cloud.google.com/go/firestore"
)

type PatientService struct {
	Firestore *firebase.FirestoreClient
	IPFS      *storage.IPFSClient
}

func NewPatientService(firestore *firebase.FirestoreClient, ipfs *storage.IPFSClient) *PatientService {
	return &PatientService{
		Firestore: firestore,
		IPFS:      ipfs,
	}
}

func (ps *PatientService) GetProfile(userID string) (*models.User, error) {
	ctx := context.Background()

	// Fetch user document from Firestore by UID
	doc, err := ps.Firestore.Client.Collection("users").Doc(userID).Get(ctx)
	if err != nil {
		log.Printf("Failed to get patient profile: %v", err)
		return nil, err
	}

	var user models.User
	if err := doc.DataTo(&user); err != nil {
		log.Printf("Failed to parse user data: %v", err)
		return nil, err
	}
	user.UID = doc.Ref.ID

	if user.Role != "patient" {
		return nil, logError("User is not a patient: " + userID)
	}

	return &user, nil
}

func (ps *PatientService) GetMedicalHistory(userID string, key []byte) ([]*models.MedicalHistoryEntry, error) {
	ctx := context.Background()

	// Query Firestore for medical history entries
	docs, err := ps.Firestore.Client.Collection("medical_history").
		Where("user_id", "==", userID).
		Documents(ctx).GetAll()
	if err != nil {
		log.Printf("Failed to query medical history: %v", err)
		return nil, err
	}

	var history []*models.MedicalHistoryEntry
	for _, doc := range docs {
		var mh models.MedicalHistory
		if err := doc.DataTo(&mh); err != nil {
			log.Printf("Failed to parse medical history: %v", err)
			continue
		}
		mh.ID = doc.Ref.ID

		// Fetch encrypted data from IPFS
		encryptedData, err := ps.IPFS.GetData(mh.CID)
		if err != nil {
			log.Printf("Failed to fetch from IPFS for CID %s: %v", mh.CID, err)
			continue
		}

		// Decrypt the data
		decryptedData, err := crypto.Decrypt(encryptedData, key)
		if err != nil {
			log.Printf("Failed to decrypt history for CID %s: %v", mh.CID, err)
			continue
		}

		// Add to results
		history = append(history, &models.MedicalHistoryEntry{
			History:   string(decryptedData),
			CreatedAt: mh.CreatedAt,
		})
	}

	return history, nil
}

func (ps *PatientService) GetPrescriptions(userID string) ([]*models.Prescription, error) {
	// TODO: Implement with blockchain NFT data
	log.Println("GetPrescriptions not implemented yet—waiting for blockchain")
	return nil, nil
}
