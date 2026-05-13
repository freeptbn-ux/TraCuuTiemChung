package engine

import (
	"testing"
	"time"
	"tracuutiemchung-engine/models"
)

func TestPhase07_LiveVaccineSpacing(t *testing.T) {
	// Scenario: MMR and Varicella are both due on 2024-06-01.
	// Since they are both live, if they are not administered on the same day,
	// they must be spaced by at least 28 days.
	// In the engine, if we recommend both on the same day, that's technically allowed
	// IF the user can actually give them both. 
	// But the post-processor logic says: 
	// "If not on the same day as last live dose... push to 28 days later"
	// Wait, if they are BOTH recommendations for the FUTURE, how do we handle them?
	// The current post-processor sorts them and then spaces them.
	
	rules := map[string]models.VaccineRule{
		"MMR": {
			DisplayName: "MMR (Live)",
			IsLive:      true,
		},
		"Varivax": {
			DisplayName: "Varivax (Live)",
			IsLive:      true,
		},
	}

	date1 := time.Date(2024, 6, 1, 0, 0, 0, 0, time.UTC)
	results := []models.MissingItem{
		{VaccineName: "MMR (Live)", EarliestNextDoseDate: &date1},
		{VaccineName: "Varivax (Live)", EarliestNextDoseDate: &date1},
	}

	record := models.PatientRecord{
		AdministeredMap: make(map[string][]models.AdministeredDose),
	}

	finalResults := ApplySpacingAndSort(results, rules, record)
	if len(finalResults) != 2 {
		t.Errorf("Expected 2 results, got %d", len(finalResults))
	}

	// One should stay at 2024-06-01, the other should stay at 2024-06-01?
	// Let's check the implementation of ApplySpacingAndSort again.
	/*
	if !isSameDay(currentDate, lastLiveDate) {
		minAllowedDate := lastLiveDate.AddDate(0, 0, 28)
		if currentDate.Before(minAllowedDate) {
			currentDate = minAllowedDate
			results[i].EarliestNextDoseDate = &currentDate
			results[i].StatusTags = append(results[i].StatusTags, "spacing_adjusted")
		}
	}
	*/
	// If both are 2024-06-01, isSameDay will be true. So it won't push.
	// Is this correct? In practice, yes, you can give live vaccines on the SAME day.
	// If you miss that day, you wait 28 days.
	
	// If one was 2024-06-01 and the other was 2024-06-05.
	date2 := time.Date(2024, 6, 5, 0, 0, 0, 0, time.UTC)
	results2 := []models.MissingItem{
		{VaccineName: "MMR (Live)", EarliestNextDoseDate: &date1},
		{VaccineName: "Varivax (Live)", EarliestNextDoseDate: &date2},
	}
	
	finalResults2 := ApplySpacingAndSort(results2, rules, record)
	
	// Varivax should be pushed to 2024-06-29 (28 days after 2024-06-01).
	foundVarivax := false
	expectedDate := time.Date(2024, 6, 29, 0, 0, 0, 0, time.UTC)
	for _, item := range finalResults2 {
		if item.VaccineName == "Varivax (Live)" {
			foundVarivax = true
			if item.EarliestNextDoseDate == nil || !item.EarliestNextDoseDate.Equal(expectedDate) {
				t.Errorf("Expected Varivax at %v, got %v", expectedDate, item.EarliestNextDoseDate)
			}
		}
	}
	if !foundVarivax {
		t.Error("Varivax recommendation not found")
	}
}

func TestPhase07_OutputSorting(t *testing.T) {
	// Scenario: Varicella (2024-10-01), BCG (2024-01-01), 6-in-1 (2024-02-01)
	dateV := time.Date(2024, 10, 1, 0, 0, 0, 0, time.UTC)
	dateB := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	date6 := time.Date(2024, 2, 1, 0, 0, 0, 0, time.UTC)

	results := []models.MissingItem{
		{VaccineName: "Varicella", EarliestNextDoseDate: &dateV},
		{VaccineName: "BCG", EarliestNextDoseDate: &dateB},
		{VaccineName: "6-in-1", EarliestNextDoseDate: &date6},
	}

	rules := make(map[string]models.VaccineRule)
	record := models.PatientRecord{}

	finalResults := ApplySpacingAndSort(results, rules, record)

	if len(finalResults) != 3 {
		t.Fatalf("Expected 3 results, got %d", len(finalResults))
	}

	if finalResults[0].VaccineName != "BCG" {
		t.Errorf("Expected 1st result to be BCG, got %s", finalResults[0].VaccineName)
	}
	if finalResults[1].VaccineName != "6-in-1" {
		t.Errorf("Expected 2nd result to be 6-in-1, got %s", finalResults[1].VaccineName)
	}
	if finalResults[2].VaccineName != "Varicella" {
		t.Errorf("Expected 3rd result to be Varicella, got %s", finalResults[2].VaccineName)
	}
}

func TestPhase07_RuleDispatcher(t *testing.T) {
	// Ensure all rule types are handled without crashing
	rules := map[string]models.VaccineRule{
		"R1": {Type: "single_series", DisplayName: "Basic"},
		"R2": {Type: "age_dependent_series", DisplayName: "AgeDep"},
		"R3": {Type: "mmr_equivalent_group", DisplayName: "MMR"},
		"R4": {Type: "flu_group", DisplayName: "Flu"},
		"R5": {Type: "unknown_type", DisplayName: "Unknown"},
	}

	dob := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	record := models.PatientRecord{
		BirthDate: dob,
		AdministeredMap: make(map[string][]models.AdministeredDose),
	}

	// Should not panic
	results := ProcessAllRules(rules, record, dob.AddDate(0, 3, 0))
	
	// R5 should probably not produce a result if it's unknown and we don't have a default
	// Let's check how many results we got.
	// Flu will produce a result at 6 months, but we are at 3 months.
	// Basic might produce a result if min_age is met.
	t.Logf("Got %d results", len(results))
}
