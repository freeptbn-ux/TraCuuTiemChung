package analyzer

import (
	"strings"
	"testing"
	"vercel-backend/pkg/models"
)

func TestPneumo_NoDoses(t *testing.T) {
	rulesPath := "../../assets/vaccine_rules.json"
	dob, _ := ParseDateDDMMYYYY("01/01/2024")
	analysisDate, _ := ParseDateDDMMYYYY("01/03/2024") // 2 months old
	
	engine, _ := NewEngine(rulesPath, dob, analysisDate)
	results := engine.Analyze([]models.VaccineRecord{})
	
	// Should show all 3 PCVs
	foundP := false
	foundV := false
	foundS := false
	
	for _, res := range results {
		if res.VaccineNameForPopup == "Prevenar 13 (Phế cầu)" {
			foundP = true
		}
		if res.VaccineNameForPopup == "Vaxneuvance (Phế cầu)" {
			foundV = true
		}
		if res.VaccineNameForPopup == "Synflorix (Phế cầu)" {
			foundS = true
		}
	}
	
	if !foundP || !foundV || !foundS {
		t.Errorf("Expected all 3 PCVs to be recommended for no doses, got P:%v V:%v S:%v", foundP, foundV, foundS)
	}
}

func TestPneumo_Mixed_Prevenar_Synflorix(t *testing.T) {
	rulesPath := "../../assets/vaccine_rules.json"
	dob, _ := ParseDateDDMMYYYY("01/01/2024")
	analysisDate, _ := ParseDateDDMMYYYY("01/05/2024")
	
	engine, _ := NewEngine(rulesPath, dob, analysisDate)
	history := []models.VaccineRecord{
		{VaccineName: "Prevenar 13", Date: dob.AddDate(0, 2, 0)},
		{VaccineName: "Synflorix", Date: dob.AddDate(0, 3, 0)},
	}
	
	results := engine.Analyze(history)
	
	foundWarning := false
	for _, res := range results {
		if res.VaccineNameForPopup == "Phế cầu (nhiều loại)" {
			foundWarning = true
			expected := "Cảnh báo: Đã ghi nhận tiêm xen kẽ các loại phế cầu (Prevenar 13 (Phế cầu) và Synflorix (Phế cầu))"
			if res.Description[:len(expected)] != expected {
				t.Errorf("Expected warning description to start with %q, got %q", expected, res.Description)
			}
		}
	}
	
	if !foundWarning {
		t.Errorf("Expected mixed series warning")
	}
}

func TestPneumo_Pneumovax23_Done(t *testing.T) {
	rulesPath := "../../assets/vaccine_rules.json"
	dob, _ := ParseDateDDMMYYYY("01/01/2020")
	analysisDate, _ := ParseDateDDMMYYYY("01/01/2024")
	
	engine, _ := NewEngine(rulesPath, dob, analysisDate)
	history := []models.VaccineRecord{
		{VaccineName: "Pneumovax 23", Date: dob.AddDate(2, 0, 0)},
	}
	
	results := engine.Analyze(history)
	
	// Should skip all PCVs
	for _, res := range results {
		name := res.VaccineNameForPopup
		if name == "Prevenar 13 (Phế cầu)" || name == "Vaxneuvance (Phế cầu)" || name == "Synflorix (Phế cầu)" {
			t.Errorf("PCV %s should be skipped when Pneumovax 23 is done", name)
		}
	}
}

func TestPneumo_Prevenar_2Doses_Over2Years(t *testing.T) {
	rulesPath := "../../assets/vaccine_rules.json"
	dob, _ := ParseDateDDMMYYYY("01/01/2022")
	analysisDate, _ := ParseDateDDMMYYYY("01/02/2024") // 25 months old
	
	engine, _ := NewEngine(rulesPath, dob, analysisDate)
	history := []models.VaccineRecord{
		{VaccineName: "Prevenar 13", Date: dob.AddDate(0, 2, 0)},
		{VaccineName: "Prevenar 13", Date: dob.AddDate(0, 4, 0)},
	}
	
	results := engine.Analyze(history)
	
	foundPneumovaxAlternative := false
	for _, res := range results {
		if strings.Contains(res.VaccineNameForPopup, "Pneumovax 23") && containsStr(res.StatusTags, "alternative_completion") {
			foundPneumovaxAlternative = true
		}
	}
	
	if !foundPneumovaxAlternative {
		t.Errorf("Expected Pneumovax 23 as alternative completion for child > 2y with < 3 doses PCV")
	}
}

