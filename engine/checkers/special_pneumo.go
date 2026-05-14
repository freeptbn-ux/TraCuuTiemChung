package checkers

import (
	"time"
	"tracuutiemchung-engine/engine/utils"
	"tracuutiemchung-engine/models"
)

// CheckSpecialPneumo handles the complex logic for pneumococcal vaccines (Synflorix, Prevenar 13, Vaxneuvance)
// and detects interchange between them.
func CheckSpecialPneumo(adminMap map[string][]models.AdministeredDose, dob time.Time, analysisDate time.Time, rules map[string]*models.VaccineRule) []models.MissingItem {
	prevenar13Key := "Prevenar13"
	synflorixKey := "Synflorix"
	vaxneuvanceKey := "Vaxneuvance"
	pneumovax23Key := "Pneumovax23"

	// Collect doses and counts for each rule
	pneumoRecords := make(map[string][]models.AdministeredDose)
	numDoses := make(map[string]int)
	var allDoses []models.AdministeredDose
	
	activeSeriesKeys := []string{}
	for _, k := range []string{prevenar13Key, synflorixKey, vaxneuvanceKey} {
		rule, ok := rules[k]
		if !ok {
			continue
		}
		var doses []models.AdministeredDose
		for _, norm := range rule.NamesNorm {
			if d, ok := adminMap[norm]; ok {
				doses = append(doses, d...)
			}
		}
		// Also check the rule ID itself as a key (normalized)
		normID := utils.NormalizeVaccineName(k)
		if d, ok := adminMap[normID]; ok {
			// Avoid double-counting if rule ID matches a normalized name
			alreadyAdded := false
			for _, norm := range rule.NamesNorm {
				if norm == normID {
					alreadyAdded = true
					break
				}
			}
			if !alreadyAdded {
				doses = append(doses, d...)
			}
		}
		pneumoRecords[k] = doses
		numDoses[k] = len(doses)
		if len(doses) > 0 {
			allDoses = append(allDoses, doses...)
			activeSeriesKeys = append(activeSeriesKeys, k)
		}
	}

	// 1. If Pneumovax 23 is present -> skip all (standard logic)
	pv23Rule, ok := rules[pneumovax23Key]
	if ok {
		for _, norm := range pv23Rule.NamesNorm {
			if d, ok := adminMap[norm]; ok && len(d) > 0 {
				return nil
			}
		}
	}

	// 2. Detect Interchange
	isInterchange := len(activeSeriesKeys) > 1
	var interchangeWarning string
	if isInterchange {
		mixedNames := []string{}
		for _, k := range activeSeriesKeys {
			mixedNames = append(mixedNames, rules[k].DisplayName)
		}
		mixedStr := ""
		for i, name := range mixedNames {
			if i > 0 {
				mixedStr += " và "
			}
			mixedStr += name
		}
		interchangeWarning = "Cảnh báo: Đã ghi nhận tiêm xen kẽ các loại phế cầu (" + mixedStr + "). Không nên sử dụng xen kẽ."
	}

	// 3. If any brand has >= 4 doses -> complete
	if len(allDoses) >= 4 {
		return nil
	}

	// 4. Determine which rule to use for next dose
	primaryKey := prevenar13Key
	if len(activeSeriesKeys) > 0 {
		// Use the most recent brand
		utils.SortDoses(allDoses)
		lastDose := allDoses[len(allDoses)-1]
		// Find which key this last dose belongs to
		for _, k := range activeSeriesKeys {
			rule := rules[k]
			for _, norm := range rule.NamesNorm {
				if lastDose.VaccineName == norm || lastDose.VaccineName == rule.DisplayName {
					primaryKey = k
					break
				}
			}
		}
	}

	// 5. Special logic for 2+ years old
	_, _, ageYears := utils.GetAgeAtDate(dob, analysisDate)
	if ageYears >= 2 && len(allDoses) < 3 {
		if pv23Rule != nil {
			return []models.MissingItem{{
				VaccineName: pv23Rule.DisplayName,
				Description: pv23Rule.DisplayName + ": Có thể tiêm 1 mũi để hoàn thành phác đồ phế cầu (do đã trên 2 tuổi và đã tiêm < 3 mũi phế cầu cộng dồn).",
				EarliestNextDoseDate: &analysisDate,
				StatusTags: []string{"info", "alternative_completion"},
			}}
		}
	}

	// 6. Use AgeDependentSeries logic with merged doses
	tempAdminMap := map[string][]models.AdministeredDose{
		rules[primaryKey].DisplayName: allDoses,
	}
	results := CheckAgeDependentSeries(rules[primaryKey], tempAdminMap, dob, analysisDate)

	// 7. Inject interchange warning if needed
	if isInterchange && len(results) > 0 {
		for i := range results {
			results[i].Description = interchangeWarning + " " + results[i].Description
			results[i].StatusTags = append(results[i].StatusTags, "error_interchange")
		}
	}

	return results
}

