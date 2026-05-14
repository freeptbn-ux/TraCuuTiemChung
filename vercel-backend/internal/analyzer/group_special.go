package analyzer

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"vercel-backend/internal/models"
)

// checkMMREquivalentGroup xử lý nhóm vắc xin MMR và kiểm tra khoảng cách với Sởi đơn (MVVAC).
func (e *Engine) checkMMREquivalentGroup(ruleKey string, rule VaccineRule, administeredMap map[string][]models.VaccineRecord) *AnalysisResult {
	// 1. Tìm tất cả các bản ghi MMR
	var mmrRecords []models.VaccineRecord
	for _, reg := range rule.Regimens {
		recs := e.getMatchingRecords(reg.NamesNorm, administeredMap)
		mmrRecords = append(mmrRecords, recs...)
	}
	sort.Slice(mmrRecords, func(i, j int) bool {
		return mmrRecords[i].Date.Before(mmrRecords[j].Date)
	})

	// 2. Tìm các bản ghi Sởi đơn (MVVAC)
	mvvacRule, ok := e.Rules["MVVAC"]
	var mvvacRecords []models.VaccineRecord
	if ok {
		mvvacRecords = e.getMatchingRecords(mvvacRule.NamesNorm, administeredMap)
	}

	numDoses := len(mmrRecords)
	monthsNow, _, _ := GetAgeAtDate(e.DOB, e.AnalysisDate)

	// Nếu chưa tiêm MMR nào, kiểm tra khoảng cách với MVVAC cuối cùng
	var earliestMMRDate time.Time
	if numDoses == 0 {
		if len(mvvacRecords) > 0 {
			lastMVVAC := mvvacRecords[len(mvvacRecords)-1]
			// Quy tắc: MMR cách MVVAC ít nhất 84 ngày
			earliestMMRDate = lastMVVAC.Date.AddDate(0, 0, 84)
		} else {
			// Nếu chưa tiêm Sởi đơn, MMR bắt đầu từ 9 tháng (theo rule)
			earliestMMRDate = AddMonths(e.DOB, rule.MinAgeMonthsOverallGroup)
		}
	}

	// Xác định phác đồ phù hợp
	var applicableRegimen *Course
	if numDoses == 0 {
		for i := range rule.Regimens {
			reg := &rule.Regimens[i]
			if monthsNow >= reg.MinAgeMonthsAtFirstDose && (reg.MaxAgeMonthsAtFirstDose == 0 || monthsNow <= reg.MaxAgeMonthsAtFirstDose) {
				applicableRegimen = reg
				break
			}
		}
	} else {
		// Dựa vào tuổi lúc tiêm mũi 1
		monthsAtFirst, _, _ := GetAgeAtDate(e.DOB, mmrRecords[0].Date)
		for i := range rule.Regimens {
			reg := &rule.Regimens[i]
			if monthsAtFirst >= reg.MinAgeMonthsAtFirstDose && (reg.MaxAgeMonthsAtFirstDose == 0 || monthsAtFirst <= reg.MaxAgeMonthsAtFirstDose) {
				applicableRegimen = reg
				break
			}
		}
	}

	if applicableRegimen == nil {
		return nil
	}

	if numDoses >= applicableRegimen.DosesRequired {
		return nil
	}

	// Tính ngày tiêm tiếp theo
	var nextDoseDate time.Time
	if numDoses == 0 {
		nextDoseDate = earliestMMRDate
		// Nếu tuổi hiện tại vẫn chưa đến tuổi bắt đầu của regimen
		startAgeDate := AddMonths(e.DOB, applicableRegimen.MinAgeMonthsAtFirstDose)
		if nextDoseDate.Before(startAgeDate) {
			nextDoseDate = startAgeDate
		}
	} else {
		lastDoseDate := mmrRecords[numDoses-1].Date
		intervalDays := 0
		if numDoses < len(applicableRegimen.MinIntervalDays) && applicableRegimen.MinIntervalDays[numDoses] != nil {
			intervalDays = *applicableRegimen.MinIntervalDays[numDoses]
		}
		nextDoseDate = lastDoseDate.AddDate(0, 0, intervalDays)
		
		// Kiểm tra quy tắc đặc biệt cho mũi 2/3 (ví dụ: MMR mũi nhắc lúc 4-7 tuổi)
		doseNumStr := fmt.Sprintf("%d", numDoses+1)
		if spec, ok := applicableRegimen.DoseSpecificRules[doseNumStr]; ok {
			if spec.AlternativeMinAgeYears > 0 {
				minAgeDate := AddYears(e.DOB, spec.AlternativeMinAgeYears)
				if nextDoseDate.Before(minAgeDate) {
					nextDoseDate = minAgeDate
				}
			}
		}
	}

	tags := []string{"due"}
	if e.AnalysisDate.Before(nextDoseDate) {
		tags = []string{"too_young"}
	}

	return &AnalysisResult{
		VaccineNameForPopup:  rule.GroupDisplayName,
		Description:          fmt.Sprintf("%s (Tiêm mũi %d)", rule.GroupDisplayName, numDoses+1),
		EarliestNextDoseDate: &nextDoseDate,
		StatusTags:           tags,
	}
}

