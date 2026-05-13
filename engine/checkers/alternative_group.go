package checkers

import (
	"fmt"
	"strings"
	"time"
	"tracuutiemchung-engine/engine/utils"
	"tracuutiemchung-engine/models"
)

// CheckAlternativeCoursesGroup handles vaccine groups where multiple alternative courses exist (e.g., Rota, Hep A, JE).
func CheckAlternativeCoursesGroup(rule *models.VaccineRule, adminMap map[string][]models.AdministeredDose, dob time.Time, analysisDate time.Time) []models.MissingItem {
	// 1. Collect all doses and identify which brands were used
	var allDoses []models.AdministeredDose
	brandsUsed := make(map[int]bool) // index of course in rule.Courses
	
	for i, course := range rule.Courses {
		for _, normName := range course.NamesNorm {
			if doses, ok := adminMap[normName]; ok {
				allDoses = append(allDoses, doses...)
				brandsUsed[i] = true
			}
		}
	}

	utils.SortDoses(allDoses)

	displayName := rule.GroupDisplayName
	if displayName == "" {
		displayName = rule.DisplayName
	}

	// 2. If no doses, pick the first applicable one
	if len(allDoses) == 0 {
		if len(rule.Courses) == 0 {
			return nil
		}
		
		ageMonths, _, _ := utils.GetAgeAtDate(dob, analysisDate)
		if rule.MaxAgeMonthsToStartFirstDoseGroup != nil && ageMonths > *rule.MaxAgeMonthsToStartFirstDoseGroup {
			return []models.MissingItem{{
				VaccineName: displayName,
				Description: displayName + ": Đã qua " + fmt.Sprintf("%d", *rule.MaxAgeMonthsToStartFirstDoseGroup) + " tháng tuổi, không còn chỉ định bắt đầu.",
				EarliestNextDoseDate: nil,
				StatusTags: []string{"info", "too_old_to_start"},
			}}
		}

		courseDisplays := []string{}
		for _, c := range rule.Courses {
			courseDisplays = append(courseDisplays, c.Display)
		}
		alternativesStr := ""
		for i, d := range courseDisplays {
			if i > 0 {
				alternativesStr += ", "
			}
			alternativesStr += d
		}

		return []models.MissingItem{{
			VaccineName:          displayName,
			Description:          displayName + ": Chưa tiêm. Lựa chọn có thể tiêm: " + alternativesStr + ".",
			EarliestNextDoseDate: &analysisDate,
			StatusTags:           []string{"due"},
		}}
	}

	// 3. If doses exist, determine which course to follow
	selectedCourseIdx := -1
	maxDosesFound := 0
	anyCourseCompleted := false
	
	for i, course := range rule.Courses {
		courseDoses := 0
		for _, normName := range course.NamesNorm {
			if ds, ok := adminMap[normName]; ok {
				courseDoses += len(ds)
			}
		}
		if courseDoses >= course.DosesRequired && course.DosesRequired > 0 {
			anyCourseCompleted = true
		}
		if courseDoses > maxDosesFound {
			maxDosesFound = courseDoses
			selectedCourseIdx = i
		}
	}
	
	if selectedCourseIdx == -1 {
		selectedCourseIdx = 0 // Fallback
	}
	
	selectedCourse := rule.Courses[selectedCourseIdx]

	// 4. Mirror Python's Rota bug: return 'too old' if over age limit AND (no group-level names_norm)
	// Even if a course was completed, if the group-level names_norm is empty, Python adds the message.
	ageMonths, _, _ := utils.GetAgeAtDate(dob, analysisDate)
	if rule.Type == "group_alternative_courses_min_age" && rule.MaxAgeMonthsToStartFirstDoseGroup != nil {
		if ageMonths > *rule.MaxAgeMonthsToStartFirstDoseGroup {
			// In Python, it checks if any_course_completed was set. 
			// But for Rota, because rule.NamesNorm is empty, it often misses completion.
			if !anyCourseCompleted || len(rule.NamesNorm) == 0 {
				return []models.MissingItem{{
					VaccineName:          displayName,
					Description:          displayName + ": Đã qua " + fmt.Sprintf("%d", *rule.MaxAgeMonthsToStartFirstDoseGroup) + " tháng tuổi, không còn chỉ định bắt đầu.",
					EarliestNextDoseDate: nil,
					StatusTags:           []string{"too_old_to_start"},
				}}
			}
		}
	}

	if anyCourseCompleted {
		return nil
	}

	return CheckSeriesInternal(
		selectedCourse.Display,
		selectedCourse.DosesRequired,
		selectedCourse.MinIntervalDays,
		selectedCourse.MinAgeMonthsAtFirstDose,
		nil, nil, nil, nil,
		selectedCourse.BoosterIntervalYears,
		allDoses,
		dob,
		analysisDate,
	)
}

