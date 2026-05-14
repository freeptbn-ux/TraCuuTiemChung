package analyzer

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"time"

	"vercel-backend/pkg/models"
)

// Engine is the main vaccine analysis engine
type Engine struct {
	Rules        map[string]VaccineRule
	AnalysisDate time.Time
	DOB          time.Time
}

// NewEngine creates a new analysis engine with rules loaded from path
func NewEngine(rulesPath string, dob, analysisDate time.Time) (*Engine, error) {
	data, err := os.ReadFile(rulesPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read rules file: %w", err)
	}

	var rawRules map[string]VaccineRule
	if err := json.Unmarshal(data, &rawRules); err != nil {
		return nil, fmt.Errorf("failed to unmarshal rules: %w", err)
	}

	// Normalize names in rules
	for key, rule := range rawRules {
		rule.NamesNorm = make([]string, 0)
		for _, name := range rule.RawNames {
			rule.NamesNorm = append(rule.NamesNorm, NormalizeVaccineName(name))
		}
		
		if rule.Courses != nil {
			for i := range rule.Courses {
				rule.Courses[i].NamesNorm = make([]string, 0)
				for _, name := range rule.Courses[i].RawNames {
					rule.Courses[i].NamesNorm = append(rule.Courses[i].NamesNorm, NormalizeVaccineName(name))
				}
			}
		}

		if rule.Regimens != nil {
			for i := range rule.Regimens {
				rule.Regimens[i].NamesNorm = make([]string, 0)
				for _, name := range rule.Regimens[i].RawNames {
					rule.Regimens[i].NamesNorm = append(rule.Regimens[i].NamesNorm, NormalizeVaccineName(name))
				}
			}
		}

		if rule.Members != nil {
			for k, m := range rule.Members {
				m.NamesNorm = make([]string, 0)
				for _, name := range m.RawNames {
					m.NamesNorm = append(m.NamesNorm, NormalizeVaccineName(name))
				}
				rule.Members[k] = m
			}
		}

		rawRules[key] = rule
	}

	return &Engine{
		Rules:        rawRules,
		DOB:          dob,
		AnalysisDate: analysisDate,
	}, nil
}

// getMatchingRecords collects and sorts records matching the given normalized names
func (e *Engine) getMatchingRecords(namesNorm []string, administeredMap map[string][]models.VaccineRecord) []models.VaccineRecord {
	var records []models.VaccineRecord
	seen := make(map[string]bool)
	for _, name := range namesNorm {
		if seen[name] {
			continue
		}
		seen[name] = true
		if recs, ok := administeredMap[name]; ok {
			records = append(records, recs...)
		}
	}
	sort.Slice(records, func(i, j int) bool {
		return records[i].Date.Before(records[j].Date)
	})
	return records
}

// Analyze processes the vaccination history and returns missing items
func (e *Engine) Analyze(history []models.VaccineRecord) []AnalysisResult {
	administeredMap := make(map[string][]models.VaccineRecord)
	for _, rec := range history {
		norm := NormalizeVaccineName(rec.VaccineName)
		administeredMap[norm] = append(administeredMap[norm], rec)
	}

	for norm := range administeredMap {
		sort.Slice(administeredMap[norm], func(i, j int) bool {
			return administeredMap[norm][i].Date.Before(administeredMap[norm][j].Date)
		})
	}

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
			if res != nil {
				results = append(results, *res)
			}
		case RuleTypeAgeDependent:
			res := e.checkAgeDependentSeries(ruleKey, rule, administeredMap)
			if res != nil {
				results = append(results, *res)
			}
		case RuleTypeGroupAlternativeMinAge:
			res := e.checkAlternativeCoursesMinAgeGroup(ruleKey, rule, administeredMap)
			if res != nil {
				results = append(results, *res)
			}
		case RuleTypeGroupAlternativeAgeRange:
			res := e.checkAlternativeCoursesAgeRangeGroup(ruleKey, rule, administeredMap)
			if res != nil {
				results = append(results, *res)
			}
		case RuleTypeFluGroup:
			res := e.checkFluGroup(ruleKey, rule, administeredMap)
			if res != nil {
				results = append(results, *res)
			}
		case RuleTypeMMREquivalentGroup:
			res := e.checkMMREquivalentGroup(ruleKey, rule, administeredMap)
			if res != nil {
				results = append(results, *res)
			}
		case RuleTypeMeningococcalACYWGroup:
			res := e.checkMeningococcalACYWGroup(ruleKey, rule, administeredMap)
			if res != nil {
				results = append(results, *res)
			}
		}
	}

	return results
}

