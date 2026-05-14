package analyzer

import (
	"strings"
	"testing"
	"vercel-backend/pkg/models"
)

func TestGroupSpecialCheckers(t *testing.T) {
	rulesPath := "../../assets/vaccine_rules.json"

	t.Run("Test MMR Interval with MVVAC", func(t *testing.T) {
		dob, _ := ParseDateDDMMYYYY("01/01/2023")
		// MVVAC lúc 9 tháng tuổi
		mvvacDate, _ := ParseDateDDMMYYYY("01/10/2023")
		
		// Ngày phân tích: 10 ngày sau MVVAC
		analysisDate := mvvacDate.AddDate(0, 0, 10)
		engine, _ := NewEngine(rulesPath, dob, analysisDate)

		history := []models.VaccineRecord{
			{VaccineName: "MVVAC", Date: mvvacDate},
		}

		results := engine.Analyze(history)
		found := false
		for _, res := range results {
			if res.VaccineNameForPopup == "Vắc xin Sởi-Quai bị-Rubella (MMR-II/Priorix)" {
				found = true
				// Phải cách MVVAC ít nhất 84 ngày
				expectedEarliest := mvvacDate.AddDate(0, 0, 84)
				if !res.EarliestNextDoseDate.Equal(expectedEarliest) {
					t.Errorf("Expected earliest MMR date %v, got %v", expectedEarliest, res.EarliestNextDoseDate)
				}
				if !sliceContains(res.StatusTags, "too_young") {
					t.Errorf("Expected 'too_young' tag, got %v", res.StatusTags)
				}
			}
		}
		if !found {
			t.Error("MMR recommendation not found")
		}
	})

	t.Run("Test Flu Annual Booster", func(t *testing.T) {
		dob, _ := ParseDateDDMMYYYY("01/01/2010")
		// Tiêm cúm 2 năm trước
		lastFluDate, _ := ParseDateDDMMYYYY("01/01/2022")
		analysisDate, _ := ParseDateDDMMYYYY("01/01/2024")
		engine, _ := NewEngine(rulesPath, dob, analysisDate)

		history := []models.VaccineRecord{
			{VaccineName: "Vaxigrip Tetra", Date: lastFluDate},
		}

		results := engine.Analyze(history)
		found := false
		for _, res := range results {
			if res.VaccineNameForPopup == "Vắc xin Cúm" {
				found = true
				if !sliceContains(res.StatusTags, "due") {
					t.Errorf("Expected 'due' tag for annual Flu booster, got %v", res.StatusTags)
				}
				if !strings.Contains(res.Description, "nhắc lại hàng năm") {
					t.Errorf("Expected 'nhắc lại hàng năm' in description, got %q", res.Description)
				}
			}
		}
		if !found {
			t.Error("Flu recommendation not found")
		}
	})

	t.Run("Test Meningococcal ACYW Interaction with BC", func(t *testing.T) {
		dob, _ := ParseDateDDMMYYYY("01/01/2020")
		// Tiêm VA-MENGOC-BC gần đây (cách 10 ngày)
		bcDate, _ := ParseDateDDMMYYYY("01/05/2024")
		analysisDate, _ := ParseDateDDMMYYYY("11/05/2024")
		engine, _ := NewEngine(rulesPath, dob, analysisDate)

		history := []models.VaccineRecord{
			{VaccineName: "VA - MENGOC - BC", Date: bcDate},
		}

		results := engine.Analyze(history)
		found := false
		for _, res := range results {
			if res.VaccineNameForPopup == "Vắc xin Não mô cầu ACYW-135 (Menactra/MenQuadfi)" {
				found = true
				// Phải có cảnh báo tương tác
				if !strings.Contains(res.Description, "VA-Mengoc BC") {
					t.Errorf("Expected interaction warning in description, got %q", res.Description)
				}
				// Phải cách BC ít nhất 60 ngày
				expectedEarliest := bcDate.AddDate(0, 0, 60)
				if !res.EarliestNextDoseDate.Equal(expectedEarliest) {
					t.Errorf("Expected earliest ACYW date %v, got %v", expectedEarliest, res.EarliestNextDoseDate)
				}
			}
		}
		if !found {
			t.Error("ACYW recommendation not found")
		}
	})
}

func sliceContains(slice []string, val string) bool {
	for _, item := range slice {
		if item == val {
			return true
		}
	}
	return false
}
