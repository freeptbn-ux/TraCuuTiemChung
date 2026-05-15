package analyzer

import (
	"testing"
	"time"
	"vercel-backend/pkg/models"
)

func ptrInt(i int) *int {
	return &i
}

func TestSingleSeries_Booster_Due(t *testing.T) {
	dob := time.Date(2010, 1, 1, 0, 0, 0, 0, time.UTC)
	analysisDate := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	engine := &Engine{
		DOB:          dob,
		AnalysisDate: analysisDate,
		Rules:        make(map[string]VaccineRule),
	}

	rule := VaccineRule{
		DisplayName:            "Td (Uốn ván - Bạch hầu)",
		DosesRequired:          3,
		BoosterIntervalYears:   10,
		BoosterAfterDoseNumber: 3,
		NamesNorm:              []string{"td"},
	}

	// 3 doses injected, last dose 12 years ago
	lastDoseDate := time.Date(2012, 1, 1, 0, 0, 0, 0, time.UTC)
	administeredMap := map[string][]models.VaccineRecord{
		"td": {
			{Date: time.Date(2011, 1, 1, 0, 0, 0, 0, time.UTC), VaccineName: "Td"},
			{Date: time.Date(2011, 3, 1, 0, 0, 0, 0, time.UTC), VaccineName: "Td"},
			{Date: lastDoseDate, VaccineName: "Td"},
		},
	}

	results := engine.checkSingleSeries("Td", rule, administeredMap)
	if len(results) == 0 {
		t.Fatal("Expected booster result, got none")
	}
	if !containsStr(results[0].StatusTags, "booster_due") {
		t.Errorf("Expected booster_due tag, got %v", results[0].StatusTags)
	}
}

func TestSingleSeries_Booster_Upcoming(t *testing.T) {
	dob := time.Date(2010, 1, 1, 0, 0, 0, 0, time.UTC)
	analysisDate := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
	engine := &Engine{
		DOB:          dob,
		AnalysisDate: analysisDate,
		Rules:        make(map[string]VaccineRule),
	}

	rule := VaccineRule{
		DisplayName:            "Td",
		DosesRequired:          3,
		BoosterIntervalYears:   10,
		BoosterAfterDoseNumber: 3,
		NamesNorm:              []string{"td"},
	}

	// 3 doses injected, last dose 5 years ago
	lastDoseDate := time.Date(2015, 1, 1, 0, 0, 0, 0, time.UTC)
	administeredMap := map[string][]models.VaccineRecord{
		"td": {
			{Date: time.Date(2014, 1, 1, 0, 0, 0, 0, time.UTC), VaccineName: "Td"},
			{Date: time.Date(2014, 3, 1, 0, 0, 0, 0, time.UTC), VaccineName: "Td"},
			{Date: lastDoseDate, VaccineName: "Td"},
		},
	}

	results := engine.checkSingleSeries("Td", rule, administeredMap)
	if len(results) == 0 {
		t.Fatal("Expected booster upcoming result, got none")
	}
	if !containsStr(results[0].StatusTags, "booster_upcoming") {
		t.Errorf("Expected booster_upcoming tag, got %v", results[0].StatusTags)
	}
}

func TestSingleSeries_Booster_MaxAge(t *testing.T) {
	dob := time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)
	analysisDate := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC) // 25 years old
	engine := &Engine{
		DOB:          dob,
		AnalysisDate: analysisDate,
		Rules:        make(map[string]VaccineRule),
	}

	rule := VaccineRule{
		DisplayName:            "Td",
		DosesRequired:          3,
		BoosterIntervalYears:   10,
		BoosterMaxAgeYears:     20, // Only booster until 20 years old
		NamesNorm:              []string{"td"},
	}

	administeredMap := map[string][]models.VaccineRecord{
		"td": {
			{Date: time.Date(2001, 1, 1, 0, 0, 0, 0, time.UTC), VaccineName: "Td"},
			{Date: time.Date(2001, 3, 1, 0, 0, 0, 0, time.UTC), VaccineName: "Td"},
			{Date: time.Date(2001, 5, 1, 0, 0, 0, 0, time.UTC), VaccineName: "Td"},
		},
	}

	results := engine.checkSingleSeries("Td", rule, administeredMap)
	if len(results) != 0 {
		t.Errorf("Expected no results (max age exceeded), got %d results", len(results))
	}
}