func (e *Engine) checkSingleSeries(ruleKey string, rule VaccineRule, administeredMap map[string][]models.VaccineRecord) *AnalysisResult {
	records := e.getMatchingRecords(rule.NamesNorm, administeredMap)
	numDoses := len(records)

	if numDoses >= rule.DosesRequired {
		return nil
	}

	displayName := rule.DisplayName
	if displayName == "" {
		displayName = rule.GroupDisplayName
	}

	var earliestDate time.Time
	if numDoses == 0 {
		// Check max age if applicable
		months, _, _ := GetAgeAtDate(e.DOB, e.AnalysisDate)
		if rule.MaxAgeMonthsAtFirstDose > 0 && months > rule.MaxAgeMonthsAtFirstDose {
			return &AnalysisResult{
				VaccineNameForPopup: displayName,
				Description:         fmt.Sprintf("%s (Đã quá tuổi chỉ định tiêm)", displayName),
				StatusTags:          []string{"too_old"},
			}
		}

		// First dose logic
		if rule.MinAgeMonthsAtFirstDose > 0 {
			earliestDate = AddMonths(e.DOB, rule.MinAgeMonthsAtFirstDose)
		} else if rule.MinAgeWeeksAtFirstDose > 0 {
			earliestDate = e.DOB.AddDate(0, 0, rule.MinAgeWeeksAtFirstDose*7)
		} else if rule.MinAgeYearsAtFirstDose > 0 {
			earliestDate = AddYears(e.DOB, rule.MinAgeYearsAtFirstDose)
		} else if rule.MinAgeDaysAtFirstDose > 0 {
			earliestDate = e.DOB.AddDate(0, 0, rule.MinAgeDaysAtFirstDose)
		} else {
			earliestDate = e.DOB
		}
	} else {
		// Subsequent dose logic
		lastDoseDate := records[numDoses-1].Date
		intervalDays := 0
		if numDoses < len(rule.MinIntervalDays) && rule.MinIntervalDays[numDoses] != nil {
			intervalDays = *rule.MinIntervalDays[numDoses]
		}
		earliestDate = lastDoseDate.AddDate(0, 0, intervalDays)
	}

	// Check if already due
	status := "tiêm mũi " + fmt.Sprint(numDoses+1)
	tags := []string{"due"}
	if e.AnalysisDate.Before(earliestDate) {
		tags = []string{"too_young"}
	}

	return &AnalysisResult{
		VaccineNameForPopup:  displayName,
		Description:          fmt.Sprintf("%s (%s)", displayName, status),
		EarliestNextDoseDate: &earliestDate,
		StatusTags:           tags,
	}
}

func (e *Engine) checkAgeDependentSeries(ruleKey string, rule VaccineRule, administeredMap map[string][]models.VaccineRecord) *AnalysisResult {
	records := e.getMatchingRecords(rule.NamesNorm, administeredMap)
	numDoses := len(records)

	// Determine which age rule applies
	var applicableRule *AgeRule
	if numDoses == 0 {
		months, _, _ := GetAgeAtDate(e.DOB, e.AnalysisDate)
		for _, ar := range rule.RulesByAge {
			if months >= ar.MinAgeMonthsAtFirstDose && (ar.MaxAgeMonthsAtFirstDose == 0 || months <= ar.MaxAgeMonthsAtFirstDose) {
				applicableRule = &ar
				break
			}
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
	}

	if applicableRule == nil {
		return nil
	}

	if numDoses >= applicableRule.DosesRequired {
		return nil
	}

	displayName := rule.DisplayName
	
	var earliestDate time.Time
	if numDoses == 0 {
		if applicableRule.MinAgeMonthsAtFirstDose > 0 {
			earliestDate = AddMonths(e.DOB, applicableRule.MinAgeMonthsAtFirstDose)
		} else {
			earliestDate = e.DOB
		}
	} else {
		lastDoseDate := records[numDoses-1].Date
		intervalDays := 0
		if numDoses < len(applicableRule.MinIntervalDays) && applicableRule.MinIntervalDays[numDoses] != nil {
			intervalDays = *applicableRule.MinIntervalDays[numDoses]
		}
		earliestDate = lastDoseDate.AddDate(0, 0, intervalDays)
	}

	status := "tiêm mũi " + fmt.Sprint(numDoses+1)
	tags := []string{"due"}
	if e.AnalysisDate.Before(earliestDate) {
		tags = []string{"too_young"}
	}

	return &AnalysisResult{
		VaccineNameForPopup:  displayName,
		Description:          fmt.Sprintf("%s (%s)", displayName, status),
		EarliestNextDoseDate: &earliestDate,
		StatusTags:           tags,
	}
}
