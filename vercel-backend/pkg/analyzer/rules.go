package analyzer

import "time"

// RuleType constants
const (
	RuleTypeSingleSeries               = "single_series"
	RuleTypeSingleDoseMinAge           = "single_dose_min_age"
	RuleTypeSingleSeriesMinAge         = "single_series_min_age"
	RuleTypeAgeDependent               = "age_dependent_series"
	RuleTypeGroupCumulativeUnique      = "group_cumulative_unique_doses"
	RuleTypeGroupAlternative           = "group_alternative_courses"
	RuleTypeGroupAlternativeMinAge     = "group_alternative_courses_min_age"
	RuleTypeGroupAlternativeAgeRange   = "group_alternative_courses_age_range"
	RuleTypeFluGroup                   = "flu_group"
	RuleTypeMMREquivalentGroup         = "mmr_equivalent_group"
	RuleTypeMeningococcalACYWGroup     = "meningococcal_acyw_group"
)

// VaccineRule represents a rule for a specific vaccine or group
type VaccineRule struct {
	DisplayName              string           `json:"display_name"`
	GroupDisplayName         string           `json:"group_display_name"`
	Type                     string           `json:"type"`
	RawNames                 []string         `json:"raw_names"`
	RawNamesMembers          map[string][]string `json:"raw_names_members"`
	IsLive                   bool             `json:"is_live"`
	DosesRequired            int              `json:"doses_required"`
	MinIntervalDays          []*int           `json:"min_interval_days"`
	MinAgeMonthsAtFirstDose  int              `json:"min_age_months_at_first_dose"`
	MinAgeWeeksAtFirstDose   int              `json:"min_age_weeks_at_first_dose"`
	MinAgeYearsAtFirstDose   int              `json:"min_age_years_at_first_dose"`
	MinAgeDaysAtFirstDose    int              `json:"min_age_days_at_first_dose"`
	MaxAgeMonthsAtFirstDose  int              `json:"max_age_months_at_first_dose"`
	MaxAgeMonthsToStartFirstDoseGroup int       `json:"max_age_months_to_start_first_dose_group"`
	MaxAgeMonthsForCompletionGroup    int       `json:"max_age_months_for_completion_group"`
	MinAgeMonthsOverallGroup int              `json:"min_age_months_overall_group"`
	MinAgeWeeksOverallGroup  int              `json:"min_age_weeks_overall_group"`
	MinAgeYearsOverallGroup  int              `json:"min_age_years_overall_group"`
	MinAgeDaysOverallGroup   int              `json:"min_age_days_overall_group"`
	MinAgeWeeksOverall       int              `json:"min_age_weeks_overall"`
	MinAgeYearsOverall       int              `json:"min_age_years_overall"`
	MinAgeDaysOverall        int              `json:"min_age_days_overall"`
	MinAgeMonthsOverall      int              `json:"min_age_months_overall"`
	RulesByAge               []AgeRule        `json:"rules_by_age"`
	Courses                  []Course         `json:"courses"`
	Regimens                 []Course         `json:"regimens"`
	Members                  map[string]Member `json:"members"`
	
	// Specialized fields
	RecognitionKeywords      []string         `json:"recognition_keywords"`
	InitialSeriesIntervalDays int              `json:"initial_series_interval_days"`
	Interactions             map[string]Interaction `json:"interactions"`
	
	// Booster fields
	BoosterIntervalYears    int              `json:"booster_interval_years"`
	BoosterAfterDoseNumber  int              `json:"booster_after_dose_number"`
	BoosterMaxAgeYears      int              `json:"booster_max_age_years"`
	
	// Interaction flags
	ProvidesMeaslesProtection bool           `json:"provides_measles_protection"`
	
	// Dose specific rules for single_series
	DoseSpecificRules       map[string]DoseRule `json:"dose_specific_rules"`

	// Internal normalization
	NamesNorm      []string `json:"-"`
	NamesNormGroup []string `json:"-"`
}

// AgeRule represents a rule that depends on the age at first dose
type AgeRule struct {
	MinAgeMonthsAtFirstDose int     `json:"min_age_at_first_dose_months"`
	MaxAgeMonthsAtFirstDose int     `json:"max_age_at_first_dose_months"`
	MinAgeWeeksAtFirstDose  int     `json:"min_age_weeks_at_first_dose"`
	DosesRequired           int     `json:"doses_required"`
	MinIntervalDays         []*int  `json:"min_interval_days"`
	DoseSpecificRules       map[string]DoseRule `json:"dose_specific_rules"`
	Booster                 *Booster `json:"booster"`
}

// Course represents an alternative vaccination course
type Course struct {
	RawNames                []string `json:"raw_names"`
	Display                 string   `json:"display"`
	DosesRequired           int      `json:"doses_required"`
	MinAgeMonthsAtFirstDose int      `json:"min_age_months_at_first_dose"`
	MaxAgeMonthsAtFirstDose int      `json:"max_age_at_first_dose_months"`
	MinAgeYearsAtFirstDose  int      `json:"min_age_years_at_first_dose"`
	MaxAgeYearsAtFirstDose  int      `json:"max_age_years_at_first_dose"`
	MinIntervalDays         []*int   `json:"min_interval_days"`
	IsLive                  bool     `json:"is_live"`
	BoosterIntervalYears    int      `json:"booster_interval_years"`
	BoosterAfterDoseNumber  int      `json:"booster_after_dose_number"`
	BoosterMaxAgeYears      int      `json:"booster_max_age_years"`
	DoseSpecificRules       map[string]DoseRule `json:"dose_specific_rules"`
	
	NamesNorm []string `json:"-"`
}

// Member represents a member of a group with its own rules
type Member struct {
	RawNames           []string `json:"raw_names"`
	Display            string   `json:"display"`
	MinAgeMonthsOverall int     `json:"min_age_months_overall"`
	MinAgeWeeksOverall  int     `json:"min_age_weeks_overall"`
	MinAgeYearsOverall  int     `json:"min_age_years_overall"`
	MinAgeDaysOverall   int     `json:"min_age_days_overall"`
	RulesByAge         []AgeRule `json:"rules_by_age"`
	Booster            *Booster  `json:"booster"`
	
	NamesNorm []string `json:"-"`
}

// Booster represents a booster dose rule
type Booster struct {
	MinAgeMonths          int    `json:"min_age_months"`
	MinIntervalDaysFromLast int    `json:"min_interval_days_from_last"`
	Description           string `json:"description"`
}

// DoseRule represents rules for a specific dose number
type DoseRule struct {
	MinAbsoluteAgeMonths int `json:"min_absolute_age_months"`
	AlternativeMinAgeYears int `json:"alternative_min_age_years"`
	AlternativeMaxAgeYears int `json:"alternative_max_age_years"`
}

// Interaction represents an interaction between vaccines
type Interaction struct {
	MinIntervalDays         int    `json:"min_interval_days"`
	AppliesWhenAgeMonthsGte int    `json:"applies_when_age_months_gte"`
	Direction               string `json:"direction"`
	Severity                string `json:"severity"`
	Message                 string `json:"message"`
}

// AnalysisResult represents the outcome of an analysis for a vaccine
type AnalysisResult struct {
	VaccineNameForPopup  string     `json:"vaccine_name_for_popup"`
	Description          string     `json:"description"`
	EarliestNextDoseDate *time.Time  `json:"earliest_next_dose_date"`
	StatusTags           []string   `json:"status_tags"`
	IsMissing            bool       `json:"is_missing"`
	DoseNumber           int        `json:"dose_number"`
	Recommendation       string     `json:"recommendation"`
}
