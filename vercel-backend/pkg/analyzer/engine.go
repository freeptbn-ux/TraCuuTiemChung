package analyzer

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strings"
	"time"

	"vercel-backend/pkg/models"
)

type Engine struct {
	Rules        map[string]VaccineRule
	DOB          time.Time
	AnalysisDate time.Time
}

func NewEngine(rulesPath string, dob, analysisDate time.Time) (*Engine, error) {
	rules, err := LoadRules(rulesPath)
	if err != nil {
		return nil, err
	}
	return &Engine{
		Rules:        rules,
		DOB:          dob,
		AnalysisDate: analysisDate,
	}, nil
}

func NewEngineFromBytes(data []byte, dob, analysisDate time.Time) (*Engine, error) {
	rules, err := LoadRulesFromBytes(data)
	if err != nil {
		return nil, err
	}
	return &Engine{
		Rules:        rules,
		DOB:          dob,
		AnalysisDate: analysisDate,
	}, nil
}

func LoadRules(path string) (map[string]VaccineRule, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return LoadRulesFromBytes(data)
}

func LoadRulesFromBytes(data []byte) (map[string]VaccineRule, error) {
	var rules map[string]VaccineRule
	if err := json.Unmarshal(data, &rules); err != nil {
		return nil, err
	}
	for key, rule := range rules {
		// Fix mangled characters in display names
		rule.DisplayName = fixEncoding(rule.DisplayName)
		rule.GroupDisplayName = fixEncoding(rule.GroupDisplayName)

		rule.NamesNorm = normalizeSlice(rule.RawNames)
		for i := range rule.Courses {
			rule.Courses[i].Display = fixEncoding(rule.Courses[i].Display)
			rule.Courses[i].NamesNorm = normalizeSlice(rule.Courses[i].RawNames)
		}
		for i := range rule.Regimens {
			rule.Regimens[i].Display = fixEncoding(rule.Regimens[i].Display)
			rule.Regimens[i].NamesNorm = normalizeSlice(rule.Regimens[i].RawNames)
		}
		for k, m := range rule.Members {
			m.Display = fixEncoding(m.Display)
			m.NamesNorm = normalizeSlice(m.RawNames)
			rule.Members[k] = m
		}
		
		// Populate NamesNorm from Courses if empty
		if len(rule.NamesNorm) == 0 && len(rule.Courses) > 0 {
			var allNames []string
			for _, c := range rule.Courses {
				allNames = append(allNames, c.RawNames...)
			}
			rule.NamesNorm = normalizeSlice(allNames)
		}
		
		// Populate NamesNorm from Members if still empty
		if len(rule.NamesNorm) == 0 && len(rule.Members) > 0 {
			var allNames []string
			for _, m := range rule.Members {
				allNames = append(allNames, m.RawNames...)
			}
			rule.NamesNorm = normalizeSlice(allNames)
		}

		// Populate NamesNormGroup from RawNamesMembers
		if rule.RawNamesMembers != nil {
			var groupNames []string
			for _, names := range rule.RawNamesMembers {
				groupNames = append(groupNames, names...)
			}
			rule.NamesNormGroup = normalizeSlice(groupNames)
		}

		rules[key] = rule
	}
	return rules, nil
}

func fixEncoding(s string) string {
	s = strings.ReplaceAll(s, "Ph c u", "Phế cầu")
	s = strings.ReplaceAll(s, "B%", "Bỉ")
	s = strings.ReplaceAll(s, "Vit Nam", "Việt Nam")
	s = strings.ReplaceAll(s, "ViAm Gan A", "Viêm Gan A")
	s = strings.ReplaceAll(s, "V_c xin", "Vắc xin")
	s = strings.ReplaceAll(s, "Vắc xin", "Vắc xin") // Ensure standard form
	s = strings.ReplaceAll(s, "CAm", "Cúm")
	s = strings.ReplaceAll(s, "PhAp", "Pháp")
	return s
}

