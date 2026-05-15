package analyzer

import (
	"strings"
	"testing"
	"vercel-backend/pkg/models"
)

func TestAltMinAge_Rota(t *testing.T) {
	rulesPath := "../../assets/vaccine_rules.json"

	t.Run("Rota_NoDoses_Eligible", func(t *testing.T) {
		dob, _ := ParseDateDDMMYYYY("01/01/2024")
		analysisDate, _ := ParseDateDDMMYYYY("01/03/2024") // 2 months
		engine, _ := NewEngine(rulesPath, dob, analysisDate)

		results := engine.Analyze([]models.VaccineRecord{})
		found := false
		for _, res := range results {
			if strings.Contains(res.VaccineNameForPopup, "Rota") {
				found = true
				// Expected: "Vắc xin Rota (Lựa chọn: Rota Teq (Mỹ) Hoặc Rotarix/ROTARIXTM (Bỉ) Hoặc Rotavin/Rotavin-M1 (Việt Nam)). đủ điều kiện tuổi"
				if !strings.Contains(res.Description, "Rotarix") && !strings.Contains(res.Description, "Rota Teq") {
					t.Errorf("Expected options in description, got %q", res.Description)
				}
				if !containsStr(res.StatusTags, "due") {
					t.Errorf("Expected 'due' tag, got %v", res.StatusTags)
				}
			}
		}
		if !found {
			t.Error("Rota result not found")
		}
	})

	t.Run("Rota_NoDoses_TooOld", func(t *testing.T) {
		dob, _ := ParseDateDDMMYYYY("01/01/2024")
		analysisDate, _ := ParseDateDDMMYYYY("01/09/2024") // 8 months
		engine, _ := NewEngine(rulesPath, dob, analysisDate)

		results := engine.Analyze([]models.VaccineRecord{})
		found := false
		for _, res := range results {
			if strings.Contains(res.VaccineNameForPopup, "Rota") {
				found = true
				if !containsStr(res.StatusTags, "too_old_to_start") {
					t.Errorf("Expected 'too_old_to_start' tag, got %v", res.StatusTags)
				}
			}
		}
		if !found {
			t.Error("Rota result not found")
		}
	})

	t.Run("Rota_1Dose_Rotarix", func(t *testing.T) {
		dob, _ := ParseDateDDMMYYYY("01/01/2024")
		analysisDate, _ := ParseDateDDMMYYYY("15/02/2024")
		engine, _ := NewEngine(rulesPath, dob, analysisDate)

		history := []models.VaccineRecord{
			{VaccineName: "Rotarix", Date: dob.AddDate(0, 0, 45)}, // 6 weeks
		}

		results := engine.Analyze(history)
		found := false
		for _, res := range results {
			if strings.Contains(res.VaccineNameForPopup, "Rotarix") {
				found = true
				if !strings.Contains(strings.ToLower(res.Description), "mũi 2") {
					t.Errorf("Expected 'mũi 2' in description, got %q", res.Description)
				}
			}
		}
		if !found {
			t.Error("Rotarix result not found")
		}
	})

	t.Run("Rota_TooOldToComplete", func(t *testing.T) {
		dob, _ := ParseDateDDMMYYYY("01/01/2024")
		analysisDate, _ := ParseDateDDMMYYYY("01/10/2024") // 9 months
		engine, _ := NewEngine(rulesPath, dob, analysisDate)

		history := []models.VaccineRecord{
			{VaccineName: "Rotarix", Date: dob.AddDate(0, 2, 0)},
		}

		results := engine.Analyze(history)
		found := false
		for _, res := range results {
			if strings.Contains(res.VaccineNameForPopup, "Rotarix") {
				found = true
				if !containsStr(res.StatusTags, "too_old_to_complete") {
					t.Errorf("Expected 'too_old_to_complete' tag, got %v", res.StatusTags)
				}
			}
		}
		if !found {
			t.Error("Rota result not found")
		}
	})
}