func TestPneumo_Prevenar_3Doses_Over2Years(t *testing.T) {
	rulesPath := "../../assets/vaccine_rules.json"
	dob, _ := ParseDateDDMMYYYY("01/01/2022")
	analysisDate, _ := ParseDateDDMMYYYY("01/02/2024") // 25 months old
	
	engine, _ := NewEngine(rulesPath, dob, analysisDate)
	history := []models.VaccineRecord{
		{VaccineName: "Prevenar 13", Date: dob.AddDate(0, 2, 0)},
		{VaccineName: "Prevenar 13", Date: dob.AddDate(0, 4, 0)},
		{VaccineName: "Prevenar 13", Date: dob.AddDate(0, 6, 0)},
	}
	
	results := engine.Analyze(history)
	
	foundPneumovaxBooster := false
	foundPCVDose4 := false
	for _, res := range results {
		if strings.Contains(res.VaccineNameForPopup, "Pneumovax 23") && containsStr(res.StatusTags, "alternative_booster") {
			foundPneumovaxBooster = true
		}
		if strings.Contains(res.VaccineNameForPopup, "Prevenar 13") && res.DoseNumber == 4 {
			foundPCVDose4 = true
		}
	}
	
	if !foundPneumovaxBooster {
		t.Errorf("Expected Pneumovax 23 as alternative booster for child > 2y with 3 doses PCV")
	}
	if foundPCVDose4 {
		t.Errorf("Did not expect PCV Dose 4 to be recommended (Pneumovax 23 should replace it)")
	}
}

func TestPneumo_Prevenar_4Doses(t *testing.T) {
	rulesPath := "../../assets/vaccine_rules.json"
	dob, _ := ParseDateDDMMYYYY("01/01/2022")
	analysisDate, _ := ParseDateDDMMYYYY("01/03/2024")
	
	engine, _ := NewEngine(rulesPath, dob, analysisDate)
	history := []models.VaccineRecord{
		{VaccineName: "Prevenar 13", Date: dob.AddDate(0, 2, 0)},
		{VaccineName: "Prevenar 13", Date: dob.AddDate(0, 4, 0)},
		{VaccineName: "Prevenar 13", Date: dob.AddDate(0, 6, 0)},
		{VaccineName: "Prevenar 13", Date: dob.AddDate(1, 0, 0)},
	}
	
	results := engine.Analyze(history)
	
	for _, res := range results {
		if res.VaccineNameForPopup == "Prevenar 13 (Phế cầu)" {
			t.Errorf("Prevenar 13 should be completed and not recommended")
		}
	}
}

func TestCumulative_NoDoses(t *testing.T) {
	// Mock a cumulative rule
	rules := map[string]VaccineRule{
		"DTP_Group": {
			DisplayName: "Bạch hầu - Ho gà - Uốn ván",
			Type:        RuleTypeGroupCumulativeUnique,
			DosesRequired: 3,
			MinIntervalDays: []*int{nil, ptr(30), ptr(30)},
			NamesNormGroup: []string{"hexaxim", "infanrix hexa", "pentaxim", "quinvaxem"},
		},
	}
	
	dob, _ := ParseDateDDMMYYYY("01/01/2024")
	analysisDate, _ := ParseDateDDMMYYYY("01/03/2024")
	
	engine := &Engine{Rules: rules, DOB: dob, AnalysisDate: analysisDate}
	results := engine.Analyze([]models.VaccineRecord{})
	
	found := false
	for _, res := range results {
		if res.VaccineNameForPopup == "Bạch hầu - Ho gà - Uốn ván" {
			found = true
			if res.DoseNumber != 1 {
				t.Errorf("Expected dose 1, got %d", res.DoseNumber)
			}
		}
	}
	if !found {
		t.Errorf("Cumulative rule not found in results")
	}
}

func ptr(i int) *int { return &i }