func TestSingleSeries_MVVAC_CoveredByMMR(t *testing.T) {
	dob := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)
	analysisDate := time.Date(2024, 6, 1, 0, 0, 0, 0, time.UTC)
	engine := &Engine{
		DOB:          dob,
		AnalysisDate: analysisDate,
		Rules: map[string]VaccineRule{
			"MMR": {
				DisplayName:              "MMR",
				ProvidesMeaslesProtection: true,
				NamesNorm:                []string{"mmr"},
			},
		},
	}

	rule := VaccineRule{
		DisplayName:   "MVVAC (Sởi)",
		DosesRequired: 1,
		NamesNorm:     []string{"mvvac"},
	}

	// Injected MMR at 13 months
	administeredMap := map[string][]models.VaccineRecord{
		"mmr": {{Date: time.Date(2024, 2, 1, 0, 0, 0, 0, time.UTC), VaccineName: "MMR"}},
	}

	results := engine.checkSingleSeries("MVVAC", rule, administeredMap)
	if len(results) == 0 {
		t.Fatal("Expected result for MVVAC, got none")
	}
	if !containsStr(results[0].StatusTags, "coverage_by_other") {
		t.Errorf("Expected coverage_by_other tag, got %v", results[0].StatusTags)
	}
}

func TestSingleSeries_Mengoc_ReverseInteraction(t *testing.T) {
	dob := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
	analysisDate := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC) // 4 years old (> 24 months)
	engine := &Engine{
		DOB:          dob,
		AnalysisDate: analysisDate,
		Rules: map[string]VaccineRule{
			"MenQuadfi": {
				DisplayName: "MenQuadfi",
				NamesNorm:   []string{"menquadfi"},
			},
		},
	}

	rule := VaccineRule{
		DisplayName:   "VA-MENGOC-BC",
		DosesRequired: 2,
		NamesNorm:     []string{"vamengocbc"},
	}

	administeredMap := map[string][]models.VaccineRecord{
		"menquadfi": {{Date: time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC), VaccineName: "MenQuadfi"}},
	}

	results := engine.checkSingleSeries("VA-MENGOC-BC", rule, administeredMap)
	if len(results) == 0 {
		t.Fatal("Expected result, got none")
	}
	if !containsStr(results[0].StatusTags, "interaction_reverse_mengoc") {
		t.Errorf("Expected interaction_reverse_mengoc tag, got %v", results[0].StatusTags)
	}
}

func TestSingleSeries_FirstDose_TooEarly_Months(t *testing.T) {
	dob := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	analysisDate := time.Date(2024, 3, 1, 0, 0, 0, 0, time.UTC)
	engine := &Engine{
		DOB:          dob,
		AnalysisDate: analysisDate,
		Rules:        make(map[string]VaccineRule),
	}

	rule := VaccineRule{
		DisplayName:             "DPT",
		DosesRequired:           3,
		MinAgeMonthsAtFirstDose: 2,
		NamesNorm:               []string{"dpt"},
	}

	// Injected at 1 month (Too early)
	administeredMap := map[string][]models.VaccineRecord{
		"dpt": {{Date: time.Date(2024, 2, 1, 0, 0, 0, 0, time.UTC), VaccineName: "DPT"}},
	}

	results := engine.checkSingleSeries("DPT", rule, administeredMap)
	if len(results) == 0 {
		t.Fatal("Expected error result, got none")
	}
	if !containsStr(results[0].StatusTags, "too_early") {
		t.Errorf("Expected too_early tag, got %v", results[0].StatusTags)
	}
}

func TestSingleSeries_NoDoses_Eligible(t *testing.T) {
	dob := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	analysisDate := time.Date(2024, 7, 1, 0, 0, 0, 0, time.UTC) // 6 months
	engine := &Engine{
		DOB:          dob,
		AnalysisDate: analysisDate,
		Rules:        make(map[string]VaccineRule),
	}

	rule := VaccineRule{
		DisplayName:             "DPT",
		DosesRequired:           3,
		MinAgeMonthsAtFirstDose: 2,
		NamesNorm:               []string{"dpt"},
	}

	results := engine.checkSingleSeries("DPT", rule, nil)
	if len(results) == 0 {
		t.Fatal("Expected result, got none")
	}
	if !containsStr(results[0].StatusTags, "eligible") {
		t.Errorf("Expected eligible tag, got %v", results[0].StatusTags)
	}
}

