package analyzer

import (
	"testing"

	"vercel-backend/pkg/models"
)

func TestEngine_PneumoRules(t *testing.T) {
	rulesPath := "../../assets/vaccine_rules.json"
	
	t.Run("Synflorix active series - Only Synflorix recommended", func(t *testing.T) {
		dob, _ := ParseDateDDMMYYYY("01/01/2024")
		analysisDate, _ := ParseDateDDMMYYYY("01/05/2024")
		
		engine, _ := NewEngine(rulesPath, dob, analysisDate)
		
		history := []models.VaccineRecord{
			{VaccineName: "Synflorix", Date: dob.AddDate(0, 2, 0)}, // 2 months old
		}
		
		results := engine.Analyze(history)
		
		foundSynflorix := false
		foundPrevenar := false
		foundVaxneuvance := false
		
		for _, res := range results {
			if res.VaccineNameForPopup == "Synflorix (Phế cầu)" {
				foundSynflorix = true
			}
			if res.VaccineNameForPopup == "Prevenar 13 (Phế cầu)" {
				foundPrevenar = true
			}
			if res.VaccineNameForPopup == "Vaxneuvance (Phế cầu)" {
				foundVaxneuvance = true
			}
		}
		
		if !foundSynflorix {
			t.Errorf("Expected Synflorix recommendation")
		}
		if foundPrevenar || foundVaxneuvance {
			t.Errorf("Expected NO Prevenar or Vaxneuvance recommendation when Synflorix is active")
		}
	})

	t.Run("Mixed series - Warning displayed", func(t *testing.T) {
		dob, _ := ParseDateDDMMYYYY("01/01/2024")
		analysisDate, _ := ParseDateDDMMYYYY("01/07/2024")
		
		engine, _ := NewEngine(rulesPath, dob, analysisDate)
		
		history := []models.VaccineRecord{
			{VaccineName: "Synflorix", Date: dob.AddDate(0, 2, 0)},
			{VaccineName: "Prevenar 13", Date: dob.AddDate(0, 3, 0)},
		}
		
		results := engine.Analyze(history)
		
		foundMixedWarning := false
		for _, res := range results {
			if res.VaccineNameForPopup == "Phế cầu (nhiều loại)" {
				foundMixedWarning = true
				if !contains(res.StatusTags, "pneumo_mixed") {
					t.Errorf("Expected pneumo_mixed tag for mixed series")
				}
			}
		}
		
		if !foundMixedWarning {
			t.Errorf("Expected mixed series warning")
		}
	})

	t.Run("Over 2 years - Suggest Pneumovax 23", func(t *testing.T) {
		dob, _ := ParseDateDDMMYYYY("01/01/2022")
		analysisDate, _ := ParseDateDDMMYYYY("01/02/2024") // 25 months old
		
		engine, _ := NewEngine(rulesPath, dob, analysisDate)
		
		history := []models.VaccineRecord{} // No vaccines
		
		results := engine.Analyze(history)
		
		foundPneumovax := false
		for _, res := range results {
			if res.VaccineNameForPopup == "Pneumovax 23 / PNEUMO 23 (Phế cầu)" {
				foundPneumovax = true
				if !contains(res.StatusTags, "eligible") {
					t.Errorf("Expected Pneumovax 23 to be eligible")
				}
			}
		}
		
		if !foundPneumovax {
			t.Errorf("Expected Pneumovax 23 recommendation for patient > 2 years")
		}
	})
}

func contains(slice []string, val string) bool {
	for _, item := range slice {
		if item == val {
			return true
		}
	}
	return false
}
