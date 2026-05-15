package analyzer

import (
	"fmt"
	"strings"
	"time"

	"vercel-backend/pkg/models"
)

// checkAlternativeCoursesMinAgeGroup xử lý nhóm vắc xin có nhiều phác đồ thay thế dựa trên độ tuổi bắt đầu (ví dụ: Rota).
func (e *Engine) checkAlternativeCoursesMinAgeGroup(ruleKey string, rule VaccineRule, administeredMap map[string][]models.VaccineRecord) []AnalysisResult {
	var potentialIncompleteMessages []AnalysisResult
	anyCourseCompleted := false

	currentAgeMonths, _, _ := GetAgeAtDate(e.DOB, e.AnalysisDate)

	// Thu thập tất cả bản ghi của group
	allRecords := e.getAllMatchingRecordsForGroup(rule, administeredMap)

	// 1. Kiểm tra tính hợp lệ tuổi mũi 1 của cả group
	if len(allRecords) > 0 {
		valid, errRes := e.checkFirstDoseAgeValidity(allRecords[0].Date, rule, rule.GroupDisplayName)
		if !valid {
			return []AnalysisResult{*errRes}
		}
	}

	// 2. Kiểm tra quá tuổi bắt đầu (nếu chưa tiêm mũi nào)
	if len(allRecords) == 0 && rule.MaxAgeMonthsToStartFirstDoseGroup > 0 {
		if currentAgeMonths > rule.MaxAgeMonthsToStartFirstDoseGroup {
			return []AnalysisResult{{
				VaccineNameForPopup: rule.GroupDisplayName,
				Description:         fmt.Sprintf("%s: Đã qua %d tháng tuổi, không còn chỉ định bắt đầu.", rule.GroupDisplayName, rule.MaxAgeMonthsToStartFirstDoseGroup),
				StatusTags:          []string{"too_old_to_start"},
			}}
		}
	}

	// 3. Duyệt qua từng phác đồ
	isRota := ruleKey == "Rota"
	for _, course := range rule.Courses {
		courseDisplayNameFull := fmt.Sprintf("%s - %s", rule.GroupDisplayName, course.Display)
		tempRule := VaccineRule{
			DisplayName:             courseDisplayNameFull,
			GroupDisplayName:        rule.GroupDisplayName,
			DosesRequired:           course.DosesRequired,
			MinIntervalDays:         course.MinIntervalDays,
			MinAgeMonthsAtFirstDose: course.MinAgeMonthsAtFirstDose,
			MaxAgeMonthsAtFirstDose: course.MaxAgeMonthsAtFirstDose,
			MinAgeWeeksAtFirstDose:  rule.MinAgeWeeksAtFirstDose,
			NamesNorm:               course.NamesNorm,
			DoseSpecificRules:       course.DoseSpecificRules,
		}

		courseResults := e.checkSingleSeries(ruleKey+"_"+course.Display, tempRule, administeredMap)

		if len(courseResults) == 0 {
			anyCourseCompleted = true
			break
		}

		// Rota transform logic: Nếu quá tuổi hoàn thành, chuyển "due/eligible" thành "too_old_to_complete"
		if isRota && rule.MaxAgeMonthsForCompletionGroup > 0 && currentAgeMonths > rule.MaxAgeMonthsForCompletionGroup {
			courseRecords := e.getMatchingRecords(course.NamesNorm, administeredMap)
			if len(courseRecords) > 0 {
				for i := range courseResults {
					if containsStr(courseResults[i].StatusTags, "due") || containsStr(courseResults[i].StatusTags, "eligible") {
						courseResults[i].Description = fmt.Sprintf("%s: Đã trên %d tháng tuổi, không còn chỉ định hoàn thành phác đồ.", course.Display, rule.MaxAgeMonthsForCompletionGroup)
						courseResults[i].StatusTags = []string{"info", "too_old_to_complete", "rota_too_old_to_complete"}
						courseResults[i].EarliestNextDoseDate = nil
					}
				}
			}
		}

		potentialIncompleteMessages = append(potentialIncompleteMessages, courseResults...)
	}

	if anyCourseCompleted {
		return nil
	}

	// 4. Nếu chưa tiêm mũi nào
	if len(allRecords) == 0 {
		var options []string
		var earliestDates []*time.Time
		for _, course := range rule.Courses {
			options = append(options, course.Display)
			tempCourseRule := VaccineRule{
				MinAgeMonthsAtFirstDose: course.MinAgeMonthsAtFirstDose,
				MinAgeWeeksAtFirstDose:  rule.MinAgeWeeksAtFirstDose,
			}
			_, ed, _ := e.getAgeStatusAndEarliestDate(tempCourseRule, "")
			if ed != nil {
				earliestDates = append(earliestDates, ed)
			}
		}

		statusMsg, earliestNext, tags := e.getAgeStatusAndEarliestDate(rule, "")

		// Thay thế 'eligible' thành 'due' cho các mũi chưa tiêm
		for i, t := range tags {
			if t == "eligible" {
				tags[i] = "due"
			}
		}

		// Tìm ngày tiêm sớm nhất từ các phác đồ (nếu đang too_young)
		if containsStr(tags, "too_young") && len(earliestDates) > 0 {
			minDate := *earliestDates[0]
			for _, d := range earliestDates {
				if d.Before(minDate) {
					minDate = *d
				}
			}
			earliestNext = &minDate
		}

		return []AnalysisResult{{
			VaccineNameForPopup:  rule.GroupDisplayName,
			Description:          fmt.Sprintf("%s (Lựa chọn: %s). %s", rule.GroupDisplayName, strings.Join(options, " Hoặc "), statusMsg),
			EarliestNextDoseDate: earliestNext,
			StatusTags:           tags,
			IsMissing:            true,
			DoseNumber:           1,
		}}
	}

	// 5. Nếu đã tiêm, lọc kết quả theo phác đồ đã bắt đầu
	var filteredResults []AnalysisResult
	for _, res := range potentialIncompleteMessages {
		for _, course := range rule.Courses {
			if res.VaccineNameForPopup == fmt.Sprintf("%s - %s", rule.GroupDisplayName, course.Display) {
				recs := e.getMatchingRecords(course.NamesNorm, administeredMap)
				if len(recs) > 0 {
					filteredResults = append(filteredResults, res)
					break
				}
			}
		}
	}

	if len(filteredResults) > 0 {
		return filteredResults
	}

	return potentialIncompleteMessages
}

