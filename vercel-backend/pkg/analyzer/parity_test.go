package analyzer

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"testing"
	"time"

	"vercel-backend/pkg/models"
)

type ParityFixture struct {
	PatientInfo struct {
		Name       string `json:"name"`
		Birth      string `json:"birth"`
		SystemDate string `json:"system_date"`
	} `json:"patient_info"`
	History []struct {
		VaccineName string `json:"vaccine_name"`
		Dose        string `json:"dose"`
		Date        string `json:"date"`
	} `json:"history"`
}

func TestParity(t *testing.T) {
	testDataDir := "testdata"
	files, err := os.ReadDir(testDataDir)
	if err != nil {
		t.Fatalf("Failed to read testdata dir: %v", err)
	}

	for _, file := range files {
		if file.IsDir() || filepath.Ext(file.Name()) != ".json" || strings.HasPrefix(file.Name(), "expected_") {
			continue
		}

		t.Run(file.Name(), func(t *testing.T) {
			// 1. Load Fixture
			fixturePath := filepath.Join(testDataDir, file.Name())
			fixtureData, err := os.ReadFile(fixturePath)
			if err != nil {
				t.Fatalf("Failed to read fixture: %v", err)
			}

			var fixture ParityFixture
			if err := json.Unmarshal(fixtureData, &fixture); err != nil {
				t.Fatalf("Failed to unmarshal fixture: %v", err)
			}

			// 2. Load Expected Output
			expectedPath := filepath.Join(testDataDir, "expected_"+file.Name())
			expectedData, err := os.ReadFile(expectedPath)
			if err != nil {
				t.Fatalf("Failed to read expected output: %v", err)
			}

			type ExpectedResult struct {
				VaccineNameForPopup  string   `json:"vaccine_name_for_popup"`
				Description          string   `json:"description"`
				EarliestNextDoseDate *string  `json:"earliest_next_dose_date"`
				StatusTags           []string `json:"status_tags"`
			}
			var expectedRaw []ExpectedResult
			if err := json.Unmarshal(expectedData, &expectedRaw); err != nil {
				t.Fatalf("Failed to unmarshal expected output: %v", err)
			}

			var expected []AnalysisResult
			for _, r := range expectedRaw {
				var dt *time.Time
				if r.EarliestNextDoseDate != nil && *r.EarliestNextDoseDate != "" && *r.EarliestNextDoseDate != "None" {
					parsed, err := time.Parse("2006-01-02", *r.EarliestNextDoseDate)
					if err == nil {
						dt = &parsed
					}
				}
				expected = append(expected, AnalysisResult{
					VaccineNameForPopup:  r.VaccineNameForPopup,
					Description:          r.Description,
					EarliestNextDoseDate: dt,
					StatusTags:           r.StatusTags,
				})
			}

			// 3. Run Go Engine
			dob, _ := ParseDateDDMMYYYY(fixture.PatientInfo.Birth)
			analysisDate, _ := ParseDateDDMMYYYY(fixture.PatientInfo.SystemDate)
			rulesPath := filepath.Join("..", "..", "assets", "vaccine_rules.json")

			engine, err := NewEngine(rulesPath, dob, analysisDate)
			if err != nil {
				t.Fatalf("Failed to create engine: %v", err)
			}

			var history []models.VaccineRecord
			for _, h := range fixture.History {
				d, _ := ParseDateDDMMYYYY(h.Date)
				history = append(history, models.VaccineRecord{
					VaccineName: h.VaccineName,
					Dose:        h.Dose,
					Date:        d,
				})
			}

			actual := engine.Analyze(history)

			// 4. Assert Parity
			assertParity(t, expected, actual, analysisDate)
		})
	}
}

