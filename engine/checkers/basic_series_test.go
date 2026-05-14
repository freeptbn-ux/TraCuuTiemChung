package checkers

import (
	"testing"
	"time"
	"tracuutiemchung-engine/models"
)

func TestCheckBasicSeries_Phase04(t *testing.T) {
	// Scenario 1: Interval vs Age Constraint (The "MAX" rule)
	// Rule: Min age 2mo, Min interval 1mo.
	// Patient: DOB: 2024-01-01. Dose 1: 2024-02-15 (at 1.5mo - INVALID).
	t.Run("MAX Rule - Invalid First Dose", func(t *testing.T) {
		dob := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
		dose1Date := time.Date(2024, 2, 15, 0, 0, 0, 0, time.UTC)
		
		dosesRequired := 2
		minInterval := 30
		minAgeMonths := 2
		
		rule := &models.VaccineRule{
			DisplayName:             "Test Vaccine",
			DosesRequired:           &dosesRequired,
			MinIntervalDays:         []*int{nil, &minInterval},
			MinAgeMonthsAtFirstDose: &minAgeMonths,
		}

		adminMap := map[string][]models.AdministeredDose{
			"Test Vaccine": {
				{VaccineName: "Test Vaccine", Date: dose1Date},
			},
		}

		results := CheckBasicSeries(rule, adminMap, dob, time.Date(2024, 5, 13, 0, 0, 0, 0, time.UTC))
		
		if len(results) == 0 {
			t.Fatal("Expected 1 missing item")
		}
		
		// Expected: Next Dose 1 scheduled for 2024-03-01 (exactly 2mo)
		expectedDate := time.Date(2024, 3, 1, 0, 0, 0, 0, time.UTC)
		if results[0].EarliestNextDoseDate == nil || !results[0].EarliestNextDoseDate.Equal(expectedDate) {
			t.Errorf("Expected EarliestNextDoseDate %v, got %v", expectedDate, results[0].EarliestNextDoseDate)
		}
		if results[0].StatusTags[0] != "due" || results[0].StatusTags[1] != "eligible" {
			t.Errorf("Expected tags [due, eligible], got %v", results[0].StatusTags)
		}
	})

	// Scenario 2: Valid First Dose, calculate Dose 2
	t.Run("MAX Rule - Valid First Dose", func(t *testing.T) {
		dob := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
		dose1Date := time.Date(2024, 3, 1, 0, 0, 0, 0, time.UTC) // Exactly 2mo
		
		dosesRequired := 2
		minInterval := 31 // 1 month approx
		minAgeMonths := 2
		
		rule := &models.VaccineRule{
			DisplayName:             "Test Vaccine",
			DosesRequired:           &dosesRequired,
			MinIntervalDays:         []*int{nil, &minInterval},
			MinAgeMonthsAtFirstDose: &minAgeMonths,
		}

		adminMap := map[string][]models.AdministeredDose{
			"Test Vaccine": {
				{VaccineName: "Test Vaccine", Date: dose1Date},
			},
		}

		results := CheckBasicSeries(rule, adminMap, dob, time.Date(2024, 5, 13, 0, 0, 0, 0, time.UTC))
		
		if len(results) == 0 {
			t.Fatal("Expected 1 missing item")
		}
		
		// Expected: Dose 2 scheduled for 2024-04-01 (D1 + 31 days)
		expectedDate := dose1Date.AddDate(0, 0, 31)
		if results[0].EarliestNextDoseDate == nil || !results[0].EarliestNextDoseDate.Equal(expectedDate) {
			t.Errorf("Expected EarliestNextDoseDate %v, got %v", expectedDate, results[0].EarliestNextDoseDate)
		}
	})

	// Scenario 3: Valid Dose Counting (Dose 2 too close)
	t.Run("Valid Dose Counting - Interval Too Short", func(t *testing.T) {
		dob := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
		dose1Date := time.Date(2024, 3, 1, 0, 0, 0, 0, time.UTC)
		dose2Date := time.Date(2024, 3, 15, 0, 0, 0, 0, time.UTC) // Only 14 days later
		
		dosesRequired := 3
		minInterval := 30
		minAgeMonths := 2
		
		rule := &models.VaccineRule{
			DisplayName:             "Test Vaccine",
			DosesRequired:           &dosesRequired,
			MinIntervalDays:         []*int{nil, &minInterval, &minInterval},
			MinAgeMonthsAtFirstDose: &minAgeMonths,
		}

		adminMap := map[string][]models.AdministeredDose{
			"Test Vaccine": {
				{VaccineName: "Test Vaccine", Date: dose1Date},
				{VaccineName: "Test Vaccine", Date: dose2Date},
			},
		}

		results := CheckBasicSeries(rule, adminMap, dob, time.Date(2024, 5, 13, 0, 0, 0, 0, time.UTC))
		
		if len(results) == 0 {
			t.Fatal("Expected 1 missing item")
		}
		
		// Expected: valid_doses = 1. Dose 2 scheduled from Dose 1.
		expectedDate := dose1Date.AddDate(0, 0, 30)
		if results[0].EarliestNextDoseDate == nil || !results[0].EarliestNextDoseDate.Equal(expectedDate) {
			t.Errorf("Expected EarliestNextDoseDate %v (from D1), got %v", expectedDate, results[0].EarliestNextDoseDate)
		}
		expectedDesc := "Test Vaccine - Mũi 2 (Cần thêm 2 liều)"
		if results[0].Description != expectedDesc {
			t.Errorf("Expected description '%s', got '%s'", expectedDesc, results[0].Description)
		}
	})

	// Scenario 4: Completion Logic
	t.Run("Completion Logic", func(t *testing.T) {
		dob := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
		dose1Date := time.Date(2024, 3, 1, 0, 0, 0, 0, time.UTC)
		dose2Date := time.Date(2024, 4, 1, 0, 0, 0, 0, time.UTC)
		dose3Date := time.Date(2024, 5, 1, 0, 0, 0, 0, time.UTC)
		
		dosesRequired := 3
		minInterval := 30
		
		rule := &models.VaccineRule{
			DisplayName:     "Test Vaccine",
			DosesRequired:   &dosesRequired,
			MinIntervalDays: []*int{nil, &minInterval, &minInterval},
		}

		adminMap := map[string][]models.AdministeredDose{
			"Test Vaccine": {
				{VaccineName: "Test Vaccine", Date: dose1Date},
				{VaccineName: "Test Vaccine", Date: dose2Date},
				{VaccineName: "Test Vaccine", Date: dose3Date},
			},
		}

		results := CheckBasicSeries(rule, adminMap, dob, time.Date(2024, 5, 13, 0, 0, 0, 0, time.UTC))
		
		if len(results) != 0 {
			t.Errorf("Expected 0 missing items for completed series, got %d", len(results))
		}
	})
}