// checkAlternativeCoursesAgeRangeGroup xử lý nhóm vắc xin có nhiều phác đồ thay thế dựa trên khoảng độ tuổi (ví dụ: JE, HepA).
func (e *Engine) checkAlternativeCoursesAgeRangeGroup(ruleKey string, rule VaccineRule, administeredMap map[string][]models.VaccineRecord) []AnalysisResult {
	allRecords := e.getAllMatchingRecordsForGroup(rule, administeredMap)
	numDoses := len(allRecords)
	currentAgeMonths, _, _ := GetAgeAtDate(e.DOB, e.AnalysisDate)

	var jeSpecialResults []AnalysisResult
	coursesToSkip := make(map[string]bool)

	// 1. Logic đặc biệt cho nhóm Nhật Bản B (JE_Group)
	if ruleKey == "JE_Group" {
		res, skipped, shouldStop := e.handleJESkipLogic(rule, administeredMap)
		if shouldStop {
			return res
		}
		jeSpecialResults = res
		for _, s := range skipped {
			coursesToSkip[s] = true
		}
	}

	// 2. Trường hợp chưa tiêm mũi nào
	if numDoses == 0 {
		var eligibleOptions []string
		var earliestDates []*time.Time

		for _, course := range rule.Courses {
			// Kiểm tra xem độ tuổi hiện tại có nằm trong dải tuổi bắt đầu của phác đồ không
			if currentAgeMonths >= course.MinAgeMonthsAtFirstDose && (course.MaxAgeYearsAtFirstDose == 0 || currentAgeMonths <= course.MaxAgeYearsAtFirstDose*12) {
				eligibleOptions = append(eligibleOptions, course.Display)

				tempRule := VaccineRule{
					MinAgeMonthsAtFirstDose: course.MinAgeMonthsAtFirstDose,
					MinAgeYearsAtFirstDose:  course.MinAgeYearsAtFirstDose,
				}
				_, ed, _ := e.getAgeStatusAndEarliestDate(tempRule, "")
				if ed != nil {
					earliestDates = append(earliestDates, ed)
				}
			}
		}

		if len(eligibleOptions) == 0 {
			// Tìm tuổi tối thiểu của tất cả các phác đồ
			minAge := 9999
			for _, course := range rule.Courses {
				if course.MinAgeMonthsAtFirstDose < minAge {
					minAge = course.MinAgeMonthsAtFirstDose
				}
			}

			if currentAgeMonths < minAge {
				// Cập nhật rule tạm thời để getAgeStatus trả về too_young
				tempRule := rule
				tempRule.MinAgeMonthsOverallGroup = minAge
				statusMsg, earliestNext, tags := e.getAgeStatusAndEarliestDate(tempRule, "")
				return []AnalysisResult{{
					VaccineNameForPopup:  rule.GroupDisplayName,
					Description:          fmt.Sprintf("%s. %s", rule.GroupDisplayName, statusMsg),
					EarliestNextDoseDate: earliestNext,
					StatusTags:           tags,
					IsMissing:            true,
					DoseNumber:           1,
				}}
			}

			return []AnalysisResult{{
				VaccineNameForPopup: rule.GroupDisplayName,
				Description:         fmt.Sprintf("%s (Hiện không có phác đồ phù hợp với lứa tuổi)", rule.GroupDisplayName),
				StatusTags:          []string{"too_old"},
			}}
		}

		statusMsg, earliestNext, tags := e.getAgeStatusAndEarliestDate(rule, "")
		// Thay thế 'eligible' thành 'due'
		for i, t := range tags {
			if t == "eligible" {
				tags[i] = "due"
			}
		}
		if containsStr(tags, "too_young") && len(earliestDates) > 0 {
			minDate := *earliestDates[0]
			for _, d := range earliestDates {
				if d.Before(minDate) {
					minDate = *d
				}
			}
			earliestNext = &minDate
		}

		return []AnalysisResult{{
			VaccineNameForPopup:  rule.GroupDisplayName,
			Description:          fmt.Sprintf("%s (Lựa chọn: %s). %s", rule.GroupDisplayName, strings.Join(eligibleOptions, " Hoặc "), statusMsg),
			EarliestNextDoseDate: earliestNext,
			StatusTags:           tags,
			IsMissing:            true,
			DoseNumber:           1,
		}}
	}

	// 3. Delegate xuống checkSingleSeries cho các phác đồ đã bắt đầu
	anyCourseCompleted := false
	var loopResults []AnalysisResult

	for _, course := range rule.Courses {
		isSkipped := false
		displayNorm := NormalizeForMatch(course.Display)
		for skipName := range coursesToSkip {
			if NormalizeForMatch(skipName) == displayNorm {
				isSkipped = true
				break
			}
		}
		
		if isSkipped {
			continue
		}

		courseRecords := e.getMatchingRecords(course.NamesNorm, administeredMap)
		if len(courseRecords) == 0 {
			continue
		}

		tempRule := VaccineRule{
			DisplayName:             fmt.Sprintf("%s - %s", rule.GroupDisplayName, course.Display),
			GroupDisplayName:        rule.GroupDisplayName,
			DosesRequired:           course.DosesRequired,
			MinIntervalDays:         course.MinIntervalDays,
			MinAgeMonthsAtFirstDose: course.MinAgeMonthsAtFirstDose,
			MaxAgeMonthsAtFirstDose: course.MaxAgeMonthsAtFirstDose,
			MinAgeYearsAtFirstDose:  course.MinAgeYearsAtFirstDose,
			BoosterIntervalYears:    course.BoosterIntervalYears,
			BoosterAfterDoseNumber:  course.BoosterAfterDoseNumber,
			BoosterMaxAgeYears:      course.BoosterMaxAgeYears,
			NamesNorm:               course.NamesNorm,
			Type:                    RuleTypeSingleSeries,
		}

		courseRes := e.checkSingleSeries("", tempRule, administeredMap)
		if len(courseRes) == 0 {
			// Only consider as group completion if this course was NOT skipped
			anyCourseCompleted = true
			break
		}
		loopResults = append(loopResults, courseRes...)
	}

	if anyCourseCompleted {
		return jeSpecialResults
	}

	if numDoses > 0 && len(loopResults) == 0 {
		loopResults = append(loopResults, AnalysisResult{
			VaccineNameForPopup: rule.GroupDisplayName,
			Description:         fmt.Sprintf("%s (Các mũi đã tiêm không phù hợp với phác đồ nào trong khoảng tuổi cho phép hoặc đã qua tuổi. Kiểm tra lại.)", rule.GroupDisplayName),
			StatusTags:          []string{"error_ambiguous_course"},
		})
	}

	jeCombinedResults := append(jeSpecialResults, loopResults...)
	return jeCombinedResults
}

