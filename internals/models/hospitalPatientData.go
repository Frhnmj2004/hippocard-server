package models

import "time"

// HospitalPatientData represents the data returned for one-time hospital access
type HospitalPatientData struct {
	Patient        *User                  `json:"patient"`
	Prescriptions  []*Prescription        `json:"prescriptions"`
	MedicalHistory []*MedicalHistoryEntry `json:"medical_history"`
	AccessTime     time.Time              `json:"access_time"`
	AccessID       string                 `json:"access_id"`
}
