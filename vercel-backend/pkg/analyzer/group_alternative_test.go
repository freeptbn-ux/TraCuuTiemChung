package analyzer

import (
	"testing"
	"vercel-backend/pkg/models"
)

func TestGroupAlternativeCheckers(t *testing.T) {
	rulesPath := "../../assets/vaccine_rules.json"

	t.Run("Test Rota Multiple Courses", func(t *testing.T) {
		dob, _ := ParseDateDDMMYYYY("01/01/2024")
		analysisDate, _ := ParseDateDDMMYYYY("01/03/2024") // 2 months old
		engine, _ := NewEngine(rulesPath, dob, analysisDate)

		results := engine.Analyze([]models.VaccineRecord{})
		found := false
		for _, res := range results {
			if res.VaccineNameForPopup == "Vắc xin Rota" {
				found = true
				// Check if contains course names
				expectedDesc := "Vắc xin Rota (Bắt đầu tiêm: [Rota Teq (Mỹ) Rotarix/ROTARIXTM (Bỉ) Rotavin/Rotavin-M1 (Việt Nam)])"
				if res.Description != expectedDesc {
					t.Errorf("Expected description %q, got %q", expectedDesc, res.Description)
				}
				// Check status
				hasDue := false
				for _, tag := range res.StatusTags {
					if tag == "due" {
						hasDue = true
					}
				}
				if !hasDue {
					t.Errorf("Expected 'due' tag for Rota")
				}
			}
		}
		if !found {
			t.Error("Rota recommendation not found for 2-month-old")
		}
	})

	t.Run("Test Rota Too Old", func(t *testing.T) {
		dob, _ := ParseDateDDMMYYYY("01/01/2024")
		analysisDate, _ := ParseDateDDMMYYYY("01/09/2024") // 8 months old
		engine, _ := NewEngine(rulesPath, dob, analysisDate)

		// MaxAgeMonthsToStartFirstDoseGroup is 6. At 8 months, should be too old.
		results := engine.Analyze([]models.VaccineRecord{})
		found := false
		for _, res := range results {
			if res.VaccineNameForPopup == "Vắc xin Rota" {
				found = true
				hasTooOld := false
				for _, tag := range res.StatusTags {
					if tag == "too_old" {
						hasTooOld = true
					}
				}
				if !hasTooOld {
					t.Errorf("Expected 'too_old' tag for Rota at 8 months, got %v", res.StatusTags)
				}
			}
		}
		if !found {
			t.Error("Rota result not found for 8-month-old")
		}
	})

	t.Run("Test JE Switch Warning", func(t *testing.T) {
		dob, _ := ParseDateDDMMYYYY("01/01/2020")
		analysisDate, _ := ParseDateDDMMYYYY("01/01/2024")
		engine, _ := NewEngine(rulesPath, dob, analysisDate)

		history := []models.VaccineRecord{
			{VaccineName: "Jevax", Date: dob.AddDate(1, 0, 0)},
			{VaccineName: "Imojev", Date: dob.AddDate(2, 0, 0)},
		}

		results := engine.Analyze(history)
		found := false
		for _, res := range results {
			// Vì Jevax tiêm trước nên phác đồ bắt đầu là Jevax
			if res.VaccineNameForPopup == "Jevax/VNNB (Việt Nam)" {
				found = true
				// Check for mixed warning
				if !stringContains(res.Description, "Cảnh báo: Đã tiêm trộn") {
					t.Errorf("Expected mixed warning in description, got %q", res.Description)
				}
			}
		}
		if !found {
			t.Error("JE recommendation not found for mixed courses")
		}
	})
}

func stringContains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || s[0:len(substr)] == substr || s[len(s)-len(substr):] == substr || find(s, substr))
}

func find(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