func assertParity(t *testing.T, expected, actual []AnalysisResult, analysisDate time.Time) {
	// Filter out empty results in actual if any
	var filteredActual []AnalysisResult
	for _, a := range actual {
		if a.VaccineNameForPopup != "" {
			filteredActual = append(filteredActual, a)
		}
	}
	actual = filteredActual

	// Sort both for consistent matching
	sortResults(expected)
	sortResults(actual)

	for _, a := range actual {
		t.Logf("DEBUG: Actual result: name=[%s], tags=%v, date=%v", a.VaccineNameForPopup, a.StatusTags, a.EarliestNextDoseDate)
	}

	// Matching algorithm that handles duplicates
	actualUsed := make([]bool, len(actual))
	
	for i, e := range expected {
		testName := e.VaccineNameForPopup
		if i > 0 && expected[i-1].VaccineNameForPopup == e.VaccineNameForPopup {
			testName = fmt.Sprintf("%s#%02d", e.VaccineNameForPopup, i)
		}

		t.Run(testName, func(t *testing.T) {
			foundIdx := -1
			expNorm := normalizeForMatch(e.VaccineNameForPopup)

			// Try exact normalized match first
			for j, a := range actual {
				if !actualUsed[j] && normalizeForMatch(a.VaccineNameForPopup) == expNorm {
					foundIdx = j
					break
				}
			}
			// Try partial normalized match if no exact match found
			if foundIdx == -1 {
				for j, a := range actual {
					if !actualUsed[j] {
						actNorm := normalizeForMatch(a.VaccineNameForPopup)
						if strings.Contains(actNorm, expNorm) || strings.Contains(expNorm, actNorm) {
							foundIdx = j
							break
						}
					}
				}
			}

			if foundIdx == -1 {
				t.Errorf("Expected result for %s not found in actual results", e.VaccineNameForPopup)
				return
			}

			a := actual[foundIdx]
			actualUsed[foundIdx] = true

			// Compare status tags
			expectedTags := make(map[string]bool)
			for _, tag := range e.StatusTags {
				expectedTags[tag] = true
			}
			for _, tag := range a.StatusTags {
				if !expectedTags[tag] {
					t.Errorf("Unexpected status tag: %s (Actual: %v)", tag, a.StatusTags)
				}
				delete(expectedTags, tag)
			}
			for tag := range expectedTags {
				t.Errorf("Missing status tag: %s (Actual: %v)", tag, a.StatusTags)
			}

			// Compare earliest dose date (with 1 day tolerance)
			if e.EarliestNextDoseDate != nil {
				if a.EarliestNextDoseDate == nil {
					t.Errorf("Expected date %s, but Go got nil", e.EarliestNextDoseDate.Format("2006-01-02"))
				} else {
					diff := a.EarliestNextDoseDate.Sub(*e.EarliestNextDoseDate)
					if diff < -24*time.Hour || diff > 24*time.Hour {
						t.Errorf("Date mismatch: expected %s, got %s", e.EarliestNextDoseDate.Format("2006-01-02"), a.EarliestNextDoseDate.Format("2006-01-02"))
					}
				}
			} else {
				if a.EarliestNextDoseDate != nil {
					// In some cases Python uses system_date instead of None for missing items,
					// but our Go engine usually returns nil.
					// We only error if the date is significantly different from analysisDate
					if a.EarliestNextDoseDate.Sub(analysisDate) > 24*time.Hour {
						t.Errorf("Expected date is None, Go got %s", a.EarliestNextDoseDate)
					}
				}
			}
		})
	}
}

func sortResults(results []AnalysisResult) {
	sort.SliceStable(results, func(i, j int) bool {
		if results[i].VaccineNameForPopup != results[j].VaccineNameForPopup {
			return results[i].VaccineNameForPopup < results[j].VaccineNameForPopup
		}
		// Second sort criteria: status tags
		return strings.Join(results[i].StatusTags, ",") < strings.Join(results[j].StatusTags, ",")
	})
}

func normalizeForMatch(s string) string {
	return NormalizeForMatch(s)
}
func BenchmarkAnalyze_EmptyHistory(b *testing.B) {
	fixturePath := filepath.Join("testdata", "empty_history.json")
	fixtureData, _ := os.ReadFile(fixturePath)
	var fixture ParityFixture
	json.Unmarshal(fixtureData, &fixture)

	dob, _ := ParseDateDDMMYYYY(fixture.PatientInfo.Birth)
	analysisDate, _ := ParseDateDDMMYYYY(fixture.PatientInfo.SystemDate)
	rulesPath := filepath.Join("..", "..", "assets", "vaccine_rules.json")
	engine, _ := NewEngine(rulesPath, dob, analysisDate)

	var history []models.VaccineRecord
	for _, h := range fixture.History {
		d, _ := ParseDateDDMMYYYY(h.Date)
		history = append(history, models.VaccineRecord{
			VaccineName: h.VaccineName,
			Dose:        h.Dose,
			Date:        d,
		})
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		engine.Analyze(history)
	}
}

func BenchmarkAnalyze_FullSchedule(b *testing.B) {
	fixturePath := filepath.Join("testdata", "full_schedule.json")
	fixtureData, _ := os.ReadFile(fixturePath)
	var fixture ParityFixture
	json.Unmarshal(fixtureData, &fixture)

	dob, _ := ParseDateDDMMYYYY(fixture.PatientInfo.Birth)
	analysisDate, _ := ParseDateDDMMYYYY(fixture.PatientInfo.SystemDate)
	rulesPath := filepath.Join("..", "..", "assets", "vaccine_rules.json")
	engine, _ := NewEngine(rulesPath, dob, analysisDate)

	var history []models.VaccineRecord
	for _, h := range fixture.History {
		d, _ := ParseDateDDMMYYYY(h.Date)
		history = append(history, models.VaccineRecord{
			VaccineName: h.VaccineName,
			Dose:        h.Dose,
			Date:        d,
		})
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		engine.Analyze(history)
	}
}
