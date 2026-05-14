package utils

import (
	"regexp"
	"sort"
	"strings"
	"time"
	"tracuutiemchung-engine/models"
)

// AddMonths adds months to a date, pinning to the last day of the month if the day overflows.
// (e.g., 31/01 + 1 month -> 28/02 or 29/02)
func AddMonths(t time.Time, months int) time.Time {
	y, m, d := t.Date()
	
	// Create a new date with the same day but different month
	// Go's time.Date will normalize if the day is too large for the month (e.g., Feb 31 -> Mar 3)
	res := time.Date(y, m+time.Month(months), d, t.Hour(), t.Minute(), t.Second(), t.Nanosecond(), t.Location())
	
	// If the day is different, it means we overflowed the month
	if res.Day() != d {
		// Go to the last day of the intended month by taking the 0th day of the NEXT month
		res = time.Date(y, m+time.Month(months+1), 0, t.Hour(), t.Minute(), t.Second(), t.Nanosecond(), t.Location())
	}
	return res
}

// AddYears adds years to a date, handling February 29th by pinning to February 28th in non-leap years.
func AddYears(t time.Time, years int) time.Time {
	y, m, d := t.Date()
	res := time.Date(y+years, m, d, t.Hour(), t.Minute(), t.Second(), t.Nanosecond(), t.Location())
	
	// If the month changed (only happens for Feb 29th in non-leap year), pin to last day of Feb
	if res.Month() != m {
		res = time.Date(y+years, m+1, 0, t.Hour(), t.Minute(), t.Second(), t.Nanosecond(), t.Location())
	}
	return res
}

// GetAgeAtDate calculates age in months, weeks, and years.
// Parity with Python engine's get_age_at_date.
func GetAgeAtDate(dob, target time.Time) (months, weeks, years int) {
	if target.Before(dob) {
		return 0, 0, 0
	}

	// Calculate total full years
	years = target.Year() - dob.Year()
	if target.Month() < dob.Month() || (target.Month() == dob.Month() && target.Day() < dob.Day()) {
		years--
	}
	if years < 0 {
		years = 0
	}

	// Calculate total full months
	months = (target.Year()-dob.Year())*12 + int(target.Month()) - int(dob.Month())
	if target.Day() < dob.Day() {
		months--
	}
	if months < 0 {
		months = 0
	}

	// Calculate total weeks (total days / 7)
	d1 := time.Date(dob.Year(), dob.Month(), dob.Day(), 0, 0, 0, 0, time.UTC)
	d2 := time.Date(target.Year(), target.Month(), target.Day(), 0, 0, 0, 0, time.UTC)
	totalDays := int(d2.Sub(d1).Hours() / 24)
	weeks = totalDays / 7

	return months, weeks, years
}

var (
	mcgRegex    = regexp.MustCompile(`(?i)\s*\d+mcg/\d+(\.\d+)?ml\s*$`)
	parenRegex  = regexp.MustCompile(`\s*\(.*?\)\s*`)
	yearRegex   = regexp.MustCompile(`\s+\d{4}/\d{4}\s*$`)
	yearXXRegex = regexp.MustCompile(`\s+20XX/20XX\s*$`)
	mlRegex     = regexp.MustCompile(`(?i)\s+\d+(\.\d+)?ml\s*$`)
	spaceRegex  = regexp.MustCompile(`\s+`)
)

// NormalizeVaccineName normalizes vaccine names by removing diacritics, suffixes, and collapsing spaces.
func NormalizeVaccineName(name string) string {
	// 1. Remove diacritics (accents)
	name = removeDiacritics(name)

	// 2. Pre-normalization steps (parity with Python)
	name = strings.ReplaceAll(name, ",", ".")
	name = mcgRegex.ReplaceAllString(name, "")
	name = parenRegex.ReplaceAllString(name, "")
	name = yearRegex.ReplaceAllString(name, "")
	name = yearXXRegex.ReplaceAllString(name, "")
	name = mlRegex.ReplaceAllString(name, "")

	// 3. Lowercase, trim, and collapse spaces
	name = strings.ToLower(name)
	name = spaceRegex.ReplaceAllString(name, " ")
	return strings.TrimSpace(name)
}

