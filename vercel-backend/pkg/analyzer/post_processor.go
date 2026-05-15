package analyzer

import (
	"sort"
	"strings"
	"time"
	"vercel-backend/pkg/models"
)

// ApplySpacingAndSort áp dụng quy tắc khoảng cách giữa các mũi tiêm (spacing) 
// và sắp xếp kết quả theo ngày tiêm sớm nhất.
func (e *Engine) ApplySpacingAndSort(results []AnalysisResult, administeredMap map[string][]models.VaccineRecord) []AnalysisResult {
	if len(results) == 0 {
		return results
	}

	// 1. Tìm ngày tiêm gần nhất của vaccine sống và vaccine bất kỳ
	var lastLiveDate *time.Time
	var lastAnyDate *time.Time

	for _, records := range administeredMap {
		for _, rec := range records {
			normName := NormalizeVaccineName(rec.VaccineName)
			isLive := e.getVaccineLiveStatusByNormName(normName)
			
			if isLive {
				if lastLiveDate == nil || rec.Date.After(*lastLiveDate) {
					d := rec.Date
					lastLiveDate = &d
				}
			}
			
			if lastAnyDate == nil || rec.Date.After(*lastAnyDate) {
				d := rec.Date
				lastAnyDate = &d
			}
		}
	}

	// 2. Áp dụng spacing cho từng kết quả
	for i := range results {
		if results[i].EarliestNextDoseDate == nil {
			continue
		}

		earliest := *results[i].EarliestNextDoseDate
		isMissingLive := e.isMissingItemLive(results[i])

		// Quy tắc 1: Cách mũi tiêm gần nhất tối thiểu 14 ngày
		if lastAnyDate != nil {
			minDateAny := lastAnyDate.AddDate(0, 0, 14)
			if minDateAny.After(earliest) {
				earliest = minDateAny
			}
		}

		// Quy tắc 2: Nếu là vaccine sống và mũi trước cũng là vaccine sống -> cách 28 ngày
		if isMissingLive && lastLiveDate != nil {
			minDateLive := lastLiveDate.AddDate(0, 0, 28)
			if minDateLive.After(earliest) {
				earliest = minDateLive
			}
		}

		// Clamp to AnalysisDate (không lùi ngày về quá khứ so với thời điểm tra cứu)
		if earliest.Before(e.AnalysisDate) {
			earliest = e.AnalysisDate
		}

		results[i].EarliestNextDoseDate = &earliest
	}

	// 3. Sắp xếp kết quả
	sort.Slice(results, func(i, j int) bool {
		// Nil date goes last
		if results[i].EarliestNextDoseDate == nil && results[j].EarliestNextDoseDate != nil {
			return false
		}
		if results[i].EarliestNextDoseDate != nil && results[j].EarliestNextDoseDate == nil {
			return true
		}
		if results[i].EarliestNextDoseDate != nil && results[j].EarliestNextDoseDate != nil {
			if !results[i].EarliestNextDoseDate.Equal(*results[j].EarliestNextDoseDate) {
				return results[i].EarliestNextDoseDate.Before(*results[j].EarliestNextDoseDate)
			}
		}
		
		// Then by description
		return results[i].Description < results[j].Description
	})

	return results
}

// getVaccineLiveStatusByNormName kiểm tra xem một vaccine (theo tên chuẩn hóa) có phải là vaccine sống hay không.
func (e *Engine) getVaccineLiveStatusByNormName(normName string) bool {
	// Kiểm tra trong Rules
	for _, rule := range e.Rules {
		for _, name := range rule.NamesNorm {
			if name == normName {
				return rule.IsLive
			}
		}
		// Kiểm tra trong các Courses
		for _, course := range rule.Courses {
			for _, name := range course.NamesNorm {
				if name == normName {
					return course.IsLive
				}
			}
		}
	}
	return false
}

// isMissingItemLive kiểm tra xem một kết quả phân tích có phải dành cho vaccine sống hay không.
func (e *Engine) isMissingItemLive(res AnalysisResult) bool {
	// Hacky logic from Python parity: check description strings or match rule
	desc := strings.ToLower(res.Description)
	
	// Check known live vaccine keywords in description
	liveKeywords := []string{"imojev", "mvvac", "mmr", "varicella", "thủy đậu", "sởi", "quai bị", "rubella", "rota", "bại liệt (uống)"}
	for _, kw := range liveKeywords {
		if strings.Contains(desc, kw) {
			return true
		}
	}

	// Match back to rules
	if res.VaccineNameForPopup != "" {
		for _, rule := range e.Rules {
			if rule.DisplayName == res.VaccineNameForPopup || rule.GroupDisplayName == res.VaccineNameForPopup {
				if rule.IsLive {
					return true
				}
			}
		}
	}

	return false
}
