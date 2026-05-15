package analyzer

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"vercel-backend/pkg/models"
)

// checkMMREquivalentGroup xử lý nhóm vắc xin MMR và kiểm tra khoảng cách với Sởi đơn (MVVAC).
func (e *Engine) checkMMREquivalentGroup(ruleKey string, rule VaccineRule, administeredMap map[string][]models.VaccineRecord) []AnalysisResult {
	// 1. Tìm tất cả các bản ghi MMR từ các thành viên trong nhóm
	var mmrRecords []models.VaccineRecord
	mmrRecords = e.getMatchingRecords(rule.NamesNormGroup, administeredMap)
	
	sort.Slice(mmrRecords, func(i, j int) bool {
		return mmrRecords[i].Date.Before(mmrRecords[j].Date)
	})

	// 2. Tìm các bản ghi Sởi đơn (MVVAC)
	mvvacRule, hasMVVACRule := e.Rules["MVVAC"]
	var mvvacRecords []models.VaccineRecord
	if hasMVVACRule {
		mvvacRecords = e.getMatchingRecords(mvvacRule.NamesNorm, administeredMap)
	}

	numDoses := len(mmrRecords)
	displayName := rule.GroupDisplayName

	// logic tương tác MVVAC
	if len(mvvacRecords) > 0 {
		lastMVVACDate := mvvacRecords[len(mvvacRecords)-1].Date
		
		if numDoses == 0 {
			// Chưa tiêm MMR, đã tiêm MVVAC
			earliestByMVVAC := lastMVVACDate.AddDate(0, 0, 84)
			
			// Lấy earliest by age
			_, earliestByAge, _ := e.getAgeStatusAndEarliestDate(rule, "")
			
			finalEarliest := earliestByMVVAC
			if earliestByAge != nil && earliestByAge.After(finalEarliest) {
				finalEarliest = *earliestByAge
			}

			statusTags := []string{"due"}
			if e.AnalysisDate.Before(finalEarliest) {
				statusTags = []string{"info", "scheduled"}
			}

			desc := fmt.Sprintf("%s (Chưa tiêm). Trẻ đã tiêm Sởi đơn (MVVAC) lúc %s. Theo quy định, mũi Sởi-Quai bị-Rubella cần cách mũi Sởi đơn ít nhất 3 tháng (84 ngày).", displayName, lastMVVACDate.Format("02/01/2006"))
			
			return []AnalysisResult{{
				VaccineNameForPopup:  displayName,
				Description:          desc,
				EarliestNextDoseDate: &finalEarliest,
				StatusTags:           statusTags,
				IsMissing:            true,
				DoseNumber:           1,
			}}

		} else if numDoses >= 1 {
			// Đã tiêm ít nhất 1 mũi MMR sau/cùng MVVAC
			firstMMRDate := mmrRecords[0].Date
			actualInterval := int(firstMMRDate.Sub(lastMVVACDate).Hours() / 24)
			
			var results []AnalysisResult
			if actualInterval < 84 && actualInterval >= 0 {
				results = append(results, AnalysisResult{
					VaccineNameForPopup: displayName,
					Description:         fmt.Sprintf("%s - Cảnh báo: Mũi MMR đầu tiên tiêm cách mũi Sởi đơn (MVVAC) chỉ %d ngày (quy định là 84 ngày). Hiệu quả có thể bị ảnh hưởng.", displayName, actualInterval),
					StatusTags:          []string{"warning", "interval_violation_mvvac_mmr"},
				})
			}
			
			if numDoses == 1 {
				nextDueDate := AddYears(firstMMRDate, 3)
				tags := []string{"due", "booster_due"}
				if e.AnalysisDate.Before(nextDueDate) {
					tags = []string{"info", "booster_upcoming"}
				} else {
					nextDueDate = e.AnalysisDate
				}

				results = append(results, AnalysisResult{
					VaccineNameForPopup:  displayName,
					Description:          fmt.Sprintf("%s - Cần tiêm mũi 2 (phác đồ MVVAC + MMR) sau 3 năm kể từ mũi MMR/Priorix đầu tiên.", displayName),
					EarliestNextDoseDate: &nextDueDate,
					StatusTags:           tags,
				})
			}
			return results
		}
	}

	// Trường hợp không có MVVAC hoặc không rơi vào case đặc biệt trên
	if numDoses == 0 {
		statusMsg, earliestDate, ageTags := e.getAgeStatusAndEarliestDate(rule, "")
		return []AnalysisResult{{
			VaccineNameForPopup:  displayName,
			Description:          fmt.Sprintf("%s (Chưa tiêm - cần 2 liều). %s", displayName, statusMsg),
			EarliestNextDoseDate: earliestDate,
			StatusTags:           ageTags,
			IsMissing:            true,
			DoseNumber:           1,
		}}
	}

	// Đã tiêm, check phác đồ thông qua checkSingleSeries
	return e.checkSingleSeries(ruleKey, rule, administeredMap)
}

