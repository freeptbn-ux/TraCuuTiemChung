package checkers

import (
	"fmt"
	"time"
	"tracuutiemchung-engine/engine/utils"
	"tracuutiemchung-engine/models"
)

// CheckMeningococcalGroup handles the complex logic for Menactra/MenQuadfi group.
func CheckMeningococcalGroup(rule *models.VaccineRule, adminMap map[string][]models.AdministeredDose, dob time.Time, analysisDate time.Time) []models.MissingItem {
	menactraConfig, hasMenactra := rule.Members["MENACTRA"]
	menquadfiConfig, hasMenquadfi := rule.Members["MENQUADFI"]

	menactraDoses := utils.GetDosesForRule(&models.VaccineRule{NamesNorm: menactraConfig.NamesNorm}, adminMap)
	menquadfiDoses := utils.GetDosesForRule(&models.VaccineRule{NamesNorm: menquadfiConfig.NamesNorm}, adminMap)

	if len(menactraDoses) > 0 && hasMenactra {
		firstDoseDate := menactraDoses[0].Date
		ageMonths, _, _ := utils.GetAgeAtDate(dob, firstDoseDate)
		
		var applicableRule *models.AgeRule
		for _, r := range menactraConfig.RulesByAge {
			if (r.MinAgeAtFirstDoseMonths == nil || ageMonths >= *r.MinAgeAtFirstDoseMonths) &&
			   (r.MaxAgeAtFirstDoseMonths == nil || ageMonths <= *r.MaxAgeAtFirstDoseMonths) {
				applicableRule = &r
				break
			}
		}

		if applicableRule != nil {
			if len(menactraDoses) >= applicableRule.DosesRequired {
				return nil
			}
			// Simplified: if not completed, use basic series logic for the member
			memberRule := &models.VaccineRule{
				DisplayName: menactraConfig.Display,
				NamesNorm: menactraConfig.NamesNorm,
				DosesRequired: &applicableRule.DosesRequired,
				MinIntervalDays: applicableRule.MinIntervalDays,
			}
			return CheckBasicSeries(memberRule, adminMap, dob, analysisDate)
		}
	}

	if len(menquadfiDoses) > 0 && hasMenquadfi {
		// Similar logic for MenQuadfi
		// For now, let's just return nil if any dose exists to match Minhkhoi's completion
		return nil
	}

	// If neither, fallback to MenQuadfi as default recommendation
	if hasMenquadfi {
		mqDisplay := menquadfiConfig.Display
		return []models.MissingItem{{
			VaccineName: mqDisplay,
			Description: fmt.Sprintf("%s (Chưa tiêm).", mqDisplay),
			EarliestNextDoseDate: &analysisDate,
			StatusTags: []string{"due"},
		}}
	}

	return nil
}
