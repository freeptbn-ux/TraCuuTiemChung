package analyzer

import (
	"encoding/json"
	"os"
	"testing"
	"time"
)

func TestNormalizeVaccineName_Parentheses(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"Varivax (Thủy đậu)", "varivax"},
		{"Priorix (Sởi - Quai bị - Rubella)", "priorix"},
		{"Heberbiovac HB 0.5ml (Viêm gan B)", "heberbiovac hb"},
		{"6 trong 1 (Infarix)", "6 trong 1"},
	}

	for _, tt := range tests {
		got := NormalizeVaccineName(tt.input)
		if got != tt.expected {
			t.Errorf("NormalizeVaccineName(%q) = %q, expected %q", tt.input, got, tt.expected)
		}
	}
}

func TestGetAgeStatusAndEarliestDate(t *testing.T) {
	dob, _ := ParseDateDDMMYYYY("01/01/2024")
	analysisDate, _ := ParseDateDDMMYYYY("01/03/2024") // 2 months old
	
	engine := &Engine{
		DOB:          dob,
		AnalysisDate: analysisDate,
	}

	t.Run("TooYoung_Months", func(t *testing.T) {
		rule := VaccineRule{MinAgeMonthsAtFirstDose: 6}
		msg, earliest, tags := engine.getAgeStatusAndEarliestDate(rule, "Test")
		if !containsStr(tags, "too_young") {
			t.Errorf("Expected too_young tag, got %v", tags)
		}
		if msg != "cần 6 tháng tuổi" {
			t.Errorf("Expected msg 'cần 6 tháng tuổi', got %q", msg)
		}
		expectedEarliest := AddMonths(dob, 6)
		if !earliest.Equal(expectedEarliest) {
			t.Errorf("Expected earliest %v, got %v", expectedEarliest, *earliest)
		}
	})

	t.Run("Eligible_Months", func(t *testing.T) {
		rule := VaccineRule{MinAgeMonthsAtFirstDose: 2}
		msg, _, tags := engine.getAgeStatusAndEarliestDate(rule, "Test")
		if !containsStr(tags, "eligible") {
			t.Errorf("Expected eligible tag, got %v", tags)
		}
		if msg != "đủ điều kiện tuổi" {
			t.Errorf("Expected msg 'đủ điều kiện tuổi', got %q", msg)
		}
	})

	t.Run("TooYoung_Weeks", func(t *testing.T) {
		rule := VaccineRule{MinAgeWeeksAtFirstDose: 10}
		msg, _, tags := engine.getAgeStatusAndEarliestDate(rule, "Test")
		if !containsStr(tags, "too_young") {
			t.Errorf("Expected too_young tag, got %v", tags)
		}
		if msg != "cần 10 tuần tuổi" {
			t.Errorf("Expected msg 'cần 10 tuần tuổi', got %q", msg)
		}
	})

	t.Run("TooYoung_Years", func(t *testing.T) {
		rule := VaccineRule{MinAgeYearsAtFirstDose: 1}
		msg, _, tags := engine.getAgeStatusAndEarliestDate(rule, "Test")
		if !containsStr(tags, "too_young") {
			t.Errorf("Expected too_young tag, got %v", tags)
		}
		if msg != "cần 1 tuổi" {
			t.Errorf("Expected msg 'cần 1 tuổi', got %q", msg)
		}
	})
}

func TestCheckFirstDoseAgeValidity(t *testing.T) {
	dob, _ := ParseDateDDMMYYYY("01/01/2024")
	
	engine := &Engine{
		DOB: dob,
	}

	t.Run("Valid", func(t *testing.T) {
		firstDose, _ := ParseDateDDMMYYYY("01/03/2024") // 2 months
		rule := VaccineRule{MinAgeMonthsAtFirstDose: 2}
		valid, errRes := engine.checkFirstDoseAgeValidity(firstDose, rule, "Test")
		if !valid {
			t.Errorf("Expected valid, got invalid: %v", errRes.Description)
		}
	})

	t.Run("TooEarly_Months", func(t *testing.T) {
		firstDose, _ := ParseDateDDMMYYYY("01/02/2024") // 1 month
		rule := VaccineRule{MinAgeMonthsAtFirstDose: 2}
		valid, errRes := engine.checkFirstDoseAgeValidity(firstDose, rule, "Test")
		if valid {
			t.Errorf("Expected invalid, got valid")
		}
		if errRes == nil || !containsStr(errRes.StatusTags, "too_early") {
			t.Errorf("Expected too_early tag, got %v", errRes.StatusTags)
		}
	})

	t.Run("TooEarly_Weeks", func(t *testing.T) {
		firstDose, _ := ParseDateDDMMYYYY("15/01/2024") // 2 weeks
		rule := VaccineRule{MinAgeWeeksAtFirstDose: 6}
		valid, errRes := engine.checkFirstDoseAgeValidity(firstDose, rule, "Test")
		if valid {
			t.Errorf("Expected invalid, got valid")
		}
		if errRes == nil || !containsStr(errRes.StatusTags, "too_early") {
			t.Errorf("Expected too_early tag, got %v", errRes.StatusTags)
		}
	})
}

func TestRulesPreprocessing_NamesNormGroup(t *testing.T) {
	// Mock rules with RawNamesMembers
	rules := map[string]VaccineRule{
		"GroupA": {
			RawNamesMembers: map[string][]string{
				"M1": {"Vaccine A (1)", "Vaccine A (2)"},
				"M2": {"Vaccine B"},
			},
		},
		"CollectiveRule": {
			Courses: []Course{
				{RawNames: []string{"Course1 Vaccine"}},
				{RawNames: []string{"Course2 Vaccine"}},
			},
		},
	}

	// We need to write to a temp file or just test the logic directly if possible.
	// Since NewEngine reads from file, let's create a temp rules.json
	// Actually, let's just test the logic by calling NewEngine with a real file if available
	// or just trust the logic since I already implemented it and it's straightforward.
	
	// I'll skip temp file for now to keep it simple, or I can use t.TempDir()
	tempDir := t.TempDir()
	rulesFile := tempDir + "/rules.json"
	
	// Need to marshal it
	data, _ := json.Marshal(rules)
	os.WriteFile(rulesFile, data, 0644)
	
	engine, err := NewEngine(rulesFile, time.Now(), time.Now())
	if err != nil {
		t.Fatalf("Failed to create engine: %v", err)
	}
	
	ruleA := engine.Rules["GroupA"]
	if len(ruleA.NamesNormGroup) != 2 {
		t.Errorf("Expected 2 NamesNormGroup, got %d: %v", len(ruleA.NamesNormGroup), ruleA.NamesNormGroup)
	}
	
	collRule := engine.Rules["CollectiveRule"]
	if len(collRule.NamesNorm) != 2 {
		t.Errorf("Expected 2 NamesNorm for CollectiveRule, got %d: %v", len(collRule.NamesNorm), collRule.NamesNorm)
	}
}
