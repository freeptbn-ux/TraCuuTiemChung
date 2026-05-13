package checkers

import (
	"testing"
	"time"
	"tracuutiemchung-engine/models"
)

func TestAgeDependentSeries_Prevenar(t *testing.T) {
	dob := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	analysisDate := time.Date(2024, 9, 1, 0, 0, 0, 0, time.UTC)

	// Prevenar 13 Rule Mock
	rule := &models.VaccineRule{
		DisplayName: "Prevenar 13 (Phế cầu)",
		Type:        "age_dependent_series",
		RulesByAge: []models.AgeRule{
			{
				MaxAgeAtFirstDoseMonths: intPtr(6),
				DosesRequired:           4,
				MinIntervalDays:         []*int{nil, intPtr(30), intPtr(30), intPtr(240)},
			},
			{
				MinAgeAtFirstDoseMonths: intPtr(7),
				MaxAgeAtFirstDoseMonths: intPtr(11),
				DosesRequired:           3,
				MinIntervalDays:         []*int{nil, intPtr(30), intPtr(180)},
			},
		},
	}

	// Case 1: First dose at 7 months (2024-08-01)
	// Patient is 7 months old at first dose.
	adminMap := map[string][]models.AdministeredDose{
		"Prevenar 13 (Phế cầu)": {
			{VaccineName: "Prevenar 13", Date: time.Date(2024, 8, 1, 0, 0, 0, 0, time.UTC)},
		},
	}

	results := CheckAgeDependentSeries(rule, adminMap, dob, analysisDate)
	if len(results) == 0 {
		t.Fatal("Expected 1 recommendation for Case 1, got 0")
	}

	// 7-11 months rule requires 3 doses. 1 dose given, 2 remaining.
	expectedDesc1 := "Prevenar 13 (Phế cầu) - Mũi 2 (Cần thêm 2 liều)"
	if results[0].Description != expectedDesc1 {
		t.Errorf("Case 1: Expected description '%s', got '%s'", expectedDesc1, results[0].Description)
	}

	// Case 2: First dose at 2 months (2024-03-01)
	// Patient is 2 months old at first dose.
	adminMap2 := map[string][]models.AdministeredDose{
		"Prevenar 13 (Phế cầu)": {
			{VaccineName: "Prevenar 13", Date: time.Date(2024, 3, 1, 0, 0, 0, 0, time.UTC)},
		},
	}

	results2 := CheckAgeDependentSeries(rule, adminMap2, dob, analysisDate)
	if len(results2) == 0 {
		t.Fatal("Expected 1 recommendation for Case 2, got 0")
	}

	// 0-6 months rule requires 4 doses. 1 dose given, 3 remaining.
	expectedDesc2 := "Prevenar 13 (Phế cầu) - Mũi 2 (Cần thêm 3 liều)"
	if results2[0].Description != expectedDesc2 {
		t.Errorf("Case 2: Expected description '%s', got '%s'", expectedDesc2, results2[0].Description)
	}
}

func TestBoosterLogic(t *testing.T) {
	dob := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
	// Analysis date is 4 years after the last dose
	analysisDate := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)

	// Mocking a course with booster
	doses := []models.AdministeredDose{
		{VaccineName: "Jevax", Date: time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC)},
		{VaccineName: "Jevax", Date: time.Date(2021, 1, 8, 0, 0, 0, 0, time.UTC)},
		{VaccineName: "Jevax", Date: time.Date(2022, 1, 15, 0, 0, 0, 0, time.UTC)},
	}

	// DosesRequired: 3, BoosterIntervalYears: 3
	results := CheckSeriesInternal(
		"Jevax/VNNB",
		3,
		[]*int{nil, intPtr(7), intPtr(365)},
		nil, nil, nil, nil, // Min age 12 months, etc.
		nil, intPtr(3), // BoosterIntervalYears: 3
		doses,
		dob,
		analysisDate,
	)

	if len(results) == 0 {
		t.Fatal("Expected booster recommendation, got 0")
	}

	expectedDesc := "Jevax/VNNB - Cần tiêm nhắc lại"
	if results[0].Description != expectedDesc {
		t.Errorf("Expected description '%s', got '%s'", expectedDesc, results[0].Description)
	}
}

