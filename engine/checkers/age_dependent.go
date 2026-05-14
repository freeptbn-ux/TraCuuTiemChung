package checkers

import (
	"time"
	"tracuutiemchung-engine/engine/utils"
	"tracuutiemchung-engine/models"
)

// CheckAgeDependentSeries selects a sub-rule based on the age at the first dose.
func CheckAgeDependentSeries(rule *models.VaccineRule, adminMap map[string][]models.AdministeredDose, dob time.Time, analysisDate time.Time) []models.MissingItem {
	// Collect doses using the centralized helper
	doses := utils.GetDosesForRule(rule, adminMap)

	var ageAtFirstDoseMonths int
	var ageAtFirstDoseWeeks int
	if len(doses) > 0 {
		m, w, _ := utils.GetAgeAtDate(dob, doses[0].Date)
		ageAtFirstDoseMonths = m
		ageAtFirstDoseWeeks = w
	} else {
		// If no dose yet, use age at analysis date to predict which regimen will apply if started today
		m, w, _ := utils.GetAgeAtDate(dob, analysisDate)
		ageAtFirstDoseMonths = m
		ageAtFirstDoseWeeks = w
	}

	// 1. Try RulesByAge (used by Meningococcal, etc.)
	for _, ageRule := range rule.RulesByAge {
		if isMatchAgeRule(&ageRule, ageAtFirstDoseMonths, ageAtFirstDoseWeeks) {
			return CheckSeriesInternal(
				rule.DisplayName,
				ageRule.DosesRequired,
				ageRule.MinIntervalDays,
				ageRule.MinAgeAtFirstDoseMonths,
				ageRule.MinAgeWeeksAtFirstDose,
				nil, // Days not specified in AgeRule usually
				nil, // Years not specified in AgeRule
				ageRule.Booster,
				nil, // BoosterIntervalYears
				doses,
				dob,
				analysisDate,
			)
		}
	}

	// 2. Try Regimens (used by MMR, etc.)
	for _, regimen := range rule.Regimens {
		if isMatchRegimen(&regimen, ageAtFirstDoseMonths) {
			return CheckSeriesInternal(
				rule.DisplayName,
				regimen.DosesRequired,
				regimen.MinIntervalDays,
				regimen.MinAgeAtFirstDoseMonths,
				nil, // Weeks
				nil, // Days
				nil, // Years
				nil, // Booster Rule
				nil, // BoosterIntervalYears
				doses,
				dob,
				analysisDate,
			)
		}
	}

	// Fallback to top-level basic series if no age-dependent rule matches
	return CheckBasicSeries(rule, adminMap, dob, analysisDate)
}

func isMatchAgeRule(rule *models.AgeRule, months, weeks int) bool {
	if rule.MinAgeAtFirstDoseMonths != nil && months < *rule.MinAgeAtFirstDoseMonths {
		return false
	}
	if rule.MaxAgeAtFirstDoseMonths != nil && months > *rule.MaxAgeAtFirstDoseMonths {
		return false
	}
	if rule.MinAgeWeeksAtFirstDose != nil && weeks < *rule.MinAgeWeeksAtFirstDose {
		return false
	}
	return true
}

func isMatchRegimen(regimen *models.Regimen, months int) bool {
	if regimen.MinAgeAtFirstDoseMonths != nil && months < *regimen.MinAgeAtFirstDoseMonths {
		return false
	}
	if regimen.MaxAgeAtFirstDoseMonths != nil && months > *regimen.MaxAgeAtFirstDoseMonths {
		return false
	}
	return true
}