func TestAltAgeRange_JE(t *testing.T) {
	rulesPath := "../../assets/vaccine_rules.json"

	t.Run("JE_NoDoses", func(t *testing.T) {
		dob, _ := ParseDateDDMMYYYY("01/01/2023")
		analysisDate, _ := ParseDateDDMMYYYY("01/02/2024") // 13 months
		engine, _ := NewEngine(rulesPath, dob, analysisDate)

		results := engine.Analyze([]models.VaccineRecord{})
		found := false
		for _, res := range results {
			if strings.Contains(res.VaccineNameForPopup, "Viêm não Nhật Bản") {
				found = true
				if !strings.Contains(res.Description, "Jevax") && !strings.Contains(res.Description, "Imojev") {
					t.Errorf("Expected Jevax and Imojev options, got %q", res.Description)
				}
			}
		}
		if !found {
			t.Error("JE result not found")
		}
	})

	t.Run("JE_Jevax3_Complete", func(t *testing.T) {
		dob, _ := ParseDateDDMMYYYY("01/01/2020")
		analysisDate, _ := ParseDateDDMMYYYY("01/01/2024")
		engine, _ := NewEngine(rulesPath, dob, analysisDate)

		history := []models.VaccineRecord{
			{VaccineName: "Jevax", Date: dob.AddDate(1, 0, 0)},
			{VaccineName: "Jevax", Date: dob.AddDate(1, 0, 7)},
			{VaccineName: "Jevax", Date: dob.AddDate(2, 0, 0)},
		}

		results := engine.Analyze(history)
		found := false
		for _, res := range results {
			if strings.Contains(res.VaccineNameForPopup, "Jevax") {
				found = true
				if !strings.Contains(res.Description, "nhắc lại") && !strings.Contains(res.Description, "Imojev") {
					t.Errorf("Expected booster or Imojev recommendation, got %q", res.Description)
				}
			}
		}
		if !found {
			t.Error("JE result not found")
		}
	})

	t.Run("JE_Jevax3_Imojev1", func(t *testing.T) {
		dob, _ := ParseDateDDMMYYYY("01/01/2020")
		analysisDate, _ := ParseDateDDMMYYYY("01/01/2024")
		engine, _ := NewEngine(rulesPath, dob, analysisDate)

		history := []models.VaccineRecord{
			{VaccineName: "Jevax", Date: dob.AddDate(1, 0, 0)},
			{VaccineName: "Jevax", Date: dob.AddDate(1, 0, 7)},
			{VaccineName: "Jevax", Date: dob.AddDate(2, 0, 0)},
			{VaccineName: "Imojev", Date: dob.AddDate(3, 0, 0)},
		}

		results := engine.Analyze(history)
		// Parity check: Should have a warning about mixed series, even if completed
		found := false
		for _, res := range results {
			if containsStr(res.StatusTags, "je_mixed_warning") {
				found = true
				break
			}
		}
		if !found {
			t.Error("Expected je_mixed_warning tag for mixed series")
		}
	})

	t.Run("JE_Imojev_Then_Jevax", func(t *testing.T) {
		dob, _ := ParseDateDDMMYYYY("01/01/2020")
		analysisDate, _ := ParseDateDDMMYYYY("01/01/2024")
		engine, _ := NewEngine(rulesPath, dob, analysisDate)

		history := []models.VaccineRecord{
			{VaccineName: "Imojev", Date: dob.AddDate(1, 0, 0)},
			{VaccineName: "Jevax", Date: dob.AddDate(2, 0, 0)},
		}

		results := engine.Analyze(history)
		found := false
		for _, res := range results {
			if containsStr(res.StatusTags, "error_interchange") {
				found = true
			}
		}
		if !found {
			t.Error("Expected error_interchange tag for Imojev then Jevax")
		}
	})
}

func TestAltAgeRange_HepA(t *testing.T) {
	rulesPath := "../../assets/vaccine_rules.json"

	t.Run("HepA_NoDoses_12m", func(t *testing.T) {
		dob, _ := ParseDateDDMMYYYY("01/01/2023")
		analysisDate, _ := ParseDateDDMMYYYY("01/01/2024") // 12 months
		engine, _ := NewEngine(rulesPath, dob, analysisDate)

		results := engine.Analyze([]models.VaccineRecord{})
		found := false
		for _, res := range results {
			if strings.Contains(res.VaccineNameForPopup, "Viêm Gan A") {
				found = true
				if !strings.Contains(res.Description, "Avaxim") {
					t.Errorf("Expected Avaxim option, got %q", res.Description)
				}
				if !containsStr(res.StatusTags, "due") {
					t.Errorf("Expected 'due' tag, got %v", res.StatusTags)
				}
			}
		}
		if !found {
			t.Error("HepA result not found")
		}
	})

	t.Run("HepA_NoDoses_6m", func(t *testing.T) {
		dob, _ := ParseDateDDMMYYYY("01/01/2024")
		analysisDate, _ := ParseDateDDMMYYYY("01/07/2024") // 6 months
		engine, _ := NewEngine(rulesPath, dob, analysisDate)

		results := engine.Analyze([]models.VaccineRecord{})
		found := false
		for _, res := range results {
			if strings.Contains(res.VaccineNameForPopup, "Viêm Gan A") {
				found = true
				if !containsStr(res.StatusTags, "too_young") {
					t.Errorf("Expected 'too_young' tag, got %v", res.StatusTags)
				}
			}
		}
		if !found {
			t.Error("HepA result not found")
		}
	})
}
