package engine

import (
	"path/filepath"
	"testing"
	"time"
	"tracuutiemchung-engine/models"
)

func TestIntegration_BasicFlow(t *testing.T) {
	filePath := filepath.Join("..", "vercel-backend", "assets", "vaccine_rules.json")
	rules, err := LoadVaccineRules(filePath)
	if err != nil {
		t.Fatalf("Failed to load rules: %v", err)
	}

	dob := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	analysisDate := time.Date(2024, 4, 15, 0, 0, 0, 0, time.UTC)

	record := models.PatientRecord{
		BirthDate: dob,
		AdministeredMap: map[string][]models.AdministeredDose{
			"Hexaxim": {
				{VaccineName: "Hexaxim", Date: time.Date(2024, 3, 1, 0, 0, 0, 0, time.UTC)},
			},
			"Synflorix": {
				{VaccineName: "synflorix", Date: time.Date(2024, 3, 1, 0, 0, 0, 0, time.UTC)},
			},
		},
	}

	finalResults := ProcessAllRules(rules, record, analysisDate)

	if len(finalResults) == 0 {
		t.Fatal("Expected recommendations, got none")
	}

	// Verify 6-in-1 Dose 2
	found6in1 := false
	for _, item := range finalResults {
		if item.VaccineName == "Vắc xin 6 trong 1 (Hexaxim/Infanrix Hexa)" {
			found6in1 = true
			expectedDate := time.Date(2024, 3, 31, 0, 0, 0, 0, time.UTC) // 3-01 + 30 days
			if item.EarliestNextDoseDate == nil || !item.EarliestNextDoseDate.Equal(expectedDate) {
				t.Errorf("6-in-1 expected date %v, got %v", expectedDate, item.EarliestNextDoseDate)
			}
		}
	}
	if !found6in1 {
		t.Error("6-in-1 recommendation not found")
	}
}

func TestIntegration_LiveVaccineSpacing(t *testing.T) {
	filePath := filepath.Join("..", "vercel-backend", "assets", "vaccine_rules.json")
	rules, err := LoadVaccineRules(filePath)
	if err != nil {
		t.Fatalf("Failed to load rules: %v", err)
	}

	dob := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)
	// Analysis date: child is 15 months old
	analysisDate := time.Date(2024, 4, 1, 0, 0, 0, 0, time.UTC)

	record := models.PatientRecord{
		BirthDate: dob,
		AdministeredMap: map[string][]models.AdministeredDose{
			"MMR-II": {
				{VaccineName: "MMR-II", Date: time.Date(2024, 3, 20, 0, 0, 0, 0, time.UTC)},
			},
		},
	}

	finalResults := ProcessAllRules(rules, record, analysisDate)

	// Varivax should be recommended. 
	// Varivax min age is 12 months (passed).
	// But MMR (live) was given on March 20.
	// Varivax (live) should be pushed to March 20 + 28 days = April 17.

	foundVarivax := false
	for _, item := range finalResults {
		if item.VaccineName == "Varivax (Thủy đậu)" {
			foundVarivax = true
			expectedMinDate := time.Date(2024, 4, 17, 0, 0, 0, 0, time.UTC)
			if item.EarliestNextDoseDate == nil || item.EarliestNextDoseDate.Before(expectedMinDate) {
				t.Errorf("Varivax expected date at least %v, got %v", expectedMinDate, item.EarliestNextDoseDate)
			}
			
			// Check for spacing_adjusted tag
			hasTag := false
			for _, tag := range item.StatusTags {
				if tag == "spacing_adjusted" {
					hasTag = true
					break
				}
			}
			if !hasTag {
				t.Error("Varivax should have 'spacing_adjusted' tag")
			}
		}
	}
	if !foundVarivax {
		t.Error("Varivax recommendation not found")
	}
}

func TestIntegration_PneumoInterchange(t *testing.T) {
	filePath := filepath.Join("..", "vercel-backend", "assets", "vaccine_rules.json")
	rules, err := LoadVaccineRules(filePath)
	if err != nil {
		t.Fatalf("Failed to load rules: %v", err)
	}

	dob := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	analysisDate := time.Date(2024, 5, 1, 0, 0, 0, 0, time.UTC)

	record := models.PatientRecord{
		BirthDate: dob,
		AdministeredMap: map[string][]models.AdministeredDose{
			"Synflorix": {
				{VaccineName: "synflorix", Date: time.Date(2024, 3, 1, 0, 0, 0, 0, time.UTC)},
			},
			"Prevenar13": {
				{VaccineName: "prevenar 13", Date: time.Date(2024, 4, 1, 0, 0, 0, 0, time.UTC)},
			},
		},
	}

	finalResults := ProcessAllRules(rules, record, analysisDate)

	foundPneumo := false
	for _, item := range finalResults {
		if item.VaccineName == "Prevenar 13 (Phế cầu)" || item.VaccineName == "Synflorix (Phế cầu)" {
			foundPneumo = true
			hasTag := false
			for _, tag := range item.StatusTags {
				if tag == "error_interchange" {
					hasTag = true
					break
				}
			}
			if !hasTag {
				t.Errorf("Pneumo should have 'error_interchange' tag, got %v", item.StatusTags)
			}
		}
	}
	if !foundPneumo {
		t.Error("Pneumo recommendation not found")
	}
}
