package analyzer

import (
	"testing"
	"time"
	"vercel-backend/pkg/models"
)

func TestMMR_MVVAC_Logic(t *testing.T) {
	dob := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	// Để test trạng thái "info", analysisDate phải trước ngày tiêm sớm nhất (2024-12-24)
	analysisDate := time.Date(2024, 11, 1, 0, 0, 0, 0, time.UTC)
	engine, err := NewEngine("../../assets/vaccine_rules.json", dob, analysisDate)
	if err != nil {
		t.Fatalf("Failed to create engine: %v", err)
	}

	mmrDisplayName := "Vắc xin Sởi-Quai bị-Rubella (MMR-II/Priorix)"

	// Test 1: MVVAC tiêm lúc 9 tháng, chưa tiêm MMR
	history1 := []models.VaccineRecord{
		{VaccineName: "MVVAC", Date: dob.AddDate(0, 9, 0)},
	}
	results1 := engine.Analyze(history1)
	foundMMR1 := false
	for _, res := range results1 {
		if res.VaccineNameForPopup == mmrDisplayName {
			foundMMR1 = true
			expectedEarliest := history1[0].Date.AddDate(0, 0, 84)
			if !res.EarliestNextDoseDate.Equal(expectedEarliest) {
				t.Errorf("Test 1: Expected earliest MMR date %s, got %s", expectedEarliest, res.EarliestNextDoseDate)
			}
			if !containsStr(res.StatusTags, "info") {
				t.Errorf("Test 1: Expected tag 'info', got %v", res.StatusTags)
			}
		}
	}
	if !foundMMR1 {
		t.Error("Test 1: MMR recommendation not found")
	}

	// Test 2: MVVAC -> MMR cách 90 ngày (Good Interval)
	history2 := []models.VaccineRecord{
		{VaccineName: "MVVAC", Date: dob.AddDate(0, 9, 0)},
		{VaccineName: "MMR-II", Date: dob.AddDate(0, 9, 0).AddDate(0, 0, 90)},
	}
	results2 := engine.Analyze(history2)
	for _, res := range results2 {
		if res.VaccineNameForPopup == mmrDisplayName {
			if containsStr(res.StatusTags, "warning") {
				t.Error("Test 2: Unexpected warning for good interval")
			}
		}
	}

	// Test 3: MVVAC -> MMR cách 60 ngày (Bad Interval)
	history3 := []models.VaccineRecord{
		{VaccineName: "MVVAC", Date: dob.AddDate(0, 9, 0)},
		{VaccineName: "MMR-II", Date: dob.AddDate(0, 9, 0).AddDate(0, 0, 60)},
	}
	results3 := engine.Analyze(history3)
	foundWarning3 := false
	for _, res := range results3 {
		if res.VaccineNameForPopup == mmrDisplayName && containsStr(res.StatusTags, "warning") {
			foundWarning3 = true
		}
	}
	if !foundWarning3 {
		t.Error("Test 3: Expected warning for bad interval (60 days)")
	}
}

func TestFlu_Logic(t *testing.T) {
	dob := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	analysisDate := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	engine, err := NewEngine("../../assets/vaccine_rules.json", dob, analysisDate)
	if err != nil {
		t.Fatalf("Failed to create engine: %v", err)
	}

	fluDisplayName := "Vắc xin Cúm"

	// Test 1: Trẻ 8 tháng, chưa tiêm gì -> Recommend mũi 1
	results1 := engine.Analyze(nil)
	foundFlu1 := false
	for _, res := range results1 {
		if res.VaccineNameForPopup == fluDisplayName {
			foundFlu1 = true
			if res.DoseNumber != 1 {
				t.Errorf("Test 1: Expected dose 1, got %d", res.DoseNumber)
			}
		}
	}
	if !foundFlu1 {
		t.Error("Test 1: Flu recommendation not found")
	}

	// Test 2: 1 mũi lúc 2 tuổi (giả sử dob lùi lại) -> Recommend mũi 2 (30 ngày theo rules.json)
	dob2 := time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC)
	engine2, err := NewEngine("../../assets/vaccine_rules.json", dob2, analysisDate)
	if err != nil {
		t.Fatalf("Failed to create engine2: %v", err)
	}
	history2 := []models.VaccineRecord{
		{VaccineName: "Vaxigrip Tetra", Date: dob2.AddDate(2, 0, 0)},
	}
	results2 := engine2.Analyze(history2)
	foundFlu2 := false
	for _, res := range results2 {
		if res.VaccineNameForPopup == fluDisplayName {
			foundFlu2 = true
			if res.DoseNumber != 2 {
				t.Errorf("Test 2: Expected dose 2, got %d", res.DoseNumber)
			}
			expectedDate := analysisDate // Adjusted to analysisDate because 2024-01-31 < 2025-01-01
			if !res.EarliestNextDoseDate.Equal(expectedDate) {
				t.Errorf("Test 2: Expected earliest date %s, got %s", expectedDate, res.EarliestNextDoseDate)
			}
		}
	}
	if !foundFlu2 {
		t.Error("Test 2: Flu dose 2 recommendation not found")
	}

	// Test 3: Keyword match raw name
	history3 := []models.VaccineRecord{
		{VaccineName: "Vaxigrip Tetra (Cúm A/B)", Date: dob.AddDate(0, 6, 0)},
	}
	results3 := engine.Analyze(history3)
	foundFlu3 := false
	for _, res := range results3 {
		if res.VaccineNameForPopup == fluDisplayName {
			foundFlu3 = true
			if res.DoseNumber != 2 {
				t.Errorf("Test 3: Expected dose 2 after keyword match, got %d", res.DoseNumber)
			}
		}
	}
	if !foundFlu3 {
		t.Error("Test 3: Flu not matched by keyword in raw name")
	}
}

