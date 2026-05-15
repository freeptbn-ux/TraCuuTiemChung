package analyzer

import (
	"testing"
	"time"
	"vercel-backend/pkg/models"
)

func TestPostProc_GeneralSpacing(t *testing.T) {
	rulesPath := "../../assets/vaccine_rules.json"
	dob, _ := ParseDateDDMMYYYY("01/01/2024")
	analysisDate, _ := ParseDateDDMMYYYY("01/07/2024")
	
	engine, _ := NewEngine(rulesPath, dob, analysisDate)
	
	// Last vaccine was 5 days ago
	lastDate := analysisDate.AddDate(0, 0, -5)
	history := []models.VaccineRecord{
		{VaccineName: "Hexaxim", Date: lastDate},
	}
	
	// Create a mock result that would be due today
	results := []AnalysisResult{
		{
			VaccineNameForPopup:  "Varivax (Thủy đậu)",
			Description:          "Varivax (Thủy đậu)",
			EarliestNextDoseDate: &analysisDate,
		},
	}
	
	administeredMap := make(map[string][]models.VaccineRecord)
	administeredMap["hexaxim"] = history
	
	processed := engine.ApplySpacingAndSort(results, administeredMap)
	
	// Should be adjusted to lastDate + 14 days
	expectedDate := lastDate.AddDate(0, 0, 14)
	if !processed[0].EarliestNextDoseDate.Equal(expectedDate) {
		t.Errorf("Expected spacing to be adjusted to %v, got %v", expectedDate, *processed[0].EarliestNextDoseDate)
	}
}

func TestPostProc_LiveVaccineSpacing(t *testing.T) {
	rulesPath := "../../assets/vaccine_rules.json"
	dob, _ := ParseDateDDMMYYYY("01/01/2023")
	analysisDate, _ := ParseDateDDMMYYYY("01/02/2024")
	
	engine, _ := NewEngine(rulesPath, dob, analysisDate)
	
	// Last vaccine was a live vaccine (Varivax) 10 days ago
	lastLiveDate := analysisDate.AddDate(0, 0, -10)
	history := []models.VaccineRecord{
		{VaccineName: "Varivax", Date: lastLiveDate},
	}
	
	// Result is another live vaccine (Priorix)
	results := []AnalysisResult{
		{
			VaccineNameForPopup:  "Vắc xin Sởi-Quai bị-Rubella (MMR-II/Priorix)",
			Description:          "Priorix (Sởi - Quai bị - Rubella)",
			EarliestNextDoseDate: &analysisDate,
		},
	}
	
	administeredMap := make(map[string][]models.VaccineRecord)
	administeredMap["varivax"] = history
	
	processed := engine.ApplySpacingAndSort(results, administeredMap)
	
	// Should be adjusted to lastLiveDate + 28 days
	expectedDate := lastLiveDate.AddDate(0, 0, 28)
	if !processed[0].EarliestNextDoseDate.Equal(expectedDate) {
		t.Errorf("Expected live-to-live spacing to be 28 days: expected %v, got %v", expectedDate, *processed[0].EarliestNextDoseDate)
	}
}

func TestPostProc_Sort(t *testing.T) {
	dob, _ := ParseDateDDMMYYYY("01/01/2024")
	analysisDate, _ := ParseDateDDMMYYYY("01/07/2024")
	engine := &Engine{AnalysisDate: analysisDate, DOB: dob}
	
	d1 := analysisDate.AddDate(0, 0, 10)
	d2 := analysisDate.AddDate(0, 0, 5)
	
	results := []AnalysisResult{
		{Description: "B", EarliestNextDoseDate: &d1},
		{Description: "A", EarliestNextDoseDate: &d2},
		{Description: "C", EarliestNextDoseDate: nil},
	}
	
	processed := engine.ApplySpacingAndSort(results, nil)
	
	if processed[0].Description != "A" {
		t.Errorf("Sort failed: expected A first, got %s", processed[0].Description)
	}
	if processed[1].Description != "B" {
		t.Errorf("Sort failed: expected B second, got %s", processed[1].Description)
	}
	if processed[2].Description != "C" {
		t.Errorf("Sort failed: expected C last, got %s", processed[2].Description)
	}
}

func TestIsLive(t *testing.T) {
	rulesPath := "../../assets/vaccine_rules.json"
	dob, _ := ParseDateDDMMYYYY("01/01/2024")
	engine, _ := NewEngine(rulesPath, dob, time.Now())
	
	// Imojev is live
	if !engine.isMissingItemLive(AnalysisResult{Description: "Imojev (Sanofi Pasteur)"}) {
		t.Error("Imojev should be recognized as live")
	}
	
	// Jevax is inactivated
	if engine.isMissingItemLive(AnalysisResult{Description: "Jevax/VNNB (Việt Nam)"}) {
		t.Error("Jevax should be recognized as inactivated")
	}
	
	// MMR is live
	if !engine.isMissingItemLive(AnalysisResult{VaccineNameForPopup: "Vắc xin Sởi-Quai bị-Rubella (MMR-II/Priorix)"}) {
		t.Error("MMR group should be recognized as live")
	}
}
