package models

import "time"

// Patient represents basic patient information from search results
type Patient struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	DOB    string `json:"dob"`
	Gender string `json:"gender"`
	Code   string `json:"code"`
}

// VaccineRecord represents a single vaccination entry
type VaccineRecord struct {
	VaccineName string    `json:"vaccine_name"`
	Dose        string    `json:"dose"`
	Date        time.Time `json:"date"`
}

// PatientInfo contains details about a patient extracted from their detail page
type PatientInfo struct {
	Name       string `json:"name"`
	Birth      string `json:"birth"`
	SystemDate string `json:"system_date"`
}

// PatientDetail combines patient info and vaccination history
type PatientDetail struct {
	PatientInfo PatientInfo     `json:"patient_info"`
	History     []VaccineRecord `json:"history"`
}