func TestSingleSeries_NoDoses_TooOld(t *testing.T) {
	dob := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
	analysisDate := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	engine := &Engine{
		DOB:          dob,
		AnalysisDate: analysisDate,
		Rules:        make(map[string]VaccineRule),
	}

	rule := VaccineRule{
		DisplayName:             "Rota",
		DosesRequired:           2,
		MaxAgeMonthsAtFirstDose: 6, // Rota usually has tight deadline
		NamesNorm:               []string{"rota"},
	}

	results := engine.checkSingleSeries("Rota", rule, nil)
	if len(results) == 0 {
		t.Fatal("Expected result, got none")
	}
	// Note: Parity expects eligible for 0-dose single series even if max age is exceeded
	if !containsStr(results[0].StatusTags, "eligible") {
		t.Errorf("Expected eligible tag, got %v", results[0].StatusTags)
	}
}

func TestSingleSeries_DoseSpecific_AltAge(t *testing.T) {
	dob := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
	analysisDate := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)
	engine := &Engine{
		DOB:          dob,
		AnalysisDate: analysisDate,
		Rules:        make(map[string]VaccineRule),
	}

	rule := VaccineRule{
		DisplayName:   "DPT Booster",
		DosesRequired: 4,
		MinIntervalDays: []*int{nil, ptrInt(30), ptrInt(30), ptrInt(30)},
		DoseSpecificRules: map[string]DoseRule{
			"4": {MinAbsoluteAgeMonths: 48},
		},
		NamesNorm: []string{"dpt"},
	}

	// 3 doses injected
	administeredMap := map[string][]models.VaccineRecord{
		"dpt": {
			{Date: time.Date(2020, 3, 1, 0, 0, 0, 0, time.UTC)},
			{Date: time.Date(2020, 4, 1, 0, 0, 0, 0, time.UTC)},
			{Date: time.Date(2020, 5, 1, 0, 0, 0, 0, time.UTC)},
		},
	}

	results := engine.checkSingleSeries("DPT", rule, administeredMap)
	expectedDate := AddYears(dob, 4)
	if !results[0].EarliestNextDoseDate.Equal(expectedDate) {
		t.Errorf("Expected earliest date %v, got %v", expectedDate, results[0].EarliestNextDoseDate)
	}
}

func TestAgeDep_NoDoses_Status(t *testing.T) {
	dob := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	analysisDate := time.Date(2024, 2, 1, 0, 0, 0, 0, time.UTC) // 1 month
	engine := &Engine{
		DOB:          dob,
		AnalysisDate: analysisDate,
		Rules:        make(map[string]VaccineRule),
	}

	rule := VaccineRule{
		DisplayName: "Prevenar13",
		Type:        RuleTypeAgeDependent,
		RulesByAge: []AgeRule{
			{MinAgeMonthsAtFirstDose: 2, MaxAgeMonthsAtFirstDose: 6, DosesRequired: 4},
			{MinAgeMonthsAtFirstDose: 7, MaxAgeMonthsAtFirstDose: 11, DosesRequired: 3},
		},
		MinAgeMonthsOverall: 2,
		NamesNorm:           []string{"prevenar13"},
	}

	results := engine.checkAgeDependentSeries("Prevenar13", rule, nil)
	if len(results) == 0 {
		t.Fatal("Expected result, got none")
	}
	// At 1 month, it should match the 2-6 months rule as earliest possible? 
	// Wait, at 1 month it doesn't match any rule yet if we strictly check current age.
	// In Python, if 0 doses, it checks current age.
	if !containsStr(results[0].StatusTags, "too_young") {
		t.Errorf("Expected too_young tag at 1 month, got %v", results[0].StatusTags)
	}
}

func TestAgeDep_FirstDose_At2Months(t *testing.T) {
	dob := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	analysisDate := time.Date(2024, 3, 15, 0, 0, 0, 0, time.UTC)
	engine := &Engine{
		DOB:          dob,
		AnalysisDate: analysisDate,
		Rules:        make(map[string]VaccineRule),
	}

	rule := VaccineRule{
		DisplayName: "Prevenar13",
		Type:        RuleTypeAgeDependent,
		RulesByAge: []AgeRule{
			{MinAgeMonthsAtFirstDose: 2, MaxAgeMonthsAtFirstDose: 6, DosesRequired: 4},
		},
		NamesNorm: []string{"prevenar13"},
	}

	// 1 dose at 2 months
	administeredMap := map[string][]models.VaccineRecord{
		"prevenar13": {{Date: time.Date(2024, 3, 1, 0, 0, 0, 0, time.UTC)}},
	}

	results := engine.checkAgeDependentSeries("Prevenar13", rule, administeredMap)
	if len(results) == 0 {
		t.Fatal("Expected result, got none")
	}
	if results[0].DoseNumber != 2 {
		t.Errorf("Expected dose 2 recommendation, got %d", results[0].DoseNumber)
	}
}