func normalizeSlice(s []string) []string {
	seen := make(map[string]bool)
	var res []string
	for _, item := range s {
		norm := NormalizeVaccineName(item)
		if norm != "" && !seen[norm] {
			res = append(res, norm)
			seen[norm] = true
		}
	}
	return res
}

func (e *Engine) Analyze(administered []models.VaccineRecord) []AnalysisResult {
	administeredMap := e.buildAdministeredMap(administered)
	var results []AnalysisResult

	// Sort rule keys for consistent output
	var ruleKeys []string
	for k := range e.Rules {
		ruleKeys = append(ruleKeys, k)
	}
	sort.Strings(ruleKeys)

	// Collect special group results first
	pneumoResults := e.processPneumoRules(administeredMap)
	results = append(results, pneumoResults...)

	for _, ruleKey := range ruleKeys {
		// Skip rules handled by processPneumoRules
		if ruleKey == "Prevenar13" || ruleKey == "Vaxneuvance" || ruleKey == "Synflorix" || ruleKey == "Pneumovax23" {
			continue
		}

		rule := e.Rules[ruleKey]
		switch rule.Type {
		case RuleTypeSingleSeries, RuleTypeSingleDoseMinAge, RuleTypeSingleSeriesMinAge:
			res := e.checkSingleSeries(ruleKey, rule, administeredMap)
			results = append(results, res...)
		case RuleTypeAgeDependent:
			res := e.checkAgeDependentSeries(ruleKey, rule, administeredMap)
			results = append(results, res...)
		case RuleTypeGroupCumulativeUnique:
			res := e.checkCumulativeGroupDoses(ruleKey, rule, administeredMap)
			results = append(results, res...)
		case RuleTypeGroupAlternative, RuleTypeGroupAlternativeMinAge:
			res := e.checkAlternativeCoursesMinAgeGroup(ruleKey, rule, administeredMap)
			results = append(results, res...)
		case RuleTypeGroupAlternativeAgeRange:
			res := e.checkAlternativeCoursesAgeRangeGroup(ruleKey, rule, administeredMap)
			results = append(results, res...)
		case RuleTypeFluGroup:
			res := e.checkFluGroup(ruleKey, rule, administeredMap)
			results = append(results, res...)
		case RuleTypeMMREquivalentGroup:
			res := e.checkMMREquivalentGroup(ruleKey, rule, administeredMap)
			results = append(results, res...)
		case RuleTypeMeningococcalACYWGroup:
			res := e.checkMeningococcalACYWGroup(ruleKey, rule, administeredMap)
			results = append(results, res...)
		}
	}

	// Post-processing: Spacing and Sorting
	results = e.ApplySpacingAndSort(results, administeredMap)

	return results
}

