package analyzer

import (
	"fmt"
	"time"
)

const GracePeriodDays = 0

// getAgeStatusAndEarliestDate xác định trạng thái tuổi và ngày tiêm sớm nhất có thể.
func (e *Engine) getAgeStatusAndEarliestDate(rule VaccineRule, forGroupDisplayName string) (string, *time.Time, []string) {
	displayNamePrefix := ""
	if forGroupDisplayName != "" {
		displayNamePrefix = forGroupDisplayName + " - "
	}

	if e.DOB.IsZero() {
		return displayNamePrefix + "Không có ngày sinh để kiểm tra tuổi", nil, []string{"error_dob"}
	}

	if e.AnalysisDate.Before(e.DOB) {
		return displayNamePrefix + "Ngày phân tích trước ngày sinh", nil, []string{"error_date"}
	}

	_, ageDays, _ := GetAgeAtDate(e.DOB, e.AnalysisDate)

	minAgeMonths := rule.MinAgeMonthsOverallGroup
	if minAgeMonths == 0 {
		minAgeMonths = rule.MinAgeMonthsAtFirstDose
	}
	if minAgeMonths == 0 {
		minAgeMonths = rule.MinAgeMonthsOverall
	}

	minAgeWeeks := rule.MinAgeWeeksOverallGroup
	if minAgeWeeks == 0 {
		minAgeWeeks = rule.MinAgeWeeksAtFirstDose
	}
	if minAgeWeeks == 0 {
		minAgeWeeks = rule.MinAgeWeeksOverall
	}

	minAgeYears := rule.MinAgeYearsOverallGroup
	if minAgeYears == 0 {
		minAgeYears = rule.MinAgeYearsAtFirstDose
	}
	if minAgeYears == 0 {
		minAgeYears = rule.MinAgeYearsOverall
	}

	minAgeDaysVal := rule.MinAgeDaysOverallGroup
	if minAgeDaysVal == 0 {
		minAgeDaysVal = rule.MinAgeDaysAtFirstDose
	}
	if minAgeDaysVal == 0 {
		minAgeDaysVal = rule.MinAgeDaysOverall
	}

	var earliestAcceptableDate time.Time
	statusMessage := ""
	statusTags := []string{"eligible"}

	// Priority based on Python logic: Months > Years > Weeks > Days
	if minAgeMonths > 0 {
		targetMinAgeDate := AddMonths(e.DOB, minAgeMonths)
		earliestAcceptableDate = targetMinAgeDate.AddDate(0, 0, -GracePeriodDays)
		if e.AnalysisDate.Before(earliestAcceptableDate) {
			statusMessage = fmt.Sprintf("cần %d tháng tuổi", minAgeMonths)
			statusTags = []string{"too_young"}
		}
	} else if minAgeYears > 0 {
		targetMinAgeDate := AddYears(e.DOB, minAgeYears)
		earliestAcceptableDate = targetMinAgeDate.AddDate(0, 0, -GracePeriodDays)
		if e.AnalysisDate.Before(earliestAcceptableDate) {
			statusMessage = fmt.Sprintf("cần %d tuổi", minAgeYears)
			statusTags = []string{"too_young"}
		}
	} else if minAgeWeeks > 0 {
		requiredTotalDaysForWeeks := minAgeWeeks * 7
		effectiveMinTotalDays := requiredTotalDaysForWeeks - GracePeriodDays
		if ageDays < effectiveMinTotalDays {
			earliestAcceptableDate = e.DOB.AddDate(0, 0, effectiveMinTotalDays)
			statusMessage = fmt.Sprintf("cần %d tuần tuổi", minAgeWeeks)
			statusTags = []string{"too_young"}
		}
	} else if minAgeDaysVal > 0 {
		effectiveMinTotalDays := minAgeDaysVal - GracePeriodDays
		if ageDays < effectiveMinTotalDays {
			earliestAcceptableDate = e.DOB.AddDate(0, 0, effectiveMinTotalDays)
			displayReq := ""
			if minAgeDaysVal >= 60 {
				displayReq = fmt.Sprintf("%d tháng", minAgeDaysVal/30)
			} else {
				displayReq = fmt.Sprintf("%d ngày", minAgeDaysVal)
			}
			statusMessage = fmt.Sprintf("cần >%s tuổi", displayReq)
			statusTags = []string{"too_young"}
		}
	}

	if containsStr(statusTags, "too_young") {
		return statusMessage, &earliestAcceptableDate, statusTags
	}
	
	// Eligible case
	var retDate *time.Time
	if !e.DOB.IsZero() {
		retDate = &e.AnalysisDate
	}
	return "đủ điều kiện tuổi", retDate, statusTags
}

