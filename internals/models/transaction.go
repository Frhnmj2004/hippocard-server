// Transaction log model
package models

import "time"

// Transaction logs hospital one-time access events
type Transaction struct {
	ID         string    `json:"id"`          // Firestore document ID (UUID)
	UserID     string    `json:"user_id"`     // Hospital’s UID
	PatientID  string    `json:"patient_id"`  // Patient’s UID accessed
	NFCID      string    `json:"nfc_id"`      // Patient’s NFC ID
	AccessTime time.Time `json:"access_time"` // When access occurred
}
