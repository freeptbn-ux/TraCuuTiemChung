package models

import (
	"sort"
	"time"
)

// AdministeredDose represents a single dose that has been given to the patient
type AdministeredDose struct {
	VaccineName string    `json:"vaccine_name"`
	Date        time.Time `json:"date"`
}

// PatientRecord contains the patient's information and vaccination history
type PatientRecord struct {
	BirthDate       time.Time                     `json:"birth_date"`
	AdministeredMap map[string][]AdministeredDose `json:"administered_map"` // Key: Rule ID or Group Name
}

// GetSortedDoses returns doses for a specific vaccine, sorted by date
func (p *PatientRecord) GetSortedDoses(ruleID string) []AdministeredDose {
	doses, ok := p.AdministeredMap[ruleID]
	if !ok {
		return nil
	}

	// Create a copy to avoid modifying the original map's slice if needed, 
	// but usually we want to sort the existing one or return a sorted copy.
	// The requirement says "sort theo ngày", so we should probably keep them sorted.
	sortedDoses := make([]AdministeredDose, len(doses))
	copy(sortedDoses, doses)

	sort.Slice(sortedDoses, func(i, j int) bool {
		return sortedDoses[i].Date.Before(sortedDoses[j].Date)
	})

	return sortedDoses
}