func removeDiacritics(s string) string {
	replacer := strings.NewReplacer(
		"à", "a", "á", "a", "ạ", "a", "ả", "a", "ã", "a", "â", "a", "ầ", "a", "ấ", "a", "ậ", "a", "ẩ", "a", "ẫ", "a", "ă", "a", "ằ", "a", "ắ", "a", "ặ", "a", "ẳ", "a", "ẵ", "a",
		"è", "e", "é", "e", "ẹ", "e", "ẻ", "e", "ẽ", "e", "ê", "e", "ề", "e", "ế", "e", "ệ", "e", "ể", "e", "ễ", "e",
		"ì", "i", "í", "i", "ị", "i", "ỉ", "i", "ĩ", "i",
		"ò", "o", "ó", "o", "ọ", "o", "ỏ", "o", "õ", "o", "ô", "o", "ồ", "o", "ố", "o", "ộ", "o", "ổ", "o", "ỗ", "o", "ơ", "o", "ờ", "o", "ớ", "o", "ợ", "o", "ở", "o", "ỡ", "o",
		"ù", "u", "ú", "u", "ụ", "u", "ủ", "u", "ũ", "u", "ư", "u", "ừ", "u", "ứ", "u", "ự", "u", "ử", "u", "ữ", "u",
		"ỳ", "y", "ý", "y", "ỵ", "y", "ỷ", "y", "ỹ", "y",
		"đ", "d",
		"À", "a", "Á", "a", "Ạ", "a", "Ả", "a", "Ã", "a", "Â", "a", "Ầ", "a", "Ấ", "a", "Ậ", "a", "Ẩ", "a", "Ẫ", "a", "Ă", "a", "Ằ", "a", "Á", "a", "Ạ", "a", "Ả", "a", "Ã", "a",
		"È", "e", "É", "e", "Ẹ", "e", "Ẻ", "e", "Ẽ", "e", "Ê", "e", "Ề", "e", "Ế", "e", "Ệ", "e", "Ể", "e", "Ễ", "e",
		"Ì", "i", "Í", "i", "Ị", "i", "Ỉ", "i", "Ĩ", "i",
		"Ò", "o", "Ó", "o", "Ọ", "o", "Ỏ", "o", "Õ", "o", "Ô", "o", "Ồ", "o", "Ố", "o", "Ộ", "o", "Ổ", "o", "Ỗ", "o", "Ơ", "o", "Ờ", "o", "Ớ", "o", "Ợ", "o", "Ở", "o", "Ỡ", "o",
		"Ù", "u", "Ú", "u", "Ụ", "u", "Ủ", "u", "Ũ", "u", "Ư", "u", "Ừ", "u", "Ứ", "u", "Ự", "u", "Ử", "u", "Ữ", "u",
		"Ỳ", "y", "Ý", "y", "Ỵ", "y", "Ỷ", "y", "Ỹ", "y",
		"Đ", "d",
	)
	return replacer.Replace(s)
}

// SortDoses sorts doses by date in ascending order.
func SortDoses(doses []models.AdministeredDose) {
	sort.Slice(doses, func(i, j int) bool {
		return doses[i].Date.Before(doses[j].Date)
	})
}

// GetDosesForRule collects all doses associated with a rule from various possible fields.
func GetDosesForRule(rule *models.VaccineRule, adminMap map[string][]models.AdministeredDose) []models.AdministeredDose {
	var allDoses []models.AdministeredDose
	seenDates := make(map[int64]bool)
	
	names := rule.NamesNorm
	names = append(names, rule.NamesNormGroup...)
	
	for _, member := range rule.Members {
		names = append(names, member.NamesNorm...)
	}
	for _, course := range rule.Courses {
		names = append(names, course.NamesNorm...)
	}

	for _, name := range names {
		if ds, ok := adminMap[name]; ok {
			for _, d := range ds {
				if !seenDates[d.Date.Unix()] {
					allDoses = append(allDoses, d)
					seenDates[d.Date.Unix()] = true
				}
			}
		}
	}
	
	SortDoses(allDoses)
	return allDoses
}