func (e *Engine) checkSingleSeries(ruleKey string, rule VaccineRule, administeredMap map[string][]models.VaccineRecord) []AnalysisResult {
	records := e.getMatchingRecords(rule.NamesNorm, administeredMap)
	numDoses := len(records)
	displayName := rule.DisplayName
	if displayName == "" {
		displayName = rule.GroupDisplayName
	}

	var results []AnalysisResult

	// 1. Interactions & Coverage Checks
	// MVVAC <-> MMR Interaction
	if ruleKey == "MVVAC" {
		for _, otherRule := range e.Rules {
			if otherRule.ProvidesMeaslesProtection {
				otherRecords := e.getMatchingRecords(otherRule.NamesNorm, administeredMap)
				for _, rec := range otherRecords {
					months, _, _ := GetAgeAtDate(e.DOB, rec.Date)
					if months >= 12 {
						return []AnalysisResult{{
							VaccineNameForPopup: displayName,
							Description:         fmt.Sprintf("%s (Đã được bảo vệ bởi %s tiêm lúc %d tháng)", displayName, otherRule.DisplayName, months),
							StatusTags:          []string{"coverage_by_other", "completed"},
						}}
					}
				}
			}
		}
	}

	// VA-MENGOC-BC Reverse Interaction
	if ruleKey == "VA-MENGOC-BC" {
		months, _, _ := GetAgeAtDate(e.DOB, e.AnalysisDate)
		if months >= 24 {
			results = append(results, AnalysisResult{
				VaccineNameForPopup: displayName,
				Description:         "Lưu ý: Trẻ từ 2 tuổi trở lên nên tiêm VA-Mengoc BC SAU khi đã tiêm MenQuadfi/Menactra để đạt hiệu quả bảo vệ tốt nhất.",
				StatusTags:          []string{"warning", "interaction_reverse_mengoc"},
			})
		}
	}

	// 2. Already Completed Case (Check Booster)
	if numDoses >= rule.DosesRequired && rule.DosesRequired > 0 {
		if rule.BoosterIntervalYears > 0 {
			boosterAfterDose := rule.BoosterAfterDoseNumber
			if boosterAfterDose == 0 {
				boosterAfterDose = rule.DosesRequired
			}

			if numDoses >= boosterAfterDose {
				lastDoseDate := records[numDoses-1].Date
				nextBoosterDate := AddYears(lastDoseDate, rule.BoosterIntervalYears)

				if rule.BoosterMaxAgeYears > 0 {
					_, _, yearsNow := GetAgeAtDate(e.DOB, e.AnalysisDate)
					if yearsNow >= rule.BoosterMaxAgeYears {
						return results
					}
					_, _, ageAtBoosterYears := GetAgeAtDate(e.DOB, nextBoosterDate)
					if ageAtBoosterYears >= rule.BoosterMaxAgeYears {
						return results
					}
				}

				tags := []string{"due", "booster_due"}
				if e.AnalysisDate.Before(nextBoosterDate) {
					tags = []string{"info", "booster_upcoming"}
				}

				results = append(results, AnalysisResult{
					VaccineNameForPopup:  displayName,
					Description:          fmt.Sprintf("%s - Cần tiêm mũi nhắc lại định kỳ %d năm.", displayName, rule.BoosterIntervalYears),
					EarliestNextDoseDate: &nextBoosterDate,
					StatusTags:           tags,
				})
				return results
			}
		}
		return results
	}

	// 3. First Dose Logic (0 doses)
	if numDoses == 0 {
		statusMsg, earliestDate, ageTags := e.getAgeStatusAndEarliestDate(rule, "")
		desc := fmt.Sprintf("%s (Chưa tiêm - cần %d liều). %s", displayName, rule.DosesRequired, statusMsg)

		results = append(results, AnalysisResult{
			VaccineNameForPopup:  displayName,
			Description:          desc,
			EarliestNextDoseDate: earliestDate,
			StatusTags:           ageTags,
			IsMissing:            true,
			DoseNumber:           1,
		})
		return results
	}

	// 4. Subsequent Dose Logic
	// Validate first dose age
	valid, errRes := e.checkFirstDoseAgeValidity(records[0].Date, rule, displayName)
	if !valid {
		results = append(results, *errRes)
		return results
	}

	// Calculate earliest next dose
	lastDoseDate := records[numDoses-1].Date
	
	// Date by interval
	intervalDays := 0
	if numDoses < len(rule.MinIntervalDays) && rule.MinIntervalDays[numDoses] != nil {
		intervalDays = *rule.MinIntervalDays[numDoses]
	}
	earliestDate := lastDoseDate.AddDate(0, 0, intervalDays)
	
	// Dose specific rules (e.g. min age for dose 4, alternative min age for MMR)
	doseNum := numDoses + 1
	if rule.DoseSpecificRules != nil {
		if dr, ok := rule.DoseSpecificRules[fmt.Sprintf("%d", doseNum)]; ok {
			if dr.MinAbsoluteAgeMonths > 0 {
				absMinDate := AddMonths(e.DOB, dr.MinAbsoluteAgeMonths)
				if earliestDate.Before(absMinDate) {
					earliestDate = absMinDate
				}
			}
			if dr.AlternativeMinAgeYears > 0 {
				altMinDate := AddYears(e.DOB, dr.AlternativeMinAgeYears)
				if earliestDate.After(altMinDate) {
					// Logic Python: Thường dùng date_by_interval, nhưng nếu có alt_min_age thì có thể lấy cái nào nhỏ hơn?
					// Thực tế Python: earliest_next_dose_date = effective_alt_date if alt_date < earliest else earliest
					// Xem series.py:202: if earliest_next_dose_date is None or effective_alt_date < earliest_next_dose_date: earliest_next_dose_date = effective_alt_date
					earliestDate = altMinDate
				}
			}
		}
	}

	statusMsg, _, _ := e.getAgeStatusAndEarliestDate(rule, "")
	tags := []string{"due"}
	if e.AnalysisDate.Before(earliestDate) {
		tags = []string{"info", "booster_upcoming"}
		// Python legacy marks regular subsequent doses as "due" even if in the future.
		// Only first dose gets "scheduled" or "too_young".
		if numDoses == 0 {
			tags = []string{"info", "scheduled"}
		} else {
			// For parity: subsequent missing doses stay "due"
			tags = []string{"due"}
		}
	}

	desc := fmt.Sprintf("%s (Đã tiêm %d mũi. Cần tiêm mũi %d). %s", displayName, numDoses, doseNum, statusMsg)
	// Specific description for parity
	if numDoses > 0 {
		remaining := rule.DosesRequired - numDoses
		desc = fmt.Sprintf("%s - Cần thêm %d liều. Mũi %d cách mũi %d tối thiểu %s.", displayName, remaining, doseNum, numDoses, formatIntervalDescription(intervalDays))
	}

	results = append(results, AnalysisResult{
		VaccineNameForPopup:  displayName,
		Description:          desc,
		EarliestNextDoseDate: &earliestDate,
		StatusTags:           tags,
		IsMissing:            true,
		DoseNumber:           doseNum,
	})

	return results
}