// checkFluGroup xử lý nhóm vắc xin Cúm với nhắc lại hàng năm.
func (e *Engine) checkFluGroup(ruleKey string, rule VaccineRule, administeredMap map[string][]models.VaccineRecord) []AnalysisResult {
	// 1. Tìm tất cả bản ghi Cúm dựa trên keywords của RAW NAME
	var fluRecords []models.VaccineRecord
	seenDates := make(map[string]bool)
	
	var keywords []string
	for _, kw := range rule.RecognitionKeywords {
		keywords = append(keywords, strings.ToLower(kw))
	}

	for _, recs := range administeredMap {
		if len(recs) == 0 {
			continue
		}
		// Match keyword trên RAW NAME của mũi đầu tiên trong group (thường cùng loại)
		// Hoặc tốt hơn là check từng record nếu group trộn lẫn (nhưng Flu group thường parse theo keyword)
		rawName := strings.ToLower(recs[0].VaccineName)
		match := false
		for _, kw := range keywords {
			if strings.Contains(rawName, kw) {
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
	displayName := rule.GroupDisplayName
	monthsNow, _, _ := GetAgeAtDate(e.DOB, e.AnalysisDate)

	if monthsNow < rule.MinAgeMonthsAtFirstDose {
		return nil
	}

	if numDoses == 0 {
		statusMsg, earliestDate, ageTags := e.getAgeStatusAndEarliestDate(rule, "")
		return []AnalysisResult{{
			VaccineNameForPopup:  displayName,
			Description:          fmt.Sprintf("%s (Chưa tiêm - cần tiêm mũi 1). %s", displayName, statusMsg),
			EarliestNextDoseDate: earliestDate,
			StatusTags:           ageTags,
			IsMissing:            true,
			DoseNumber:           1,
		}}
	}

	// Đã tiêm ít nhất 1 mũi. Check validity mũi 1.
	valid, errRes := e.checkFirstDoseAgeValidity(fluRecords[0].Date, rule, displayName)
	if !valid {
		return []AnalysisResult{*errRes}
	}

	var results []AnalysisResult
	firstDoseDate := fluRecords[0].Date
	monthsAtFirst, _, _ := GetAgeAtDate(e.DOB, firstDoseDate)
	
	// Trẻ < 9 tuổi lần đầu tiêm cúm cần 2 mũi cách nhau 4 tuần
	needsInitialSecondDose := monthsAtFirst < 108 // 9 years = 108 months

	if needsInitialSecondDose && numDoses == 1 {
		interval := rule.InitialSeriesIntervalDays
		if interval == 0 {
			interval = 28 // Default 4 weeks
		}
		nextDoseDate := firstDoseDate.AddDate(0, 0, interval)
		tags := []string{"due", "flu_second_dose"}
		if e.AnalysisDate.Before(nextDoseDate) {
			tags = []string{"too_young"}
		} else {
			nextDoseDate = e.AnalysisDate
		}
		
		results = append(results, AnalysisResult{
			VaccineNameForPopup:  displayName,
			Description:          fmt.Sprintf("%s (Cần mũi 2 khởi đầu do lần đầu tiêm lúc < 9 tuổi, cách mũi 1 khoảng 4 tuần).", displayName),
			EarliestNextDoseDate: &nextDoseDate,
			StatusTags:           tags,
			DoseNumber:           2,
		})
		return results // Nếu đang thiếu mũi 2 thì không gợi ý annual booster vội
	}

	// Nhắc lại hàng năm (Booster)
	lastDoseDate := fluRecords[numDoses-1].Date
	nextBoosterDate := lastDoseDate.AddDate(1, 0, 0) // +1 year
	
	tags := []string{"due"}
	if e.AnalysisDate.Before(nextBoosterDate) {
		tags = []string{"info", "booster_upcoming"}
	} else {
		tags = []string{"due", "flu_annual"}
	}

	desc := ""
	if containsStr(tags, "booster_due") {
		desc = fmt.Sprintf("%s (Cần tiêm mũi nhắc lại hàng năm - cách mũi cuối > 1 năm).", displayName)
	} else {
		desc = fmt.Sprintf("%s (Đã tiêm. Mũi nhắc tiếp theo vào tháng %s).", displayName, nextBoosterDate.Format("01/2006"))
	}

	results = append(results, AnalysisResult{
		VaccineNameForPopup:  displayName,
		Description:          desc,
		EarliestNextDoseDate: &nextBoosterDate,
		StatusTags:           tags,
	})

	return results
}

// checkMeningococcalACYWGroup xử lý nhóm vắc xin Não mô cầu ACYW và các tương tác.
func (e *Engine) checkMeningococcalACYWGroup(ruleKey string, rule VaccineRule, administeredMap map[string][]models.VaccineRecord) []AnalysisResult {
	// 1. Phân loại bản ghi
	var menactraRecords []models.VaccineRecord
	var menquadfiRecords []models.VaccineRecord

	menactraMember, hasMenactra := rule.Members["MENACTRA"]
	menquadfiMember, hasMenquadfi := rule.Members["MENQUADFI"]

	if hasMenactra {
		menactraRecords = e.getMatchingRecords(menactraMember.NamesNorm, administeredMap)
	}
	if hasMenquadfi {
		menquadfiRecords = e.getMatchingRecords(menquadfiMember.NamesNorm, administeredMap)
	}

	allRecords := append(menactraRecords, menquadfiRecords...)
	sort.Slice(allRecords, func(i, j int) bool {
		return allRecords[i].Date.Before(allRecords[j].Date)
	})

	numDoses := len(allRecords)
	var results []AnalysisResult

	// Helper function for interactions
	applyInteractionsAndAppend := func(res *AnalysisResult, currentRecords []models.VaccineRecord) {
		// Menin ACYW tương tác với VA-MENGOC-BC và 6in1 (Infanrix Hexa/Hexaxim)
		for otherRuleKey, inter := range rule.Interactions {
			otherRule, ok := e.Rules[otherRuleKey]
			if !ok {
				continue
			}
			otherRecords := e.getMatchingRecords(otherRule.NamesNorm, administeredMap)
			if len(otherRecords) == 0 {
				continue
			}

			lastOtherDate := otherRecords[len(otherRecords)-1].Date
			monthsNow, _, _ := GetAgeAtDate(e.DOB, e.AnalysisDate)

			// Only apply interaction if age constraint is met
			if inter.AppliesWhenAgeMonthsGte > 0 && monthsNow < inter.AppliesWhenAgeMonthsGte {
				continue
			}

			if inter.Direction == "reverse" {
				// Cảnh báo ngược (ví dụ MenQuadfi ảnh hưởng đến MenBC)
				tagName := "interaction_" + strings.ToLower(otherRuleKey)
				if otherRuleKey == "VA-MENGOC-BC" {
					tagName = "interaction_reverse_mengoc" // Already correct in parity
				}
				
				results = append(results, AnalysisResult{
					VaccineNameForPopup:  res.VaccineNameForPopup,
					Description:          inter.Message,
					EarliestNextDoseDate: nil,
					StatusTags:           []string{inter.Severity, tagName},
				})
			} else {
				earliestByInteraction := lastOtherDate.AddDate(0, 0, inter.MinIntervalDays)
				if inter.MinIntervalDays == 60 {
					earliestByInteraction = AddMonths(lastOtherDate, 2)
				}
				
				if res.EarliestNextDoseDate == nil || earliestByInteraction.After(*res.EarliestNextDoseDate) {
					res.EarliestNextDoseDate = &earliestByInteraction
				}

				if e.AnalysisDate.Before(earliestByInteraction) {
					tagName := "interaction_constraint"
					if otherRuleKey == "VA-MENGOC-BC" {
						tagName = "interaction_mengoc_bc"
					} else if otherRuleKey == "Six_In_One_Combined" {
						tagName = "interaction_6in1"
					}
					
					// Thêm item cảnh báo tương tác nếu chưa đến ngày tiêm an toàn
					results = append(results, AnalysisResult{
						VaccineNameForPopup: res.VaccineNameForPopup,
						Description:         inter.Message,
						StatusTags:          []string{"warning", tagName},
					})
				}
			}
		}
	}

	// 2. Case: Chưa tiêm mũi nào
	if numDoses == 0 {
		monthsNow, _, _ := GetAgeAtDate(e.DOB, e.AnalysisDate)
		
		// Gợi ý MenQuadfi (ưu tiên, từ 6 tuần)
		if hasMenquadfi {
			mqEarliest := e.DOB.AddDate(0, 0, menquadfiMember.MinAgeWeeksOverall*7)
			mqRes := AnalysisResult{
				VaccineNameForPopup:  menquadfiMember.Display,
				Description:          fmt.Sprintf("%s (Gợi ý ưu tiên - tiêm được từ sớm).", menquadfiMember.Display),
				EarliestNextDoseDate: &mqEarliest,
				StatusTags:           []string{"due"},
				IsMissing:            true,
				DoseNumber:           1,
			}
			if e.AnalysisDate.Before(mqEarliest) {
				mqRes.StatusTags = []string{"too_young"}
			}
			applyInteractionsAndAppend(&mqRes, nil)
			results = append(results, mqRes)
		}

		// Gợi ý Menactra (nếu đủ 9 tháng)
		if hasMenactra && monthsNow >= 9 {
			maEarliest := AddMonths(e.DOB, 9)
			maRes := AnalysisResult{
				VaccineNameForPopup:  menactraMember.Display,
				Description:          fmt.Sprintf("%s (Gợi ý thêm - cho trẻ từ 9 tháng).", menactraMember.Display),
				EarliestNextDoseDate: &maEarliest,
				StatusTags:           []string{"due"},
				IsMissing:            true,
				DoseNumber:           1,
			}
			applyInteractionsAndAppend(&maRes, nil)
			results = append(results, maRes)
		}
		return results
	}

	// 3. Case: Đã tiêm Menactra (hoặc mũi 1 là Menactra)
	// Kiểm tra xem mũi 1 là loại nào
	firstDose := allRecords[0]
	isMenactraFirst := false
	for _, norm := range menactraMember.NamesNorm {
		if NormalizeVaccineName(firstDose.VaccineName) == norm {
			isMenactraFirst = true
			break
		}
	}

	if isMenactraFirst {
		// Delegate Menactra to checkAgeDependentSeries logic
		// Tạo 1 rule giả để gọi
		maRule := VaccineRule{
			DisplayName: menactraMember.Display,
			RulesByAge:  menactraMember.RulesByAge,
			NamesNorm:   menactraMember.NamesNorm,
		}
		maResults := e.checkAgeDependentSeries("MENACTRA", maRule, administeredMap)
		for i := range maResults {
			applyInteractionsAndAppend(&maResults[i], allRecords)
		}
		results = append(results, maResults...)
		return results
	}

	// 4. Case: Đã tiêm MenQuadfi
	// Tìm AgeRule phù hợp cho MenQuadfi
	monthsAtFirst, _, _ := GetAgeAtDate(e.DOB, allRecords[0].Date)
	var applicableMQRule *AgeRule
	for i := range menquadfiMember.RulesByAge {
		ar := &menquadfiMember.RulesByAge[i]
		if (ar.MinAgeMonthsAtFirstDose == 0 || monthsAtFirst >= ar.MinAgeMonthsAtFirstDose) &&
		   (ar.MaxAgeMonthsAtFirstDose == 0 || monthsAtFirst <= ar.MaxAgeMonthsAtFirstDose) &&
		   (ar.MinAgeWeeksAtFirstDose == 0 || int(allRecords[0].Date.Sub(e.DOB).Hours()/24/7) >= ar.MinAgeWeeksAtFirstDose) {
			applicableMQRule = ar
			break
		}
	}

	if applicableMQRule == nil {
		return []AnalysisResult{{
			VaccineNameForPopup: menquadfiMember.Display,
			Description:         fmt.Sprintf("%s - Không tìm thấy phác đồ phù hợp.", menquadfiMember.Display),
			StatusTags:          []string{"no_applicable_rule"},
		}}
	}

	// Nếu chưa đủ liều cơ bản
	if numDoses < applicableMQRule.DosesRequired {
		lastDoseDate := allRecords[numDoses-1].Date
		var nextDoseDate time.Time
		intervalDays := 0
		if numDoses < len(applicableMQRule.MinIntervalDays) && applicableMQRule.MinIntervalDays[numDoses] != nil {
			intervalDays = *applicableMQRule.MinIntervalDays[numDoses]
		}
		
		// Python fix: 60 days = 2 months
		if intervalDays == 60 {
			nextDoseDate = AddMonths(lastDoseDate, 2)
		} else {
			nextDoseDate = lastDoseDate.AddDate(0, 0, intervalDays)
		}

		res := AnalysisResult{
			VaccineNameForPopup:  menquadfiMember.Display,
			Description:          fmt.Sprintf("%s (Tiêm mũi %d).", menquadfiMember.Display, numDoses+1),
			EarliestNextDoseDate: &nextDoseDate,
			StatusTags:           []string{"due"},
			DoseNumber:           numDoses + 1,
		}
		if e.AnalysisDate.Before(nextDoseDate) {
			res.StatusTags = []string{"too_young"}
		}
		applyInteractionsAndAppend(&res, allRecords)
		results = append(results, res)
		return results
	}

	// Nếu đã đủ liều cơ bản -> Check Booster
	if applicableMQRule.Booster != nil {
		lastDoseDate := allRecords[numDoses-1].Date
		earliestByInterval := lastDoseDate.AddDate(0, 0, applicableMQRule.Booster.MinIntervalDaysFromLast)
		earliestByAge := AddMonths(e.DOB, applicableMQRule.Booster.MinAgeMonths)
		
		finalEarliest := earliestByInterval
		if earliestByAge.After(finalEarliest) {
			finalEarliest = earliestByAge
		}

		res := AnalysisResult{
			VaccineNameForPopup:  menquadfiMember.Display,
			Description:          fmt.Sprintf("%s (%s).", menquadfiMember.Display, applicableMQRule.Booster.Description),
			EarliestNextDoseDate: &finalEarliest,
			StatusTags:           []string{"booster_due"},
		}
		
		if e.AnalysisDate.Before(finalEarliest) {
			res.StatusTags = []string{"booster_upcoming", "completed"}
		}
		
		applyInteractionsAndAppend(&res, allRecords)
		results = append(results, res)
		return results
	}

	return nil
}
