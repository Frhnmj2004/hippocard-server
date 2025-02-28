package services

import (
	"context"
	"log"
	"time"

	"github.com/Frhnmj2004/hippocard-server/internals/models"
	"github.com/Frhnmj2004/hippocard-server/pkg/crypto"
	"github.com/Frhnmj2004/hippocard-server/pkg/firebase"
	"github.com/Frhnmj2004/hippocard-server/pkg/storage"

	//"cloud.google.com/go/firestore"
	"github.com/google/uuid"
)

type HospitalService struct {
	Firestore *firebase.FirestoreClient
	IPFS      *storage.IPFSClient
}

func NewHospitalService(firestore *firebase.FirestoreClient, ipfs *storage.IPFSClient) *HospitalService {
	return &HospitalService{
		Firestore: firestore,
		IPFS:      ipfs,
	}
}

func (hs *HospitalService) GetPatientData(nfcID string, key []byte) (*models.HospitalPatientData, error) {
	ctx := context.Background()

	// Step 1: Find patient by NFC ID
	userDocs, err := hs.Firestore.Client.Collection("users").
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

	// Step 2: Fetch prescriptions (placeholder until blockchain)
	prescriptions, err := hs.getPrescriptions(patient.UID)
	if err != nil {
		log.Printf("Failed to fetch prescriptions: %v", err)
		return nil, err
	}

	// Step 3: Fetch medical history
	medicalHistory, err := hs.getMedicalHistory(patient.UID, key)
	if err != nil {
		log.Printf("Failed to fetch medical history: %v", err)
		return nil, err
	}

	// Step 4: Prepare response for one-time access
	result := &models.HospitalPatientData{
		Patient:        &patient,
		Prescriptions:  prescriptions,
		MedicalHistory: medicalHistory,
		AccessTime:     time.Now().UTC(),
		AccessID:       uuid.New().String(),
	}

	// Step 5: Simulate one-time access by logging
	log.Println("One-time access granted for patient:", nfcID)

	return result, nil
}

func (hs *HospitalService) getPrescriptions(userID string) ([]*models.Prescription, error) {
	// TODO: Implement with blockchain NFT data
	log.Println("GetPrescriptions not implemented yet—waiting for blockchain")
	return nil, nil
}

func (hs *HospitalService) getMedicalHistory(userID string, key []byte) ([]*models.MedicalHistoryEntry, error) {
	ctx := context.Background()

	docs, err := hs.Firestore.Client.Collection("medical_history").
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

		encryptedData, err := hs.IPFS.GetData(mh.CID)
		if err != nil {
			log.Printf("Failed to fetch from IPFS for CID %s: %v", mh.CID, err)
			continue
		}

		decryptedData, err := crypto.Decrypt(encryptedData, key)
		if err != nil {
			log.Printf("Failed to decrypt history for CID %s: %v", mh.CID, err)
			continue
		}

		history = append(history, &models.MedicalHistoryEntry{
			History:   string(decryptedData),
			CreatedAt: mh.CreatedAt,
		})
	}

	return history, nil
}
