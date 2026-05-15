package analyzer

import (
	"fmt"
	"vercel-backend/pkg/models"
)

// checkCumulativeGroupDoses xử lý RuleTypeGroupCumulativeUnique.
// Đếm tổng số liều đã tiêm trong nhóm và kiểm tra xem đã đủ số lượng yêu cầu chưa.
func (e *Engine) checkCumulativeGroupDoses(ruleKey string, rule VaccineRule, administeredMap map[string][]models.VaccineRecord) []AnalysisResult {
	// 1. Thu thập tất cả bản ghi trong nhóm
	records := e.getMatchingRecords(rule.NamesNormGroup, administeredMap)
	numDoses := len(records)
	
	displayName := rule.DisplayName
	if displayName == "" {
		displayName = rule.GroupDisplayName
	}

	// 2. Nếu đã đủ liều
	if numDoses >= rule.DosesRequired && rule.DosesRequired > 0 {
		return nil
	}

	// 3. Nếu chưa tiêm mũi nào
	if numDoses == 0 {
		statusMsg, earliestDate, ageTags := e.getAgeStatusAndEarliestDate(rule, "")
		desc := fmt.Sprintf("%s (Chưa tiêm - cần tổng %d liều). %s", displayName, rule.DosesRequired, statusMsg)

		return []AnalysisResult{{
			VaccineNameForPopup:  displayName,
			Description:          desc,
			EarliestNextDoseDate: earliestDate,
			StatusTags:           ageTags,
			IsMissing:            true,
			DoseNumber:           1,
		}}
	}

	// 4. Nếu đã tiêm nhưng chưa đủ
	// Kiểm tra tính hợp lệ của mũi 1
	valid, errRes := e.checkFirstDoseAgeValidity(records[0].Date, rule, displayName)
	if !valid {
		return []AnalysisResult{*errRes}
	}

	// Tính ngày tiêm tiếp theo
	lastDoseDate := records[numDoses-1].Date
	
	intervalDays := 0
	if numDoses < len(rule.MinIntervalDays) && rule.MinIntervalDays[numDoses] != nil {
		intervalDays = *rule.MinIntervalDays[numDoses]
	}
	earliestDate := lastDoseDate.AddDate(0, 0, intervalDays)
	intervalDesc := formatIntervalDescription(intervalDays)

	finalEarliestDate := earliestDate
	if finalEarliestDate.Before(e.AnalysisDate) {
		finalEarliestDate = e.AnalysisDate
	}

	nextDoseNum := numDoses + 1
	status := fmt.Sprintf("tiêm thêm %d liều để hoàn thành (tổng %d liều)", rule.DosesRequired-numDoses, rule.DosesRequired)
	if intervalDesc != "" {
		status += fmt.Sprintf(" - mũi tiếp theo cách mũi gần nhất %s", intervalDesc)
	}

	tags := []string{"due"}
	if e.AnalysisDate.Before(earliestDate) {
		tags = []string{"too_young"}
	}

	return []AnalysisResult{{
		VaccineNameForPopup:  displayName,
		Description:          fmt.Sprintf("%s (%s)", displayName, status),
		EarliestNextDoseDate: &finalEarliestDate,
		StatusTags:           tags,
		DoseNumber:           nextDoseNum,
	}}
}
