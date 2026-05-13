package main

import (
	"encoding/json"
	"fmt"
	"time"

	"tracuutiemchung-engine/engine"
	"tracuutiemchung-engine/engine/checkers"
	"tracuutiemchung-engine/engine/utils"
	"tracuutiemchung-engine/models"
	"tracuutiemchung-engine/tests/parity"
)

func collectDoses(rule *models.VaccineRule, adminMap map[string][]models.AdministeredDose) []models.AdministeredDose {
	var doses []models.AdministeredDose
	for _, normName := range rule.NamesNorm {
		if d, ok := adminMap[normName]; ok {
			doses = append(doses, d...)
		}
	}
	return doses
}

func main() {
	rules, _ := engine.LoadVaccineRules("vercel-backend/assets/vaccine_rules.json")
	patient, records, _ := parity.ParseMinhkhoiHTML("test/Minhkhoi.html")
	dob, _ := time.Parse("02/01/2006", patient.Birth)
	analysisDate, _ := time.Parse("02/01/2006", patient.SystemDate)

	adminMap := make(map[string][]models.AdministeredDose)
	for _, rec := range records {
		norm := utils.NormalizeVaccineName(rec.RawName)
		adminMap[norm] = append(adminMap[norm], models.AdministeredDose{
			VaccineName: rec.RawName,
			Date:        rec.Date,
		})
	}
	
	rule := rules["Prevenar13"]
	tempMap := map[string][]models.AdministeredDose{
		rule.DisplayName: collectDoses(&rule, adminMap),
	}
	res := checkers.CheckAgeDependentSeries(&rule, tempMap, dob, analysisDate)
	b, _ := json.MarshalIndent(res, "", "  ")
	fmt.Println(string(b))
}
