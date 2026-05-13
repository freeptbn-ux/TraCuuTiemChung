package checkers

import (
	"fmt"
	"sort"
	"time"
	"tracuutiemchung-engine/engine/utils"
	"tracuutiemchung-engine/models"
)

// CheckBasicSeries checks the basic vaccination series and recommends the next dose.
func CheckBasicSeries(rule *models.VaccineRule, adminMap map[string][]models.AdministeredDose, dob time.Time, analysisDate time.Time) []models.MissingItem {
	// Collect doses using the centralized helper
	allDoses := utils.GetDosesForRule(rule, adminMap)

	dosesRequired := 0
	if rule.DosesRequired != nil {
		dosesRequired = *rule.DosesRequired
	}

	return CheckSeriesInternal(
		rule.DisplayName,
		dosesRequired,
		rule.MinIntervalDays,
		rule.MinAgeMonthsAtFirstDose,
		rule.MinAgeWeeksAtFirstDose,
		rule.MinAgeDaysAtFirstDose,
		rule.MinAgeYearsAtFirstDose,
		nil, // boosterRule
		nil, // boosterIntervalYears
		allDoses,
		dob,
		analysisDate,
	)
}

// CheckSeriesInternal is a shared implementation for basic series logic.
func CheckSeriesInternal(
	displayName string,
	dosesRequired int,
	minIntervalDays []*int,
	minAgeMonthsAtFirstDose *int,
	minAgeWeeksAtFirstDose *int,
	minAgeDaysAtFirstDose *int,
	minAgeYearsAtFirstDose *int,
	boosterRule *models.BoosterRule,
	boosterIntervalYears *int,
	doses []models.AdministeredDose,
	dob time.Time,
	analysisDate time.Time,
) []models.MissingItem {
	// 0. Sort doses by date just in case
	sort.Slice(doses, func(i, j int) bool {
		return doses[i].Date.Before(doses[j].Date)
	})

	// 1. Count valid doses (Mirror Python: count all doses if first dose age is valid)
	validDoses := 0
	var lastValidDoseDate time.Time

	if len(doses) > 0 {
		// Check first dose age constraints
		dose := doses[0]
		months, weeks, _ := utils.GetAgeAtDate(dob, dose.Date)
		isValid := true
		if minAgeMonthsAtFirstDose != nil && months < *minAgeMonthsAtFirstDose {
			isValid = false
		}
		if minAgeWeeksAtFirstDose != nil && weeks < *minAgeWeeksAtFirstDose {
			isValid = false
		}
		if minAgeDaysAtFirstDose != nil {
			days := int(dose.Date.Sub(dob).Hours() / 24)
			if days < *minAgeDaysAtFirstDose {
				isValid = false
			}
		}
		if minAgeYearsAtFirstDose != nil {
			dateByAge := utils.AddYears(dob, *minAgeYearsAtFirstDose)
			if dose.Date.Before(dateByAge) {
				isValid = false
			}
		}

		if isValid {
			validDoses = len(doses)
			lastValidDoseDate = doses[len(doses)-1].Date
		}
	}

	// 2. Check if series complete
	if validDoses >= dosesRequired && dosesRequired > 0 {
		// Handle Booster if requested (Phase 05+)
		if boosterIntervalYears != nil && !lastValidDoseDate.IsZero() {
			earliestBoosterDate := utils.AddYears(lastValidDoseDate, *boosterIntervalYears)
			statusTags := []string{"booster"}
			if analysisDate.After(earliestBoosterDate) || analysisDate.Equal(earliestBoosterDate) {
				statusTags = append(statusTags, "due")
			} else {
				statusTags = append(statusTags, "future")
			}

			return []models.MissingItem{{
				VaccineName:          displayName,
				Description:          fmt.Sprintf("%s - Cần tiêm nhắc lại", displayName),
				EarliestNextDoseDate: &earliestBoosterDate,
				StatusTags:           statusTags,
			}}
		}
		return nil
	}

	// 3. Calculate next dose date
	nextDoseIdx := validDoses

	var earliestNextDoseDate time.Time
	
	if validDoses == 0 {
		earliestNextDoseDate = dob
	} else {
		earliestNextDoseDate = lastValidDoseDate
	}

	// Rule 1: Interval Constraint
	if nextDoseIdx < len(minIntervalDays) && minIntervalDays[nextDoseIdx] != nil {
		intervalDays := *minIntervalDays[nextDoseIdx]
		if !lastValidDoseDate.IsZero() {
			dateByInterval := lastValidDoseDate.AddDate(0, 0, intervalDays)
			if dateByInterval.After(earliestNextDoseDate) {
				earliestNextDoseDate = dateByInterval
			}
		}
	}

	// Rule 2: Age Constraint
	if nextDoseIdx == 0 {
		if minAgeMonthsAtFirstDose != nil {
			dateByAge := utils.AddMonths(dob, *minAgeMonthsAtFirstDose)
			if dateByAge.After(earliestNextDoseDate) {
				earliestNextDoseDate = dateByAge
			}
		}
		if minAgeWeeksAtFirstDose != nil {
			dateByAge := dob.AddDate(0, 0, *minAgeWeeksAtFirstDose * 7)
			if dateByAge.After(earliestNextDoseDate) {
				earliestNextDoseDate = dateByAge
			}
		}
		if minAgeDaysAtFirstDose != nil {
			dateByAge := dob.AddDate(0, 0, *minAgeDaysAtFirstDose)
			if dateByAge.After(earliestNextDoseDate) {
				earliestNextDoseDate = dateByAge
			}
		}
		if minAgeYearsAtFirstDose != nil {
			dateByAge := utils.AddYears(dob, *minAgeYearsAtFirstDose)
			if dateByAge.After(earliestNextDoseDate) {
				earliestNextDoseDate = dateByAge
			}
		}
	}

	// Floor to analysisDate for recommendations (Python parity)
	if earliestNextDoseDate.Before(analysisDate) {
		earliestNextDoseDate = analysisDate
	}

	remainingDoses := dosesRequired - validDoses
	statusTags := []string{"due"}
	
	description := ""
	if validDoses == 0 {
		description = fmt.Sprintf("%s (Chưa tiêm - cần %d liều). đủ điều kiện tuổi", displayName, dosesRequired)
		statusTags = append(statusTags, "eligible")
	} else {
		// Matching Python: {VaccineName} - Mũi {X} (Cần thêm {Y} liều)
		description = fmt.Sprintf("%s - Mũi %d (Cần thêm %d liều)", displayName, validDoses+1, remainingDoses)
	}

	return []models.MissingItem{{
		VaccineName:          displayName,
		Description:          description,
		EarliestNextDoseDate: &earliestNextDoseDate,
		StatusTags:           statusTags,
	}}
}