func TestComplexInteractions(t *testing.T) {
	adminMap := map[string][]models.AdministeredDose{
		"Vắc xin Sởi-Quai bị-Rubella (MMR-II/Priorix)": {
			{VaccineName: "MMR-II", Date: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)},
		},
	}

	results := []models.MissingItem{
		{VaccineName: "MVVAC (Sởi đơn)", Description: "Cần tiêm Sởi"},
		{VaccineName: "VA - MENGOC - BC (Não mô cầu BC)", Description: "Cần tiêm BC"},
	}

	// Add ACYW to adminMap for second interaction
	adminMap["Vắc xin Não mô cầu ACYW-135 (Menactra/MenQuadfi)"] = []models.AdministeredDose{
		{VaccineName: "MenQuadfi", Date: time.Date(2024, 2, 1, 0, 0, 0, 0, time.UTC)},
	}

	filtered := ApplyComplexInteractions(results, adminMap)

	// 1. MVVAC should be removed because MMR is present
	// 2. BC should be removed because MenQuadfi is present
	
	if len(filtered) != 0 {
		t.Errorf("Expected 0 results after filtering, got %d", len(filtered))
	}
}

func TestBoosterFuture(t *testing.T) {
	dob := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
	// Analysis date is only 2 years after last dose (requires 3 years)
	analysisDate := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)

	doses := []models.AdministeredDose{
		{VaccineName: "Jevax", Date: time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC)},
		{VaccineName: "Jevax", Date: time.Date(2021, 1, 8, 0, 0, 0, 0, time.UTC)},
		{VaccineName: "Jevax", Date: time.Date(2022, 1, 15, 0, 0, 0, 0, time.UTC)},
	}

	results := CheckSeriesInternal(
		"Jevax/VNNB",
		3,
		[]*int{nil, intPtr(7), intPtr(365)},
		intPtr(12), nil, nil, nil,
		nil, intPtr(3), // BoosterIntervalYears: 3
		doses,
		dob,
		analysisDate,
	)

	if len(results) == 0 {
		t.Fatal("Expected booster recommendation (future), got 0")
	}

	// Verify status tag contains "future"
	hasFuture := false
	for _, tag := range results[0].StatusTags {
		if tag == "future" {
			hasFuture = true
			break
		}
	}
	if !hasFuture {
		t.Errorf("Expected 'future' tag for booster not yet due, got %v", results[0].StatusTags)
	}
}

func TestBoosterObjectLogic(t *testing.T) {
	dob := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	analysisDate := time.Date(2025, 2, 1, 0, 0, 0, 0, time.UTC) // 13 months old

	// MenQuadfi Rule Mock for 6w-5m
	booster := &models.BoosterRule{
		MinAgeMonths:            12,
		MinIntervalDaysFromLast: 60,
		Description:             "Mũi nhắc: từ 12 tháng tuổi, cách mũi 3 ít nhất 2 tháng",
	}

	doses := []models.AdministeredDose{
		{VaccineName: "MenQuadfi", Date: time.Date(2024, 2, 15, 0, 0, 0, 0, time.UTC)}, // ~6 weeks
		{VaccineName: "MenQuadfi", Date: time.Date(2024, 4, 15, 0, 0, 0, 0, time.UTC)}, // +60 days approx
		{VaccineName: "MenQuadfi", Date: time.Date(2024, 6, 15, 0, 0, 0, 0, time.UTC)}, // +60 days approx (total 3 doses)
	}

	results := CheckSeriesInternal(
		"MenQuadfi",
		3,
		[]*int{nil, intPtr(60), intPtr(60)},
		nil, intPtr(6), nil, nil, // 6 weeks
		booster, nil,
		doses,
		dob,
		analysisDate,
	)

	if len(results) == 0 {
		t.Fatal("Expected booster recommendation, got 0")
	}

	expectedDesc := "MenQuadfi - " + booster.Description
	if results[0].Description != expectedDesc {
		t.Errorf("Expected description '%s', got '%s'", expectedDesc, results[0].Description)
	}
}

func contains(s, substr string) bool {

	return len(s) >= len(substr) && (s == substr || (len(substr) > 0 && (len(s) > len(substr) && (s[:len(substr)] == substr || contains(s[1:], substr)))))
}


