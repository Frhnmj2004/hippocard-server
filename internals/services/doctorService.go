package services

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/Frhnmj2004/hippocard-server/internals/models"
	"github.com/Frhnmj2004/hippocard-server/pkg/crypto"
	"github.com/Frhnmj2004/hippocard-server/pkg/firebase"
	"github.com/Frhnmj2004/hippocard-server/pkg/storage"

	//"cloud.google.com/go/firestore"
	"github.com/google/uuid"
)

// DoctorService handles doctor-related operations
type DoctorService struct {
	Firestore *firebase.FirestoreClient
	IPFS      *storage.IPFSClient
}

// NewDoctorService creates a new DoctorService instance
func NewDoctorService(firestore *firebase.FirestoreClient, ipfs *storage.IPFSClient) *DoctorService {
	return &DoctorService{
		Firestore: firestore,
		IPFS:      ipfs,
	}
}

// GetPatientByNFC retrieves a patient’s profile by NFC ID from Firestore
func (ds *DoctorService) GetPatientByNFC(nfcID string) (*models.User, error) {
	ctx := context.Background()

	// Query Firestore for user with matching NFC ID and "patient" role
	docs, err := ds.Firestore.Client.Collection("users").
		Where("nfc_id", "==", nfcID).
		Where("role", "==", "patient").
		Documents(ctx).GetAll()
	if err != nil {
		log.Printf("Failed to query patient by NFC ID: %v", err)
		return nil, err
	}

	if len(docs) == 0 {
		return nil, logError("Patient not found with NFC ID: " + nfcID)
	}

	var user models.User
	if err := docs[0].DataTo(&user); err != nil {
		log.Printf("Failed to parse user data: %v", err)
		return nil, err
	}
	user.UID = docs[0].Ref.ID // Set the Firestore document ID
	return &user, nil
}

// AddMedicalHistory encrypts and stores a patient’s medical history
func (ds *DoctorService) AddMedicalHistory(patientID, history string, key []byte) (string, error) {
	ctx := context.Background()

	// Step 1: Encrypt the medical history
	encryptedData, err := crypto.Encrypt([]byte(history), key)
	if err != nil {
		log.Printf("Failed to encrypt medical history: %v", err)
		return "", err
	}

	// Step 2: Upload encrypted data to IPFS
	cid, err := ds.IPFS.AddData(encryptedData)
	if err != nil {
		log.Printf("Failed to upload to IPFS: %v", err)
		return "", err
	}

	// Step 3: Save CID and metadata to Firestore
	docID := uuid.New().String()
	_, err = ds.Firestore.Client.Collection("medical_history").Doc(docID).Set(ctx, models.MedicalHistory{
		ID:        docID,
		UserID:    patientID,
		CID:       cid,
		CreatedAt: time.Now().UTC(),
	})
	if err != nil {
		log.Printf("Failed to save medical history to Firestore: %v", err)
		return "", err
	}

	return docID, nil
}

// SearchPatients finds patients by name (partial match)
func (ds *DoctorService) SearchPatients(name string) ([]*models.User, error) {
	ctx := context.Background()

	// Query Firestore for patients with matching name (case-insensitive partial match)
	iter := ds.Firestore.Client.Collection("users").
		Where("role", "==", "patient").
		Documents(ctx)

	docs, err := iter.GetAll()
	if err != nil {
		log.Printf("Failed to search patients: %v", err)
		return nil, err
	}

	var results []*models.User
	for _, doc := range docs {
		var user models.User
		if err := doc.DataTo(&user); err != nil {
			log.Printf("Failed to parse user data: %v", err)
			continue
		}
		user.UID = doc.Ref.ID
		if strings.Contains(strings.ToLower(user.Name), strings.ToLower(name)) { // Use user.Name directly
			results = append(results, &user)
		}
	}

	return results, nil
}

// CreatePrescription is a placeholder until blockchain is implemented
func (ds *DoctorService) CreatePrescription(patientID, medication string, dosage uint64) (string, error) {
	// TODO: Implement with blockchain NFT minting
	log.Println("CreatePrescription not implemented yet—waiting for blockchain")
	return "", nil
}

func logError(msg string) error {
	err := fmt.Errorf(msg)
	log.Println(err)
	return err
}
