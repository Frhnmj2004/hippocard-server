package models

import "time"

// User represents a user in the system (patient, doctor, pharmacist, hospital)
type User struct {
	UID           string    `json:"uid"`                      // Firestore document ID (Firebase UID)
	NFCID         string    `json:"nfc_id"`                   // Unique NFC card identifier
	Name          string    `json:"name"`                     // User’s full name
	Role          string    `json:"role"`                     // Role: "patient", "doctor", "pharmacist", "hospital"
	WalletAddress string    `json:"wallet_address,omitempty"` // Optional for NFT interactions
	CreatedAt     time.Time `json:"created_at"`               // When the user was registered
}