// handleJESkipLogic xử lý các trường hợp đặc biệt của nhóm Viêm não Nhật Bản và xác định các phác đồ cần bỏ qua
func (e *Engine) handleJESkipLogic(rule VaccineRule, administeredMap map[string][]models.VaccineRecord) ([]AnalysisResult, []string, bool) {
	jevaxRecs := e.getMatchingRecordsByDisplay("Jevax", rule, administeredMap)
	imojevRecs := e.getMatchingRecordsByDisplay("Imojev", rule, administeredMap)
	jeevRecs := e.getMatchingRecordsByDisplay("JEEV", rule, administeredMap)

	numJevax := len(jevaxRecs)
	numImojev := len(imojevRecs)
	numJEEV := len(jeevRecs)

	var results []AnalysisResult
	var skipped []string

	// 2. JEEV check
	if numJevax > 0 && numJEEV > 0 {
		if jeevRecs[0].Date.After(jevaxRecs[len(jevaxRecs)-1].Date) {
			results = append(results, AnalysisResult{
				VaccineNameForPopup: rule.GroupDisplayName,
				Description:         fmt.Sprintf("%s - Cảnh báo: Đã tiêm VNNB (Jevax) sau đó chuyển sang JEEV. Sẽ ưu tiên nhắc lịch theo phác đồ JEEV.", rule.GroupDisplayName),
				StatusTags:          []string{"info", "je_mixed_warning", "switched_to_jeev"},
			})
			skipped = append(skipped, "Jevax/VNNB (Việt Nam)")
		}
	}

	// 3. Imojev check
	if numImojev > 0 && numJevax > 0 {
		if imojevRecs[0].Date.Before(jevaxRecs[0].Date) {
			results = append(results, AnalysisResult{
				VaccineNameForPopup: rule.GroupDisplayName,
				Description:         fmt.Sprintf("%s - Cảnh báo: Đã bắt đầu với Imojev, không nên tiêm Jevax/VNNB. Cần hoàn thành phác đồ Imojev.", rule.GroupDisplayName),
				StatusTags:          []string{"error_interchange", "je_mixed_warning"},
			})
			skipped = append(skipped, "Jevax/VNNB (Việt Nam)")
		} else if imojevRecs[0].Date.After(jevaxRecs[len(jevaxRecs)-1].Date) {
			results = append(results, AnalysisResult{
				VaccineNameForPopup: rule.GroupDisplayName,
				Description:         fmt.Sprintf("%s - Cảnh báo: Đã tiêm VNNB (Jevax) sau đó chuyển sang Imojev. Sẽ ưu tiên nhắc lịch theo phác đồ Imojev.", rule.GroupDisplayName),
				StatusTags:          []string{"info", "je_mixed_warning", "switched_to_imojev"},
			})
			skipped = append(skipped, "Jevax/VNNB (Việt Nam)")
		}
	}

	// 1. Trường hợp đã hoàn thành (Mix 3 Jevax + 1 Imojev)
	if numJevax >= 3 && numImojev >= 1 {
		return results, skipped, true
	}

	// 4. Jevax completion booster logic
	if numJevax >= 3 && numImojev == 0 && numJEEV == 0 {
		lastJevaxDate := jevaxRecs[len(jevaxRecs)-1].Date
		nextBoosterDate := AddYears(lastJevaxDate, 3)

		finalDate := nextBoosterDate
		if e.AnalysisDate.After(nextBoosterDate) {
			finalDate = e.AnalysisDate
		}

		tags := []string{"due", "booster_due"}
		if e.AnalysisDate.Before(nextBoosterDate) {
			tags = []string{"info", "booster_upcoming"}
		}

		results = append(results, AnalysisResult{
			VaccineNameForPopup:  rule.GroupDisplayName,
			Description:          fmt.Sprintf("%s - Cần tiêm mũi nhắc lại. Tùy chọn: dùng Jevax (3 năm/lần đến 15 tuổi) hoặc tiêm 1 mũi Imojev để hoàn thành.", rule.GroupDisplayName),
			EarliestNextDoseDate: &finalDate,
			StatusTags:           tags,
		})

		skipped = append(skipped, "Jevax/VNNB (Việt Nam)", "Imojev (Sanofi Pasteur)")
	} else if 0 < numJevax && numJevax < 3 && numImojev == 0 && numJEEV == 0 {
		// Optional switch logic
		if numJevax == 1 {
			nextDate := jevaxRecs[0].Date.AddDate(0, 0, 30)
			if nextDate.Before(e.AnalysisDate) {
				nextDate = e.AnalysisDate
			}
			results = append(results, AnalysisResult{
				VaccineNameForPopup:  "Tùy chọn VNNB -> Imojev",
				Description:          fmt.Sprintf("%s - Tùy chọn: Sau 1 mũi Jevax, có thể chuyển sang tiêm 2 mũi Imojev (mũi Imojev đầu tiên cách mũi Jevax 1 tháng).", rule.GroupDisplayName),
				EarliestNextDoseDate: &nextDate,
				StatusTags:           []string{"info", "alternative_course"},
			})
		} else if numJevax == 2 {
			nextDate := AddYears(jevaxRecs[1].Date, 1)
			if nextDate.Before(e.AnalysisDate) {
				nextDate = e.AnalysisDate
			}
			results = append(results, AnalysisResult{
				VaccineNameForPopup:  "Tùy chọn VNNB -> Imojev",
				Description:          fmt.Sprintf("%s - Tùy chọn: Sau 2 mũi Jevax, có thể tiêm 1 mũi Imojev sau 1 năm để hoàn thành.", rule.GroupDisplayName),
				EarliestNextDoseDate: &nextDate,
				StatusTags:           []string{"info", "alternative_course"},
			})
		}
	}

	return results, skipped, false
}

// getAllMatchingRecordsForGroup lấy tất cả bản ghi cho cả nhóm vắc xin
func (e *Engine) getAllMatchingRecordsForGroup(rule VaccineRule, administeredMap map[string][]models.VaccineRecord) []models.VaccineRecord {
	var allNames []string
	for _, course := range rule.Courses {
		allNames = append(allNames, course.NamesNorm...)
	}
	return e.getMatchingRecords(allNames, administeredMap)
}

// getMatchingRecordsByDisplay tìm bản ghi dựa trên chuỗi hiển thị của phác đồ
func (e *Engine) getMatchingRecordsByDisplay(display string, rule VaccineRule, administeredMap map[string][]models.VaccineRecord) []models.VaccineRecord {
	for _, course := range rule.Courses {
		if strings.Contains(strings.ToLower(course.Display), strings.ToLower(display)) {
			return e.getMatchingRecords(course.NamesNorm, administeredMap)
		}
	}
	return nil
}