// CheckFluGroup handles the special logic for seasonal influenza vaccines.
func CheckFluGroup(rule *models.VaccineRule, adminMap map[string][]models.AdministeredDose, dob time.Time, analysisDate time.Time) []models.MissingItem {
	// Find all flu doses across possible brands/recognition keywords
	var allFluDoses []models.AdministeredDose
	seenDates := make(map[int64]bool)

	// Collect doses matching names_norm or recognition keywords (substring match)
	for adminNormName, doses := range adminMap {
		match := false
		for _, name := range rule.NamesNorm {
			if adminNormName == name {
				match = true
				break
			}
		}
		if !match {
			for _, kw := range rule.RecognitionKeywords {
				if strings.Contains(adminNormName, utils.NormalizeVaccineName(kw)) {
					match = true
					break
				}
			}
		}

		if match {
			for _, d := range doses {
				if !seenDates[d.Date.Unix()] {
					allFluDoses = append(allFluDoses, d)
					seenDates[d.Date.Unix()] = true
				}
			}
		}
	}

	utils.SortDoses(allFluDoses)

	_, _, ageYears := utils.GetAgeAtDate(dob, analysisDate)
	
	if len(allFluDoses) == 0 {
		// If never vaccinated
		earliest := utils.AddMonths(dob, 6)
		if analysisDate.Before(earliest) {
			return nil
		}
		return []models.MissingItem{{
			VaccineName: "Vắc xin Cúm",
			Description: "Vắc xin Cúm (Chưa tiêm. Lần đầu (nếu <9 tuổi) có thể cần 2 mũi cách nhau ~1 tháng, sau đó nhắc lại hàng năm). đủ điều kiện tuổi",
			EarliestNextDoseDate: &analysisDate,
			StatusTags: []string{"eligible"},
		}}
	}

	// Simplified logic for Phase 06: 
	// If child < 9 years and has only 1 dose, recommend dose 2 after 30 days.
	// Otherwise, recommend annual booster if last dose > 1 year ago.
	
	lastDoseDate := allFluDoses[len(allFluDoses)-1].Date
	
	if ageYears < 9 && len(allFluDoses) == 1 {
		interval := 30
		if rule.InitialSeriesIntervalDays != nil {
			interval = *rule.InitialSeriesIntervalDays
		}
		nextDate := lastDoseDate.AddDate(0, 0, interval)
		return []models.MissingItem{{
			VaccineName: "Cúm",
			Description: "Cúm - Cần tiêm Mũi 2 (Cách mũi 1 ít nhất 1 tháng)",
			EarliestNextDoseDate: &nextDate,
			StatusTags: []string{"due"},
		}}
	}

	// Annual booster logic
	nextAnniversary := lastDoseDate.AddDate(1, 0, 0)
	displayName := rule.GroupDisplayName
	if displayName == "" {
		displayName = "Vắc xin Cúm"
	}

	if analysisDate.After(nextAnniversary) || analysisDate.Equal(nextAnniversary) {
		return []models.MissingItem{{
			VaccineName:          displayName,
			Description:          fmt.Sprintf("%s - Cần tiêm nhắc lại hàng năm (đã đến lịch).", displayName),
			EarliestNextDoseDate: &analysisDate,
			StatusTags:           []string{"due", "flu_annual"},
		}}
	}

	// Python Logic: Return as info if upcoming
	return []models.MissingItem{{
		VaccineName:          displayName,
		Description:          fmt.Sprintf("%s - Lịch tiêm nhắc lại hàng năm tiếp theo.", displayName),
		EarliestNextDoseDate: &nextAnniversary,
		StatusTags:           []string{"info", "booster_upcoming"},
	}}
}

// CheckCumulativeGroup handles groups where doses of different members count toward the same total (e.g., MMR).
func CheckCumulativeGroup(rule *models.VaccineRule, adminMap map[string][]models.AdministeredDose, dob time.Time, analysisDate time.Time) []models.MissingItem {
	// Collect doses using the centralized helper
	allDoses := utils.GetDosesForRule(rule, adminMap)

	// Delegate to AgeDependent logic but with cumulative doses
	// We wrap the cumulative doses in a temporary admin map
	displayName := rule.GroupDisplayName
	if displayName == "" {
		displayName = rule.DisplayName
	}
	
	tempAdminMap := make(map[string][]models.AdministeredDose)
	allNames := rule.NamesNorm
	allNames = append(allNames, rule.NamesNormGroup...)
	for _, name := range allNames {
		tempAdminMap[name] = allDoses
	}
	tempAdminMap[displayName] = allDoses
	
	// Create a shallow copy of rule to use the group display name as the key
	ruleCopy := *rule
	ruleCopy.DisplayName = displayName

	return CheckAgeDependentSeries(&ruleCopy, tempAdminMap, dob, analysisDate)
}


