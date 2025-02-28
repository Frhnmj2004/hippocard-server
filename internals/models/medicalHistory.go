package models

import "time"

// MedicalHistory represents a patient’s medical history entry, linked to IPFS
type MedicalHistory struct {
	ID        string    `json:"id"`         // Firestore document ID (UUID)
	UserID    string    `json:"user_id"`    // Patient’s UID
	CID       string    `json:"cid"`        // IPFS Content Identifier for encrypted data
	CreatedAt time.Time `json:"created_at"` // When the history was added
}

type MedicalHistoryEntry struct {
	History   string    `json:"history"`
	CreatedAt time.Time `json:"created_at"`
}