// checkFirstDoseAgeValidity kiểm tra tuổi mũi 1.
func (e *Engine) checkFirstDoseAgeValidity(firstDoseDate time.Time, rule VaccineRule, ruleDisplayName string) (bool, *AnalysisResult) {
	if e.DOB.IsZero() {
		return true, nil
	}

	if firstDoseDate.Before(e.DOB) {
		return false, &AnalysisResult{
			VaccineNameForPopup: ruleDisplayName,
			Description:         fmt.Sprintf("%s - Lỗi tính tuổi cho mũi đầu (ngày tiêm có thể trước ngày sinh).", ruleDisplayName),
			EarliestNextDoseDate: nil,
			StatusTags:           []string{"error_age_calculation"},
		}
	}

	_, ageDays, _ := GetAgeAtDate(e.DOB, firstDoseDate)

	minAgeMonths := rule.MinAgeMonthsOverallGroup
	if minAgeMonths == 0 {
		minAgeMonths = rule.MinAgeMonthsAtFirstDose
	}
	if minAgeMonths == 0 {
		minAgeMonths = rule.MinAgeMonthsOverall
	}

	minAgeWeeks := rule.MinAgeWeeksOverallGroup
	if minAgeWeeks == 0 {
		minAgeWeeks = rule.MinAgeWeeksAtFirstDose
	}
	if minAgeWeeks == 0 {
		minAgeWeeks = rule.MinAgeWeeksOverall
	}

	minAgeYears := rule.MinAgeYearsOverallGroup
	if minAgeYears == 0 {
		minAgeYears = rule.MinAgeYearsAtFirstDose
	}
	if minAgeYears == 0 {
		minAgeYears = rule.MinAgeYearsOverall
	}

	minAgeDaysVal := rule.MinAgeDaysOverallGroup
	if minAgeDaysVal == 0 {
		minAgeDaysVal = rule.MinAgeDaysAtFirstDose
	}
	if minAgeDaysVal == 0 {
		minAgeDaysVal = rule.MinAgeDaysOverall
	}

	errorDetail := ""
	// Priority based on Python logic for first dose: Days > Weeks > Months > Years
	if minAgeDaysVal > 0 {
		effectiveMinDays := minAgeDaysVal - GracePeriodDays
		if ageDays < effectiveMinDays {
			displayAge := ""
			if minAgeDaysVal >= 60 {
				displayAge = fmt.Sprintf("%d tháng", minAgeDaysVal/30)
			} else {
				displayAge = fmt.Sprintf("%d ngày", minAgeDaysVal)
			}
			errorDetail = fmt.Sprintf("Mũi 1 tiêm quá sớm (cần >%s, thực tế %d ngày tuổi).", displayAge, ageDays)
		}
	} else if minAgeWeeks > 0 {
		effectiveMinDays := (minAgeWeeks * 7) - GracePeriodDays
		if ageDays < effectiveMinDays {
			errorDetail = fmt.Sprintf("Mũi 1 tiêm quá sớm (cần %d tuần, thực tế %d ngày tuổi).", minAgeWeeks, ageDays)
		}
	} else if minAgeMonths > 0 {
		earliestAllowed := AddMonths(e.DOB, minAgeMonths).AddDate(0, 0, -GracePeriodDays)
		if firstDoseDate.Before(earliestAllowed) {
			errorDetail = fmt.Sprintf("Mũi 1 tiêm quá sớm (cần %d tháng tuổi).", minAgeMonths)
		}
	} else if minAgeYears > 0 {
		earliestAllowed := AddYears(e.DOB, minAgeYears).AddDate(0, 0, -GracePeriodDays)
		if firstDoseDate.Before(earliestAllowed) {
			errorDetail = fmt.Sprintf("Mũi 1 tiêm quá sớm (cần %d tuổi).", minAgeYears)
		}
	}

	if errorDetail != "" {
		return false, &AnalysisResult{
			VaccineNameForPopup: ruleDisplayName,
			Description:         fmt.Sprintf("%s - %s", ruleDisplayName, errorDetail),
			EarliestNextDoseDate: nil,
			StatusTags:           []string{"error_age_first_dose", "too_early"},
			IsMissing:            true,
		}
	}

	return true, nil
}

func containsStr(slice []string, val string) bool {
	for _, item := range slice {
		if item == val {
			return true
		}
	}
	return false
}
