package engine

import (
	"time"
	"tracuutiemchung-engine/engine/checkers"
	"tracuutiemchung-engine/engine/utils"
	"tracuutiemchung-engine/models"
)

// ProcessAllRules runs all vaccine rules against the patient record.
func ProcessAllRules(rules map[string]models.VaccineRule, record models.PatientRecord, analysisDate time.Time) []models.MissingItem {
	var allResults []models.MissingItem
	visitedRules := make(map[string]bool)

	// Normalize AdministeredMap keys
	normalizedAdminMap := make(map[string][]models.AdministeredDose)
	for k, v := range record.AdministeredMap {
		normKey := utils.NormalizeVaccineName(k)
		normalizedAdminMap[normKey] = v
	}
	record.AdministeredMap = normalizedAdminMap

	// 1. Handle Special Groups first (like Pneumo)
	pneumoRules := []string{"Prevenar13", "Synflorix", "Vaxneuvance", "Pneumovax23"}
	hasPneumo := false
	for _, name := range pneumoRules {
		if _, ok := rules[name]; ok {
			hasPneumo = true
			break
		}
	}

	if hasPneumo {
		// Pneumo rules are handled together
		rulesPtr := make(map[string]*models.VaccineRule)
		for k := range rules {
			r := rules[k]
			rulesPtr[k] = &r
		}
		
		pneumoResults := checkers.CheckSpecialPneumo(record.AdministeredMap, record.BirthDate, analysisDate, rulesPtr)
		if len(pneumoResults) > 0 {
			allResults = append(allResults, pneumoResults...)
		}
		
		for _, name := range pneumoRules {
			visitedRules[name] = true
		}
	}

	// 2. Ordered keys from vaccine_rules.json for parity
	orderedKeys := []string{
		"MMR_Group",
		"Varivax",
		"MeningococcalACYW_Group",
		"JE_Group",
		"VA-MENGOC-BC",
		"MVVAC",
		"TyphimVi",
		"Morcvax",
		"Six_In_One_Combined",
		"Rota",
		"Prevenar13",
		"Vaxneuvance",
		"Synflorix",
		"Pneumovax23",
		"HepA",
		"Flu",
		"BCG",
	}

	for _, id := range orderedKeys {
		rule, ok := rules[id]
		if !ok || visitedRules[id] {
			continue
		}

		var results []models.MissingItem

		// Prepare full map for checkers
		fullAdminMap := record.AdministeredMap
		
		results = nil
		switch rule.Type {
		case "single_series", "single_dose_min_age", "single_series_min_age":
			results = checkers.CheckBasicSeries(&rule, fullAdminMap, record.BirthDate, analysisDate)
		case "age_dependent_series":
			results = checkers.CheckAgeDependentSeries(&rule, fullAdminMap, record.BirthDate, analysisDate)
		case "mmr_equivalent_group":
			results = checkers.CheckCumulativeGroup(&rule, fullAdminMap, record.BirthDate, analysisDate)
		case "meningococcal_acyw_group":
			results = checkers.CheckMeningococcalGroup(&rule, fullAdminMap, record.BirthDate, analysisDate)
		case "group_alternative_courses_age_range", "group_alternative_courses_min_age":
			results = checkers.CheckAlternativeCoursesGroup(&rule, fullAdminMap, record.BirthDate, analysisDate)
		case "flu_group":
			results = checkers.CheckFluGroup(&rule, fullAdminMap, record.BirthDate, analysisDate)
		}

		if len(results) > 0 {
			allResults = append(allResults, results...)
		}
	}

	// 3. Apply Complex Interactions (e.g., MMR vs Measles)
	allResults = checkers.ApplyComplexInteractions(allResults, record.AdministeredMap)

	return allResults
}

