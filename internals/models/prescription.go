package models

import "time"

// Prescription represents a medical prescription stored as an NFT
type Prescription struct {
	ID          string     `json:"id"`                     // Firestore document ID (token ID or UUID)
	UserID      string     `json:"user_id"`                // Patient’s UID
	TokenID     string     `json:"token_id"`               // NFT token ID on Polygon (string for simplicity)
	Medication  string     `json:"medication"`             // Medication name (e.g., "Aspirin")
	Dosage      uint64     `json:"dosage"`                 // Dosage amount (e.g., 200 mg)
	IsActive    bool       `json:"is_active"`              // Whether the prescription is still active
	CreatedAt   time.Time  `json:"created_at"`             // When the prescription was created
	DispensedAt *time.Time `json:"dispensed_at,omitempty"` // When dispensed (null if active)
}
