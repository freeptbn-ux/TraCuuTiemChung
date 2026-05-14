package parity

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"
	"time"

	"tracuutiemchung-engine/engine"
	"tracuutiemchung-engine/engine/utils"
	"tracuutiemchung-engine/models"
)

type PythonOutput struct {
	MissingVaccines []PythonMissingItem `json:"missing_vaccines"`
}

type PythonMissingItem struct {
	VaccineName string   `json:"vaccine_name"`
	Status      string   `json:"status"`
	NextDose    *string  `json:"next_dose"`
	Message     string   `json:"message"`
	StatusTags  []string `json:"status_tags"`
}

func TestMinhKhoiParity(t *testing.T) {
	// 1. Load Rules
	rules, err := engine.LoadVaccineRules("../../vercel-backend/assets/vaccine_rules.json")
	if err != nil {
		t.Fatalf("Failed to load rules: %v", err)
	}

	// 2. Parse HTML
	patient, records, err := ParseMinhkhoiHTML("../../test/Minhkhoi.html")
	if err != nil {
		t.Fatalf("Failed to parse html: %v", err)
	}

	dob, _ := time.Parse("02/01/2006", patient.Birth)
	analysisDate, _ := time.Parse("02/01/2006", patient.SystemDate)

	fmt.Printf("Parsed DOB: %v, AnalysisDate: %v\n", dob.Format("02/01/2006"), analysisDate.Format("02/01/2006"))

	adminMap := make(map[string][]models.AdministeredDose)
	for _, rec := range records {
		norm := utils.NormalizeVaccineName(rec.RawName)
		adminMap[norm] = append(adminMap[norm], models.AdministeredDose{
			VaccineName: rec.RawName,
			Date:        rec.Date,
		})
	}

	patientRecord := models.PatientRecord{
		BirthDate:       dob,
		AdministeredMap: adminMap,
	}

	// 3. Run Engine
	goResults := engine.ProcessAllRules(rules, patientRecord, analysisDate)

	// Save Go output to JSON for parity
	goOutFile, _ := os.Create("../../testdata/go_output.json")
	json.NewEncoder(goOutFile).Encode(models.AnalysisResult{MissingItems: goResults})
	goOutFile.Close()

	// 4. Load Python Output
	pyFile, err := os.Open("../../testdata/python_minhkhoi_output.json")
	if err != nil {
		t.Fatalf("Failed to open python output: %v", err)
	}
	defer pyFile.Close()

	var pyOut PythonOutput
	if err := json.NewDecoder(pyFile).Decode(&pyOut); err != nil {
		t.Fatalf("Failed to decode python output: %v", err)
	}

	// 5. Compare Length
	if len(goResults) != len(pyOut.MissingVaccines) {
		t.Errorf("Length mismatch: Go %d, Py %d", len(goResults), len(pyOut.MissingVaccines))
	}

	// Compare Items
	for i, pyItem := range pyOut.MissingVaccines {
		if i >= len(goResults) {
			break
		}
		goItem := goResults[i]

		if goItem.VaccineName != pyItem.VaccineName {
			t.Errorf("Item %d Name mismatch: Go %s, Py %s", i, goItem.VaccineName, pyItem.VaccineName)
		}

		var goNextDoseStr *string
		if goItem.EarliestNextDoseDate != nil {
			s := goItem.EarliestNextDoseDate.Format("02/01/2006")
			goNextDoseStr = &s
		}

		if (goNextDoseStr == nil && pyItem.NextDose != nil) || (goNextDoseStr != nil && pyItem.NextDose == nil) {
			t.Errorf("Item %d NextDose mismatch: Go %v, Py %v", i, goNextDoseStr, pyItem.NextDose)
		} else if goNextDoseStr != nil && pyItem.NextDose != nil && *goNextDoseStr != *pyItem.NextDose {
			t.Errorf("Item %d NextDose mismatch: Go %s, Py %s", i, *goNextDoseStr, *pyItem.NextDose)
		}

		if goItem.Description != pyItem.Message {
			t.Errorf("Item %d Message mismatch:\nGo: %s\nPy: %s", i, goItem.Description, pyItem.Message)
		}
	}
}
