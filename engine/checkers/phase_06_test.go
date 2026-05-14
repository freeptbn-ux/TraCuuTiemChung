package checkers

import (
	"strings"
	"testing"
	"time"
	"tracuutiemchung-engine/models"
)

func TestSpecialPneumo_Interchange(t *testing.T) {
	dob := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	analysisDate := time.Date(2024, 9, 1, 0, 0, 0, 0, time.UTC)

	// Mock Rules
	prevenarRule := &models.VaccineRule{
		DisplayName: "Prevenar13",
		Type:        "age_dependent_series",
		NamesNorm:   []string{"prevenar13"},
		RulesByAge: []models.AgeRule{
			{
				MaxAgeAtFirstDoseMonths: intPtr(6),
				DosesRequired:           4,
				MinIntervalDays:         []*int{nil, intPtr(30), intPtr(30), intPtr(240)},
			},
		},
	}
	synflorixRule := &models.VaccineRule{
		DisplayName: "Synflorix",
		Type:        "age_dependent_series",
		NamesNorm:   []string{"synflorix"},
		RulesByAge: []models.AgeRule{
			{
				MaxAgeAtFirstDoseMonths: intPtr(6),
				DosesRequired:           4,
				MinIntervalDays:         []*int{nil, intPtr(30), intPtr(30), intPtr(180)},
			},
		},
	}

	rules := map[string]*models.VaccineRule{
		"Prevenar13": prevenarRule,
		"Synflorix":  synflorixRule,
	}

	// Case: 1 Synflorix (2 months) and 1 Prevenar13 (4 months)
	adminMap := map[string][]models.AdministeredDose{
		"synflorix": {
			{VaccineName: "synflorix", Date: time.Date(2024, 3, 1, 0, 0, 0, 0, time.UTC)},
		},
		"prevenar13": {
			{VaccineName: "prevenar13", Date: time.Date(2024, 5, 1, 0, 0, 0, 0, time.UTC)},
		},
	}

	results := CheckSpecialPneumo(adminMap, dob, analysisDate, rules)

	if len(results) == 0 {
		t.Fatal("Expected recommendations, got 0")
	}

	foundInterchange := false
	for _, res := range results {
		if strings.Contains(res.Description, "Cảnh báo: Đã ghi nhận tiêm xen kẽ") {
			foundInterchange = true
		}
		for _, tag := range res.StatusTags {
			if tag == "error_interchange" {
				foundInterchange = true
			}
		}
	}

	if !foundInterchange {
		t.Errorf("Expected 'error_interchange' tag or warning description, but not found")
	}
}

func TestAlternativeGroup_Rota(t *testing.T) {
	dob := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	analysisDate := time.Date(2024, 5, 1, 0, 0, 0, 0, time.UTC)

	rotaRule := &models.VaccineRule{
		DisplayName: "Rota",
		Type:        "group_alternative_courses_min_age",
		Courses: []models.Course{
			{
				Display:       "Rotarix",
				RawNames:      []string{"Rotarix"},
				NamesNorm:     []string{"rotarix"},
				DosesRequired: 2,
				MinIntervalDays: []*int{nil, intPtr(30)},
			},
			{
				Display:       "RotaTeq",
				RawNames:      []string{"RotaTeq"},
				NamesNorm:     []string{"rotateq"},
				DosesRequired: 3,
				MinIntervalDays: []*int{nil, intPtr(30), intPtr(30)},
			},
		},
	}

	// Case 1: No doses, should recommend first course (generic "Rota")
	results := CheckAlternativeCoursesGroup(rotaRule, nil, dob, analysisDate)
	if len(results) == 0 || results[0].VaccineName != "Rota" {
		t.Errorf("Expected Rota recommendation, got %v", results)
	}

	// Case 2: Started RotaTeq
	adminMap := map[string][]models.AdministeredDose{
		"rotateq": {
			{VaccineName: "RotaTeq", Date: time.Date(2024, 2, 1, 0, 0, 0, 0, time.UTC)},
		},
	}
	results = CheckAlternativeCoursesGroup(rotaRule, adminMap, dob, analysisDate)
	if len(results) == 0 || results[0].VaccineName != "RotaTeq" {
		t.Errorf("Expected RotaTeq recommendation, got %v", results)
	}
	if !strings.Contains(results[0].Description, "Mũi 2") {
		t.Errorf("Expected 'Mũi 2' in description, got %s", results[0].Description)
	}
}