func (e *Engine) checkAgeDependentSeries(ruleKey string, rule VaccineRule, administeredMap map[string][]models.VaccineRecord) []AnalysisResult {
	records := e.getMatchingRecords(rule.NamesNorm, administeredMap)
	numDoses := len(records)
	displayName := rule.DisplayName

	// Determine which age rule applies
	var applicableRule *AgeRule
	if numDoses == 0 {
		// Python legacy does not filter by age when recommending the first dose of an age-dependent series.
		// It just uses the first rule as a placeholder for doses_required.
		if len(rule.RulesByAge) > 0 {
			applicableRule = &rule.RulesByAge[0]
		}
		
		if applicableRule == nil {
			return []AnalysisResult{{
				VaccineNameForPopup: displayName,
				Description:         fmt.Sprintf("%s (Không tìm thấy phác đồ phù hợp)", displayName),
				StatusTags:          []string{"no_applicable_rule"},
			}}
		}
	} else {
		// Use age at first dose
		months, _, _ := GetAgeAtDate(e.DOB, records[0].Date)
		for _, ar := range rule.RulesByAge {
			if months >= ar.MinAgeMonthsAtFirstDose && (ar.MaxAgeMonthsAtFirstDose == 0 || months <= ar.MaxAgeMonthsAtFirstDose) {
				applicableRule = &ar
				break
			}
		}

		if applicableRule == nil {
			return []AnalysisResult{{
				VaccineNameForPopup: displayName,
				Description:         fmt.Sprintf("%s (Mũi 1 tiêm lúc %d tháng không khớp phác đồ nào)", displayName, months),
				StatusTags:          []string{"error_invalid_first_dose_age"},
			}}
		}
		
		displayName = fmt.Sprintf("%s (phác đồ cho mũi 1 lúc %d tháng)", displayName, months)
	}

	// Now check current status based on applicableRule
	if numDoses >= applicableRule.DosesRequired {
		// Check booster if defined
		if applicableRule.Booster != nil {
			lastDoseDate := records[numDoses-1].Date
			nextBoosterDate := lastDoseDate.AddDate(0, 0, applicableRule.Booster.MinIntervalDaysFromLast)
			
			// Age check for booster
			minBoosterAgeDate := AddMonths(e.DOB, applicableRule.Booster.MinAgeMonths)
			if nextBoosterDate.Before(minBoosterAgeDate) {
				nextBoosterDate = minBoosterAgeDate
			}

			tags := []string{"due", "booster_due"}
			if e.AnalysisDate.Before(nextBoosterDate) {
				tags = []string{"info", "booster_upcoming"}
			}

			return []AnalysisResult{{
				VaccineNameForPopup:  displayName,
				Description:          fmt.Sprintf("%s - %s", displayName, applicableRule.Booster.Description),
				EarliestNextDoseDate: &nextBoosterDate,
				StatusTags:           tags,
			}}
		}
		return nil
	}

	// Missing doses
	lastDoseDate := e.AnalysisDate // Default if 0 doses
	intervalDays := 0
	if numDoses > 0 {
		lastDoseDate = records[numDoses-1].Date
		if numDoses < len(applicableRule.MinIntervalDays) && applicableRule.MinIntervalDays[numDoses] != nil {
			intervalDays = *applicableRule.MinIntervalDays[numDoses]
		}
	}
	earliestDate := lastDoseDate.AddDate(0, 0, intervalDays)
	
	// Dose specific rules for age-dependent
	doseNum := numDoses + 1
	if applicableRule.DoseSpecificRules != nil {
		if dr, ok := applicableRule.DoseSpecificRules[fmt.Sprintf("%d", doseNum)]; ok {
			if dr.MinAbsoluteAgeMonths > 0 {
				absMinDate := AddMonths(e.DOB, dr.MinAbsoluteAgeMonths)
				if earliestDate.Before(absMinDate) {
					earliestDate = absMinDate
				}
			}
		}
	}

	statusMsg, _, tags := e.getAgeStatusAndEarliestDate(rule, "")
	desc := fmt.Sprintf("%s (Đã tiêm %d mũi. Cần tiêm mũi %d). %s", displayName, numDoses, doseNum, statusMsg)
	if numDoses == 0 {
		desc = fmt.Sprintf("%s (Chưa tiêm - cần %d liều theo phác đồ). %s", displayName, applicableRule.DosesRequired, statusMsg)
	} else {
		// Python legacy always marks subsequent doses as "due", even if age status says "eligible"
		tags = []string{"due"}
	}

	if e.AnalysisDate.Before(earliestDate) {
		if numDoses > 0 {
			tags = []string{"due"}
		} else {
			// For dose 0, if before earliestDate, it's scheduled/too_young which getAgeStatus handled.
			// But wait, getAgeStatusAndEarliestDate already returned too_young/scheduled.
		}
	}

	return []AnalysisResult{{
		VaccineNameForPopup:  displayName,
		Description:          desc,
		EarliestNextDoseDate: &earliestDate,
		StatusTags:           tags,
		IsMissing:            true,
		DoseNumber:           doseNum,
	}}
}

func (e *Engine) buildAdministeredMap(administered []models.VaccineRecord) map[string][]models.VaccineRecord {
	m := make(map[string][]models.VaccineRecord)
	for _, rec := range administered {
		norm := NormalizeVaccineName(rec.VaccineName)
		m[norm] = append(m[norm], rec)
	}
	return m
}

func (e *Engine) getMatchingRecords(namesNorm []string, administeredMap map[string][]models.VaccineRecord) []models.VaccineRecord {
	var records []models.VaccineRecord
	seenNames := make(map[string]bool)
	for _, name := range namesNorm {
		if seenNames[name] {
			continue
		}
		seenNames[name] = true
		if recs, ok := administeredMap[name]; ok {
			records = append(records, recs...)
		}
	}
	// Sort records by date
	sort.Slice(records, func(i, j int) bool {
		return records[i].Date.Before(records[j].Date)
	})
	return records
}

func formatIntervalDescription(days int) string {
	if days == 0 {
		return ""
	}
	if days%365 == 0 {
		return fmt.Sprintf("%d năm", days/365)
	}
	if days%30 == 0 {
		return fmt.Sprintf("%d tháng", days/30)
	}
	return fmt.Sprintf("%d ngày", days)
}