func TestMeningococcalACYW_Logic(t *testing.T) {
	dob := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	analysisDate := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	engine, err := NewEngine("../../assets/vaccine_rules.json", dob, analysisDate)
	if err != nil {
		t.Fatalf("Failed to create engine: %v", err)
	}

	mqDisplayName := "MenQuadfi (Sanofi Pasteur)"
	maDisplayName := "Menactra (Sanofi Pasteur)"

	// Test 1: Trẻ 3 tháng -> MenQuadfi gợi ý
	results1 := engine.Analyze(nil)
	foundMQ1 := false
	for _, res := range results1 {
		if res.VaccineNameForPopup == mqDisplayName {
			foundMQ1 = true
		}
	}
	if !foundMQ1 {
		t.Error("Test 1: MenQuadfi suggestion not found for 3m infant")
	}

	// Test 2: Trẻ 9 tháng -> Cả 2 gợi ý
	analysisDate2 := dob.AddDate(0, 9, 0)
	engine2, err := NewEngine("../../assets/vaccine_rules.json", dob, analysisDate2)
	if err != nil {
		t.Fatalf("Failed to create engine2: %v", err)
	}
	results2 := engine2.Analyze(nil)
	foundMQ2 := false
	foundMA2 := false
	for _, res := range results2 {
		if res.VaccineNameForPopup == mqDisplayName {
			foundMQ2 = true
		}
		if res.VaccineNameForPopup == maDisplayName {
			foundMA2 = true
		}
	}
	if !foundMQ2 || !foundMA2 {
		t.Errorf("Test 2: Expected both MenQuadfi (%v) and Menactra (%v) for 9m infant", foundMQ2, foundMA2)
	}

	// Test 3: Interaction with VA-MENGOC-BC
	dob3 := analysisDate.AddDate(-3, 0, 0) // 3 years old
	engine3, _ := NewEngine("../../assets/vaccine_rules.json", dob3, analysisDate)
	history3 := []models.VaccineRecord{
		{VaccineName: "VA - MENGOC - BC", Date: analysisDate.AddDate(0, 0, -10)},
	}
	results3 := engine3.Analyze(history3)
	foundWarning3 := false
	for _, res := range results3 {
		if (res.VaccineNameForPopup == mqDisplayName || res.VaccineNameForPopup == maDisplayName) && containsStr(res.StatusTags, "warning") {
			foundWarning3 = true
		}
	}
	if !foundWarning3 {
		t.Error("Test 3: Expected interaction warning with VA-MENGOC-BC")
	}

	// Test 4: 60-day calendar months
	history4 := []models.VaccineRecord{
		{VaccineName: "MenQuadfi", Date: time.Date(2025, 1, 31, 0, 0, 0, 0, time.UTC)},
	}
	// MenQuadfi tiêm lúc 3 tháng (dob lùi lại)
	dob4 := history4[0].Date.AddDate(0, -3, 0)
	engine4, _ := NewEngine("../../assets/vaccine_rules.json", dob4, history4[0].Date.AddDate(0, 0, 1))
	results4 := engine4.Analyze(history4)
	for _, res := range results4 {
		if res.VaccineNameForPopup == mqDisplayName {
			// MenQuadfi tiêm lúc < 6 tháng cần 3 mũi cách 2 tháng (60 ngày)
			expectedDate := AddMonths(history4[0].Date, 2) // 2025-01-31 + 2m = 2025-03-31
			if !res.EarliestNextDoseDate.Equal(expectedDate) {
				t.Errorf("Test 4: Expected calendar month earliest date %s, got %s", expectedDate.Format("2006-01-02"), res.EarliestNextDoseDate.Format("2006-01-02"))
			}
		}
	}
}