// checkFluGroup xử lý nhóm vắc xin Cúm với nhắc lại hàng năm.
func (e *Engine) checkFluGroup(ruleKey string, rule VaccineRule, administeredMap map[string][]models.VaccineRecord) *AnalysisResult {
	// 1. Tìm tất cả bản ghi Cúm dựa trên keywords
	var fluRecords []models.VaccineRecord
	seenDates := make(map[string]bool)
	
	// Normalize keywords
	var keywords []string
	for _, kw := range rule.RecognitionKeywords {
		keywords = append(keywords, strings.ToLower(kw))
	}

	for normName, recs := range administeredMap {
		match := false
		for _, kw := range keywords {
			if strings.Contains(normName, kw) {
				match = true
				break
			}
		}
		if match {
			for _, rec := range recs {
				dateStr := rec.Date.Format("2006-01-02")
				if !seenDates[dateStr] {
					fluRecords = append(fluRecords, rec)
					seenDates[dateStr] = true
				}
			}
		}
	}

	sort.Slice(fluRecords, func(i, j int) bool {
		return fluRecords[i].Date.Before(fluRecords[j].Date)
	})

	numDoses := len(fluRecords)
	monthsNow, _, _ := GetAgeAtDate(e.DOB, e.AnalysisDate)

	if monthsNow < rule.MinAgeMonthsAtFirstDose {
		return nil
	}

	if numDoses == 0 {
		earliestDate := AddMonths(e.DOB, rule.MinAgeMonthsAtFirstDose)
		tags := []string{"due"}
		if e.AnalysisDate.Before(earliestDate) {
			tags = []string{"too_young"}
		}
		return &AnalysisResult{
			VaccineNameForPopup:  rule.GroupDisplayName,
			Description:          fmt.Sprintf("%s (Bắt đầu tiêm mũi 1)", rule.GroupDisplayName),
			EarliestNextDoseDate: &earliestDate,
			StatusTags:           tags,
		}
	}

	// Đã tiêm ít nhất 1 mũi
	firstDoseDate := fluRecords[0].Date
	monthsAtFirst, _, _ := GetAgeAtDate(e.DOB, firstDoseDate)
	
	// Kiểm tra xem có cần mũi 2 khởi đầu không (Trẻ < 9 tuổi lần đầu tiêm cúm cần 2 mũi)
	needsInitialSecondDose := monthsAtFirst < 108 // 9 years = 108 months

	if needsInitialSecondDose && numDoses == 1 {
		interval := rule.InitialSeriesIntervalDays
		if interval == 0 {
			interval = 30
		}
		nextDoseDate := firstDoseDate.AddDate(0, 0, interval)
		tags := []string{"due"}
		if e.AnalysisDate.Before(nextDoseDate) {
			tags = []string{"too_young"}
		}
		return &AnalysisResult{
			VaccineNameForPopup:  rule.GroupDisplayName,
			Description:          fmt.Sprintf("%s (Tiêm mũi 2 khởi đầu)", rule.GroupDisplayName),
			EarliestNextDoseDate: &nextDoseDate,
			StatusTags:           tags,
		}
	}

	// Nhắc lại hàng năm
	lastDoseDate := fluRecords[numDoses-1].Date
	nextDoseDate := lastDoseDate.AddDate(1, 0, 0) // +1 year
	
	tags := []string{"due"}
	if e.AnalysisDate.Before(nextDoseDate) {
		tags = []string{"too_young"}
	}

	return &AnalysisResult{
		VaccineNameForPopup:  rule.GroupDisplayName,
		Description:          fmt.Sprintf("%s (Tiêm mũi nhắc lại hàng năm)", rule.GroupDisplayName),
		EarliestNextDoseDate: &nextDoseDate,
		StatusTags:           tags,
	}
}

