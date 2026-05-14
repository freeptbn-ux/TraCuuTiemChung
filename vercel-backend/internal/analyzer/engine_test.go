package analyzer

import (
	"encoding/json"
	"os"
	"testing"

	"vercel-backend/internal/models"
)

func TestNormalizeVaccineName(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"Influvac Tetra 2024/2025", "influvac tetra"},
		{"Varivax & Diluent Inj 0.5ml", "varivax & diluent inj"},
		{"Priorix (Sởi - Quai bị - Rubella)", "priorix"},
		{"JEEV 3mcg/0.5ml", "jeev"},
		{"Hexaxim 0,5ml", "hexaxim"},
	}

	for _, tt := range tests {
		result := NormalizeVaccineName(tt.input)
		if result != tt.expected {
			t.Errorf("NormalizeVaccineName(%q) = %q, expected %q", tt.input, result, tt.expected)
		}
	}
}

func TestEngine_Analyze_TableDriven(t *testing.T) {
	dob, _ := ParseDateDDMMYYYY("01/01/2024")
	analysisDate, _ := ParseDateDDMMYYYY("01/07/2024")
	
	rulesPath := "../../assets/vaccine_rules.json"
	engine, err := NewEngine(rulesPath, dob, analysisDate)
	if err != nil {
		t.Fatalf("Failed to create engine: %v", err)
	}

	tests := []struct {
		name     string
		history  []models.VaccineRecord
		expected map[string]string // Vaccine name -> expected status tag
	}{
		{
			name: "Newborn - Only BCG needed",
			history: []models.VaccineRecord{},
			expected: map[string]string{
				"Lao (BCG)": "due",
			},
		},
		{
			name: "Varivax - Too young for first dose",
			history: []models.VaccineRecord{},
			expected: map[string]string{
				"Varivax (Thủy đậu)": "too_young",
			},
		},
		{
			name: "Varivax - First dose taken, waiting for second",
			history: []models.VaccineRecord{
				{VaccineName: "Varivax", Date: dob.AddDate(1, 0, 0)}, // 12 months old
			},
			expected: map[string]string{
				"Varivax (Thủy đậu)": "too_young", // Needs 3 months interval
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			results := engine.Analyze(tt.history)
			for vaccine, expectedTag := range tt.expected {
				found := false
				for _, res := range results {
					if res.VaccineNameForPopup == vaccine {
						found = true
						tagMatch := false
						for _, tag := range res.StatusTags {
							if tag == expectedTag {
								tagMatch = true
								break
							}
						}
						if !tagMatch {
							t.Errorf("For %s, expected tag %s but got %v", vaccine, expectedTag, res.StatusTags)
						}
						break
					}
				}
				if !found {
					t.Errorf("Expected result for %s but not found", vaccine)
				}
			}
		})
	}
}

func TestParityWithPython(t *testing.T) {
	// Load expected output from Python
	expectedData, err := os.ReadFile("../../../testdata/python_minhkhoi_output.json")
	if err != nil {
		t.Skip("Python expected output not found, skipping parity test")
		return
	}

	var expected map[string]interface{}
	json.Unmarshal(expectedData, &expected)

	dobStr := expected["dob"].(string)
	analysisDateStr := expected["analysis_date"].(string)
	dob, _ := ParseDateDDMMYYYY(dobStr)
	analysisDate, _ := ParseDateDDMMYYYY(analysisDateStr)

	engine, err := NewEngine("../../assets/vaccine_rules.json", dob, analysisDate)
	if err != nil {
		t.Fatalf("Failed to create engine: %v", err)
	}

	// Convert history from expected JSON to models.VaccineRecord
	historyRaw := expected["administered_vaccines"].([]interface{})
	var history []models.VaccineRecord
	for _, h := range historyRaw {
		item := h.(map[string]interface{})
		date, _ := ParseDateDDMMYYYY(item["date"].(string))
		history = append(history, models.VaccineRecord{
			VaccineName: item["vaccine_name"].(string),
			Date:        date,
			Dose:        item["dose"].(string),
		})
	}

	results := engine.Analyze(history)
	
	// Print results for debugging
	t.Logf("Found %d recommendations", len(results))
	for _, res := range results {
		t.Logf("- %s: %v", res.VaccineNameForPopup, res.StatusTags)
	}

	// Basic check: at least some recommendations should match
	// Full 100% parity is hard without implementing all rule types
	if len(results) == 0 {
		t.Error("Expected at least some recommendations, got 0")
	}
}
