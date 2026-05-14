package models

// VaccineRule represents a single rule definition from vaccine_rules.json
type VaccineRule struct {
	Type                                     string                 `json:"type"`
	DisplayName                              string                 `json:"display_name,omitempty"`
	GroupDisplayName                         string                 `json:"group_display_name,omitempty"`
	IsLive                                   bool                   `json:"is_live,omitempty"`
	RawNames                                 []string               `json:"raw_names,omitempty"`
	RawNamesMembers                          map[string][]string    `json:"raw_names_members,omitempty"`
	DosesRequired                            *int                   `json:"doses_required,omitempty"`
	MinIntervalDays                          []*int                 `json:"min_interval_days,omitempty"`
	MinAgeMonthsAtFirstDose                  *int                   `json:"min_age_months_at_first_dose,omitempty"`
	MinAgeWeeksAtFirstDose                   *int                   `json:"min_age_weeks_at_first_dose,omitempty"`
	MinAgeDaysAtFirstDose                    *int                   `json:"min_age_days_at_first_dose,omitempty"`
	MaxAgeMonthsAtFirstDose                  *int                   `json:"max_age_months_at_first_dose,omitempty"`
	MinAgeYearsAtFirstDose                   *int                   `json:"min_age_years_at_first_dose,omitempty"`
	MinAgeWeeksOverall                       *int                   `json:"min_age_weeks_overall,omitempty"`
	Regimens                                 []Regimen              `json:"regimens,omitempty"`
	Members                                  map[string]Member      `json:"members,omitempty"`
	Courses                                  []Course               `json:"courses,omitempty"`
	RulesByAge                               []AgeRule              `json:"rules_by_age,omitempty"`
	Interactions                             map[string]Interaction `json:"interactions,omitempty"`
	RecognitionKeywords                      []string               `json:"recognition_keywords,omitempty"`
	InitialSeriesIntervalDays                *int                   `json:"initial_series_interval_days,omitempty"`
	MaxAgeMonthsToStartFirstDoseGroup        *int                   `json:"max_age_months_to_start_first_dose_group,omitempty"`
	MaxAgeMonthsForCompletionGroup           *int                   `json:"max_age_months_for_completion_group,omitempty"`
	MinAgeMonthsOverallGroup                 *int                   `json:"min_age_months_overall_group,omitempty"`
	ProvidesMeaslesProtectionGroup           bool                   `json:"provides_measles_protection_group,omitempty"`
	BoosterIntervalYears                    *int                   `json:"booster_interval_years,omitempty"`
	NamesNorm                      []string               `json:"names_norm,omitempty"`
	NamesNormGroup                 []string               `json:"names_norm_group,omitempty"`
}

// Regimen defines a specific vaccination schedule based on age
type Regimen struct {
	RegimenName             string                  `json:"regimen_name"`
	MinAgeAtFirstDoseMonths *int                    `json:"min_age_at_first_dose_months,omitempty"`
	MaxAgeAtFirstDoseMonths *int                    `json:"max_age_at_first_dose_months,omitempty"`
	DosesRequired           int                     `json:"doses_required"`
	MinIntervalDays         []*int                  `json:"min_interval_days"`
	DoseSpecificRules       map[string]SpecificRule `json:"dose_specific_rules,omitempty"`
}

// Member represents a specific vaccine brand or subgroup in a rule
type Member struct {
	RawNames            []string  `json:"raw_names,omitempty"`
	Display             string    `json:"display,omitempty"`
	MinAgeMonthsOverall *int      `json:"min_age_months_overall,omitempty"`
	MinAgeWeeksOverall  *int      `json:"min_age_weeks_overall,omitempty"`
	RulesByAge          []AgeRule `json:"rules_by_age,omitempty"`
	NamesNorm           []string  `json:"names_norm,omitempty"`
}

// Course represents an alternative vaccination course
type Course struct {
	RawNames                []string `json:"raw_names"`
	DosesRequired           int      `json:"doses_required"`
	Display                 string   `json:"display"`
	MinAgeMonthsAtFirstDose *int     `json:"min_age_months_at_first_dose,omitempty"`
	MaxAgeYearsAtFirstDose  *int     `json:"max_age_years_at_first_dose,omitempty"`
	MinIntervalDays         []*int   `json:"min_interval_days"`
	IsLive                  bool     `json:"is_live,omitempty"`
	BoosterIntervalYears    *int     `json:"booster_interval_years,omitempty"`
	BoosterAfterDoseNumber  *int     `json:"booster_after_dose_number,omitempty"`
	BoosterMaxAgeYears      *int     `json:"booster_max_age_years,omitempty"`
	NamesNorm               []string `json:"names_norm,omitempty"`
}

// AgeRule defines rules for a specific age range
type AgeRule struct {
	MinAgeAtFirstDoseMonths *int                    `json:"min_age_at_first_dose_months,omitempty"`
	MaxAgeAtFirstDoseMonths *int                    `json:"max_age_at_first_dose_months,omitempty"`
	MinAgeWeeksAtFirstDose  *int                    `json:"min_age_weeks_at_first_dose,omitempty"`
	DosesRequired           int                     `json:"doses_required"`
	MinIntervalDays         []*int                  `json:"min_interval_days"`
	Booster                 *BoosterRule            `json:"booster,omitempty"`
	DoseSpecificRules       map[string]SpecificRule `json:"dose_specific_rules,omitempty"`
}

// BoosterRule defines rules for a booster dose
type BoosterRule struct {
	MinAgeMonths            int    `json:"min_age_months"`
	MinIntervalDaysFromLast int    `json:"min_interval_days_from_last"`
	Description             string `json:"description"`
}

// SpecificRule defines rules for a specific dose number
type SpecificRule struct {
	AlternativeMinAgeYears *int `json:"alternative_min_age_years,omitempty"`
	AlternativeMaxAgeYears *int `json:"alternative_max_age_years,omitempty"`
	MinAbsoluteAgeMonths   *int `json:"min_absolute_age_months,omitempty"`
}

// Interaction defines rules between different vaccines
type Interaction struct {
	MinIntervalDays         int    `json:"min_interval_days"`
	AppliesWhenAgeMonthsGte *int   `json:"applies_when_age_months_gte,omitempty"`
	Direction               string `json:"direction"`
	Severity                string `json:"severity"`
	Message                 string `json:"message"`
}