// checkMeningococcalACYWGroup xử lý nhóm vắc xin Não mô cầu ACYW và các tương tác.
func (e *Engine) checkMeningococcalACYWGroup(ruleKey string, rule VaccineRule, administeredMap map[string][]models.VaccineRecord) *AnalysisResult {
	// Thu thập bản ghi Menactra và MenQuadfi
	var allRecords []models.VaccineRecord
	memberMap := make(map[string]*Member)

	for _, member := range rule.Members {
		recs := e.getMatchingRecords(member.NamesNorm, administeredMap)
		allRecords = append(allRecords, recs...)
		for _, norm := range member.NamesNorm {
			memberMap[norm] = &member
		}
	}

	sort.Slice(allRecords, func(i, j int) bool {
		return allRecords[i].Date.Before(allRecords[j].Date)
	})

	numDoses := len(allRecords)
	monthsNow, _, _ := GetAgeAtDate(e.DOB, e.AnalysisDate)

	// Kiểm tra tương tác với các vắc xin khác
	var interactionWarnings []string
	var earliestByInteractions time.Time

	for otherRuleKey, inter := range rule.Interactions {
		if otherRule, ok := e.Rules[otherRuleKey]; ok {
			otherRecords := e.getMatchingRecords(otherRule.NamesNorm, administeredMap)
			if len(otherRecords) > 0 {
				lastOther := otherRecords[len(otherRecords)-1]
				if monthsNow >= inter.AppliesWhenAgeMonthsGte {
					earliestInter := lastOther.Date.AddDate(0, 0, inter.MinIntervalDays)
					if earliestInter.After(earliestByInteractions) {
						earliestByInteractions = earliestInter
					}
					if e.AnalysisDate.Before(earliestInter) {
						interactionWarnings = append(interactionWarnings, inter.Message)
					}
				}
			}
		}
	}

	var interactionWarning string
	if len(interactionWarnings) > 0 {
		interactionWarning = fmt.Sprintf(" [Lưu ý: %s]", strings.Join(interactionWarnings, "; "))
	}

	if numDoses == 0 {
		// Gợi ý MenQuadfi (ưu tiên vì tiêm sớm được)
		mq, ok := rule.Members["MENQUADFI"]
		if !ok { return nil }
		
		earliestDate := e.DOB.AddDate(0, 0, mq.MinAgeWeeksOverall*7)
		if !earliestByInteractions.IsZero() && earliestDate.Before(earliestByInteractions) {
			earliestDate = earliestByInteractions
		}

		tags := []string{"due"}
		if e.AnalysisDate.Before(earliestDate) {
			tags = []string{"too_young"}
		}

		return &AnalysisResult{
			VaccineNameForPopup:  rule.GroupDisplayName,
			Description:          fmt.Sprintf("%s (Bắt đầu tiêm)%s", rule.GroupDisplayName, interactionWarning),
			EarliestNextDoseDate: &earliestDate,
			StatusTags:           tags,
		}
	}

	// Đã tiêm. Xác định rule áp dụng của member đã tiêm mũi 1.
	firstRec := allRecords[0]
	normFirst := NormalizeVaccineName(firstRec.VaccineName)
	member := memberMap[normFirst]
	if member == nil { return nil }

	monthsAtFirst, _, _ := GetAgeAtDate(e.DOB, firstRec.Date)
	var applicableAgeRule *AgeRule
	for i := range member.RulesByAge {
		ar := &member.RulesByAge[i]
		if (ar.MinAgeMonthsAtFirstDose == 0 || monthsAtFirst >= ar.MinAgeMonthsAtFirstDose) &&
		   (ar.MaxAgeMonthsAtFirstDose == 0 || monthsAtFirst <= ar.MaxAgeMonthsAtFirstDose) &&
		   (ar.MinAgeWeeksAtFirstDose == 0 || int(firstRec.Date.Sub(e.DOB).Hours()/24/7) >= ar.MinAgeWeeksAtFirstDose) {
			applicableAgeRule = ar
			break
		}
	}

	if applicableAgeRule == nil { return nil }

	if numDoses >= applicableAgeRule.DosesRequired {
		// Kiểm tra booster
		if applicableAgeRule.Booster != nil {
			lastDoseDate := allRecords[numDoses-1].Date
			earliestBooster := lastDoseDate.AddDate(0, 0, applicableAgeRule.Booster.MinIntervalDaysFromLast)
			minAgeBooster := AddMonths(e.DOB, applicableAgeRule.Booster.MinAgeMonths)
			if earliestBooster.Before(minAgeBooster) {
				earliestBooster = minAgeBooster
			}

			if monthsNow >= applicableAgeRule.Booster.MinAgeMonths || numDoses == applicableAgeRule.DosesRequired {
				tags := []string{"due"}
				if e.AnalysisDate.Before(earliestBooster) {
					tags = []string{"too_young"}
				}
				return &AnalysisResult{
					VaccineNameForPopup:  member.Display,
					Description:          fmt.Sprintf("%s (%s)%s", member.Display, applicableAgeRule.Booster.Description, interactionWarning),
					EarliestNextDoseDate: &earliestBooster,
					StatusTags:           tags,
				}
			}
		}
		return nil
	}

	// Tiếp tục phác đồ
	lastDoseDate := allRecords[numDoses-1].Date
	intervalDays := 0
	if numDoses < len(applicableAgeRule.MinIntervalDays) && applicableAgeRule.MinIntervalDays[numDoses] != nil {
		intervalDays = *applicableAgeRule.MinIntervalDays[numDoses]
	}
	nextDoseDate := lastDoseDate.AddDate(0, 0, intervalDays)
	if !earliestByInteractions.IsZero() && nextDoseDate.Before(earliestByInteractions) {
		nextDoseDate = earliestByInteractions
	}

	tags := []string{"due"}
	if e.AnalysisDate.Before(nextDoseDate) {
		tags = []string{"too_young"}
	}

	return &AnalysisResult{
		VaccineNameForPopup:  member.Display,
		Description:          fmt.Sprintf("%s (Tiêm mũi %d)%s", member.Display, numDoses+1, interactionWarning),
		EarliestNextDoseDate: &nextDoseDate,
		StatusTags:           tags,
	}
}