func TestFluGroup(t *testing.T) {
	dob := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	analysisDate := time.Date(2025, 2, 1, 0, 0, 0, 0, time.UTC) // 13 months old

	fluRule := &models.VaccineRule{
		DisplayName: "Flu",
		Type:        "flu_group",
		RawNames:    []string{"Vaxigrip"},
		NamesNorm:   []string{"vaxigrip"},
	}

	// Case 1: Child > 6 months, no doses
	results := CheckFluGroup(fluRule, nil, dob, analysisDate)
	if len(results) == 0 || !strings.Contains(results[0].Description, "Mũi 1") {
		t.Errorf("Expected Flu dose 1 recommendation, got %v", results)
	}

	// Case 2: Child 13 months, 1 dose given at 12 months
	adminMap := map[string][]models.AdministeredDose{
		"vaxigrip": {
			{VaccineName: "Vaxigrip", Date: time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)},
		},
	}
	results = CheckFluGroup(fluRule, adminMap, dob, analysisDate)
	if len(results) == 0 || !strings.Contains(results[0].Description, "Mũi 2") {
		t.Errorf("Expected Flu dose 2 recommendation, got %v", results)
	}
}

func TestAlternativeGroup_Rota_Mixed(t *testing.T) {
	dob := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	analysisDate := time.Date(2024, 6, 1, 0, 0, 0, 0, time.UTC)

	rotaRule := &models.VaccineRule{
		DisplayName: "Rota",
		Type:        "group_alternative_courses_min_age",
		Courses: []models.Course{
			{
				Display:       "Rotarix",
				NamesNorm:     []string{"rotarix"},
				DosesRequired: 2,
				MinIntervalDays: []*int{nil, intPtr(30)},
			},
			{
				Display:       "RotaTeq",
				NamesNorm:     []string{"rotateq"},
				DosesRequired: 3,
				MinIntervalDays: []*int{nil, intPtr(30), intPtr(30)},
			},
		},
	}

	// Mixed: Dose 1 Rotarix, Dose 2 RotaTeq
	adminMap := map[string][]models.AdministeredDose{
		"rotarix": {
			{VaccineName: "rotarix", Date: time.Date(2024, 2, 1, 0, 0, 0, 0, time.UTC)},
		},
		"rotateq": {
			{VaccineName: "rotateq", Date: time.Date(2024, 3, 10, 0, 0, 0, 0, time.UTC)},
		},
	}

	results := CheckAlternativeCoursesGroup(rotaRule, adminMap, dob, analysisDate)
	if len(results) == 0 {
		t.Fatal("Expected results")
	}

	// Should follow RotaTeq (3 doses) because RotaTeq was used
	if results[0].VaccineName != "RotaTeq" {
		t.Errorf("Expected RotaTeq course, got %s", results[0].VaccineName)
	}
	if !strings.Contains(results[0].Description, "Mũi 3") {
		t.Errorf("Expected Dose 3 recommendation, got %s", results[0].Description)
	}
}

func TestFluGroup_Over9(t *testing.T) {
	dob := time.Date(2010, 1, 1, 0, 0, 0, 0, time.UTC) // 15 years old
	analysisDate := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)

	fluRule := &models.VaccineRule{
		DisplayName: "Flu",
		Type:        "flu_group",
		NamesNorm:   []string{"vaxigrip"},
	}

	// Case 1: Over 9 years, no doses
	results := CheckFluGroup(fluRule, nil, dob, analysisDate)
	if len(results) == 0 || !strings.Contains(results[0].Description, "Mũi 1") {
		t.Errorf("Expected Flu dose 1, got %v", results)
	}

	// Case 2: Over 9 years, 1 dose given
	adminMap := map[string][]models.AdministeredDose{
		"vaxigrip": {
			{VaccineName: "vaxigrip", Date: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)},
		},
	}
	results = CheckFluGroup(fluRule, adminMap, dob, analysisDate)
	if len(results) == 0 || !strings.Contains(results[0].Description, "tiêm nhắc lại hàng năm") {
		t.Errorf("Expected annual booster, got %v", results)
	}
}

func TestCumulativeGroup_MMR(t *testing.T) {
	dob := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	analysisDate := time.Date(2025, 10, 1, 0, 0, 0, 0, time.UTC)

	mmrRule := &models.VaccineRule{
		DisplayName: "MMR",
		Type:        "age_dependent_series",
		NamesNormGroup: []string{"mmr2", "priorix"},
		RulesByAge: []models.AgeRule{
			{
				MinAgeAtFirstDoseMonths: intPtr(9),
				DosesRequired:           2,
				MinIntervalDays:         []*int{nil, intPtr(90)},
			},
		},
	}

	// Patient has 1 dose of MMR II
	adminMap := map[string][]models.AdministeredDose{
		"mmr2": {
			{VaccineName: "mmr2", Date: time.Date(2024, 10, 1, 0, 0, 0, 0, time.UTC)},
		},
	}

	results := CheckCumulativeGroup(mmrRule, adminMap, dob, analysisDate)
	if len(results) == 0 {
		t.Fatal("Expected results")
	}
	if !strings.Contains(results[0].Description, "Mũi 2") {
		t.Errorf("Expected Dose 2 recommendation, got %s", results[0].Description)
	}
}



