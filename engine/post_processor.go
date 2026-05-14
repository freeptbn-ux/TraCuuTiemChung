package engine

import (
	"sort"
	"time"
	"tracuutiemchung-engine/engine/utils"
	"tracuutiemchung-engine/models"
)

// ApplySpacingAndSort adjusts dates for live vaccines and sorts the final results.
func ApplySpacingAndSort(results []models.MissingItem, rules map[string]models.VaccineRule, record models.PatientRecord) []models.MissingItem {
	if len(results) == 0 {
		return results
	}

	// 1. Identify live vaccine rules
	isLiveMap := make(map[string]bool)
	for _, rule := range rules {
		if rule.IsLive {
			// Map by DisplayName or GroupDisplayName
			name := rule.DisplayName
			if rule.GroupDisplayName != "" {
				name = rule.GroupDisplayName
			}
			isLiveMap[name] = true
			
			// Also map members if any
			if rule.RawNamesMembers != nil {
				// handled by cumulative group display name
			}
			for _, course := range rule.Courses {
				if course.IsLive {
					isLiveMap[course.Display] = true
				}
			}
		}
	}

	// 2. Collect all administered live vaccine dates
	var adminLiveDates []time.Time
	for _, rule := range rules {
		if rule.IsLive {
			// Check top-level NamesNorm
			for _, name := range rule.NamesNorm {
				if doses, ok := record.AdministeredMap[name]; ok {
					for _, d := range doses {
						adminLiveDates = append(adminLiveDates, d.Date)
					}
				}
			}
			// Check member NamesNorm
			if rule.RawNamesMembers != nil {
				for _, memberNames := range rule.RawNamesMembers {
					for _, raw := range memberNames {
						name := utils.NormalizeVaccineName(raw)
						if doses, ok := record.AdministeredMap[name]; ok {
							for _, d := range doses {
								adminLiveDates = append(adminLiveDates, d.Date)
							}
						}
					}
				}
			}
			for _, member := range rule.Members {
				for _, name := range member.NamesNorm {
					if doses, ok := record.AdministeredMap[name]; ok {
						for _, d := range doses {
							adminLiveDates = append(adminLiveDates, d.Date)
						}
					}
				}
			}
			// Check courses
			for _, course := range rule.Courses {
				if course.IsLive {
					for _, name := range course.NamesNorm {
						if doses, ok := record.AdministeredMap[name]; ok {
							for _, d := range doses {
								adminLiveDates = append(adminLiveDates, d.Date)
							}
						}
					}
				}
			}
		}
	}
	
	// Sort admin live dates
	sort.Slice(adminLiveDates, func(i, j int) bool {
		return adminLiveDates[i].Before(adminLiveDates[j])
	})

	// 3. Sort recommendations by EarliestNextDoseDate
	sort.Slice(results, func(i, j int) bool {
		di := results[i].EarliestNextDoseDate
		dj := results[j].EarliestNextDoseDate
		if di == nil { return false }
		if dj == nil { return true }
		return di.Before(*dj)
	})

	// 4. Adjust spacing for live vaccines
	// We'll track the "last used live date" (administered or recommended)
	var lastLiveDate time.Time
	if len(adminLiveDates) > 0 {
		lastLiveDate = adminLiveDates[len(adminLiveDates)-1]
	}

	for i := range results {
		if results[i].EarliestNextDoseDate == nil {
			continue
		}

		if isLiveMap[results[i].VaccineName] {
			currentDate := *results[i].EarliestNextDoseDate
			
			if !lastLiveDate.IsZero() {
				// If not on the same day as last live dose
				if !isSameDay(currentDate, lastLiveDate) {
					minAllowedDate := lastLiveDate.AddDate(0, 0, 28)
					if currentDate.Before(minAllowedDate) {
						currentDate = minAllowedDate
						results[i].EarliestNextDoseDate = &currentDate
						results[i].StatusTags = append(results[i].StatusTags, "spacing_adjusted")
					}
				}
			}
			
			// Update lastLiveDate for the next item in the sorted list
			// Wait, if multiple live vaccines are recommended on the SAME day, that's allowed.
			// So we only update lastLiveDate if this one is AFTER the current lastLiveDate.
			if currentDate.After(lastLiveDate) {
				lastLiveDate = currentDate
			}
		}
	}

	// 5. Final sort (in case some dates were pushed forward)
	sort.Slice(results, func(i, j int) bool {
		di := results[i].EarliestNextDoseDate
		dj := results[j].EarliestNextDoseDate
		if di == nil { return false }
		if dj == nil { return true }
		if di.Equal(*dj) {
			return results[i].VaccineName < results[j].VaccineName
		}
		return di.Before(*dj)
	})

	return results
}

func isSameDay(t1, t2 time.Time) bool {
	y1, m1, d1 := t1.Date()
	y2, m2, d2 := t2.Date()
	return y1 == y2 && m1 == m2 && d1 == d2
}
