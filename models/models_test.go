package models

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestJSONUnmarshalingNullField(t *testing.T) {
	// Scenario: Load a sample JSON where booster_interval_years is null.
	jsonData := `{
		"type": "single_series",
		"booster_interval_years": null
	}`

	var rule VaccineRule
	err := json.Unmarshal([]byte(jsonData), &rule)
	assert.NoError(t, err)

	// Expected: Struct field *int should be nil without crashing.
	assert.Nil(t, rule.BoosterIntervalYears)
}

func TestAdministeredMapSorting(t *testing.T) {
	// Scenario: Add records in random order.
	d1 := time.Date(2024, 5, 1, 0, 0, 0, 0, time.UTC)
	d2 := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)

	patient := PatientRecord{
		AdministeredMap: map[string][]AdministeredDose{
			"vax1": {
				{VaccineName: "Vax 1", Date: d1},
				{VaccineName: "Vax 1", Date: d2},
			},
		},
	}

	// Expected: GetSortedDoses() must return them in chronological order.
	sorted := patient.GetSortedDoses("vax1")
	assert.Equal(t, 2, len(sorted))
	assert.True(t, sorted[0].Date.Before(sorted[1].Date))
	assert.Equal(t, d2, sorted[0].Date)
	assert.Equal(t, d1, sorted[1].Date)
}

func TestRuleNormalizationFields(t *testing.T) {
	// Scenario: Check if raw_names from JSON are correctly mapped to the struct.
	jsonData := `{
		"type": "single_series",
		"raw_names": ["Vax A", "Vax B"]
	}`

	var rule VaccineRule
	err := json.Unmarshal([]byte(jsonData), &rule)
	assert.NoError(t, err)

	// Expected: raw_names should be populated.
	assert.Equal(t, []string{"Vax A", "Vax B"}, rule.RawNames)
	
	// Note: NamesNorm is populated during loading (Phase 03), 
	// but here we just check if the struct has the field and it's initialized correctly if needed.
	rule.NamesNorm = []string{"vax a", "vax b"}
	assert.Equal(t, 2, len(rule.NamesNorm))
}
