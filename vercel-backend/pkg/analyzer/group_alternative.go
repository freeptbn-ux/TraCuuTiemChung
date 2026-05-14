package analyzer

import (
	"fmt"
	"sort"

	"vercel-backend/pkg/models"
)

// checkAlternativeCoursesMinAgeGroup xử lý nhóm vắc xin có nhiều phác đồ thay thế dựa trên độ tuổi bắt đầu (ví dụ: Rota).
func (e *Engine) checkAlternativeCoursesMinAgeGroup(ruleKey string, rule VaccineRule, administeredMap map[string][]models.VaccineRecord) *AnalysisResult {
	// Thu thập tất cả các bản ghi liên quan đến tất cả các phác đồ trong nhóm
	var allRecords []models.VaccineRecord
	for _, course := range rule.Courses {
		recs := e.getMatchingRecords(course.NamesNorm, administeredMap)
		allRecords = append(allRecords, recs...)
	}

	// Sắp xếp theo ngày tiêm
	sort.Slice(allRecords, func(i, j int) bool {
		return allRecords[i].Date.Before(allRecords[j].Date)
	})

	numDoses := len(allRecords)
	monthsNow, _, _ := GetAgeAtDate(e.DOB, e.AnalysisDate)

	// Kiểm tra xem đã quá tuổi hoàn thành chưa
	if rule.MaxAgeMonthsForCompletionGroup > 0 && monthsNow > rule.MaxAgeMonthsForCompletionGroup {
		if numDoses == 0 {
			return &AnalysisResult{
				VaccineNameForPopup: rule.GroupDisplayName,
				Description:         fmt.Sprintf("%s (Đã quá tuổi chỉ định tiêm)", rule.GroupDisplayName),
				StatusTags:          []string{"too_old"},
			}
		}
		// Nếu đã tiêm rồi nhưng chưa xong và đã quá tuổi hoàn thành?
		// Tùy theo logic, nhưng thường là không tiêm tiếp.
		return nil 
	}

	if numDoses == 0 {
		// Chưa tiêm mũi nào
		if rule.MaxAgeMonthsToStartFirstDoseGroup > 0 && monthsNow > rule.MaxAgeMonthsToStartFirstDoseGroup {
			return &AnalysisResult{
				VaccineNameForPopup: rule.GroupDisplayName,
				Description:         fmt.Sprintf("%s (Đã quá tuổi bắt đầu tiêm)", rule.GroupDisplayName),
				StatusTags:          []string{"too_old"},
			}
		}

		// Gợi ý tất cả các phác đồ khả dụng
		var courseNames []string
		for _, course := range rule.Courses {
			courseNames = append(courseNames, course.Display)
		}

		earliestDate := e.DOB
		if rule.MinAgeWeeksAtFirstDose > 0 {
			earliestDate = e.DOB.AddDate(0, 0, rule.MinAgeWeeksAtFirstDose*7)
		} else if rule.MinAgeMonthsAtFirstDose > 0 {
			earliestDate = AddMonths(e.DOB, rule.MinAgeMonthsAtFirstDose)
		}

		tags := []string{"due"}
		if e.AnalysisDate.Before(earliestDate) {
			tags = []string{"too_young"}
		}

		return &AnalysisResult{
			VaccineNameForPopup:  rule.GroupDisplayName,
			Description:          fmt.Sprintf("%s (Bắt đầu tiêm: %v)", rule.GroupDisplayName, courseNames),
			EarliestNextDoseDate: &earliestDate,
			StatusTags:           tags,
		}
	}

	// Đã tiêm ít nhất 1 mũi. Xác định phác đồ đã bắt đầu.
	// Tìm phác đồ chứa mũi tiêm đầu tiên.
	var startedCourse *Course
	firstRecord := allRecords[0]
	normFirst := NormalizeVaccineName(firstRecord.VaccineName)

	for i := range rule.Courses {
		for _, nameNorm := range rule.Courses[i].NamesNorm {
			if nameNorm == normFirst {
				startedCourse = &rule.Courses[i]
				break
			}
		}
		if startedCourse != nil {
			break
		}
	}

	if startedCourse == nil {
		// Không tìm thấy phác đồ tương ứng (có thể là vắc xin lạ)
		return nil
	}

	if numDoses >= startedCourse.DosesRequired {
		return nil
	}

	// Tính ngày tiêm tiếp theo dựa trên phác đồ đã bắt đầu
	lastDoseDate := allRecords[numDoses-1].Date
	intervalDays := 0
	if numDoses < len(startedCourse.MinIntervalDays) && startedCourse.MinIntervalDays[numDoses] != nil {
		intervalDays = *startedCourse.MinIntervalDays[numDoses]
	}
	earliestDate := lastDoseDate.AddDate(0, 0, intervalDays)

	tags := []string{"due"}
	if e.AnalysisDate.Before(earliestDate) {
		tags = []string{"too_young"}
	}

	return &AnalysisResult{
		VaccineNameForPopup:  startedCourse.Display,
		Description:          fmt.Sprintf("%s (Tiêm mũi %d)", startedCourse.Display, numDoses+1),
		EarliestNextDoseDate: &earliestDate,
		StatusTags:           tags,
	}
}

// checkAlternativeCoursesAgeRangeGroup xử lý nhóm vắc xin có nhiều phác đồ thay thế dựa trên khoảng độ tuổi (ví dụ: JE, HepA).
func (e *Engine) checkAlternativeCoursesAgeRangeGroup(ruleKey string, rule VaccineRule, administeredMap map[string][]models.VaccineRecord) *AnalysisResult {
	// Thu thập tất cả các bản ghi
	var allRecords []models.VaccineRecord
	courseMap := make(map[string]*Course) // nameNorm -> Course

	for i := range rule.Courses {
		recs := e.getMatchingRecords(rule.Courses[i].NamesNorm, administeredMap)
		allRecords = append(allRecords, recs...)
		for _, nameNorm := range rule.Courses[i].NamesNorm {
			courseMap[nameNorm] = &rule.Courses[i]
		}
	}

	sort.Slice(allRecords, func(i, j int) bool {
		return allRecords[i].Date.Before(allRecords[j].Date)
	})

	numDoses := len(allRecords)
	monthsNow, _, _ := GetAgeAtDate(e.DOB, e.AnalysisDate)

	if numDoses == 0 {
		// Tìm phác đồ phù hợp với tuổi hiện tại
		var bestCourse *Course
		for i := range rule.Courses {
			c := &rule.Courses[i]
			if monthsNow >= c.MinAgeMonthsAtFirstDose && (c.MaxAgeYearsAtFirstDose == 0 || monthsNow <= c.MaxAgeYearsAtFirstDose*12) {
				bestCourse = c
				break
			}
		}

		if bestCourse == nil {
			return &AnalysisResult{
				VaccineNameForPopup: rule.GroupDisplayName,
				Description:         fmt.Sprintf("%s (Hiện không có phác đồ phù hợp với lứa tuổi)", rule.GroupDisplayName),
				StatusTags:          []string{"too_old"},
			}
		}

		earliestDate := AddMonths(e.DOB, bestCourse.MinAgeMonthsAtFirstDose)
		tags := []string{"due"}
		if e.AnalysisDate.Before(earliestDate) {
			tags = []string{"too_young"}
		}

		return &AnalysisResult{
			VaccineNameForPopup:  rule.GroupDisplayName,
			Description:          fmt.Sprintf("%s (Bắt đầu tiêm phác đồ %s)", rule.GroupDisplayName, bestCourse.Display),
			EarliestNextDoseDate: &earliestDate,
			StatusTags:           tags,
		}
	}

	// Đã tiêm. Kiểm tra trộn phác đồ.
	var startedCourse *Course
	mixed := false
	for _, rec := range allRecords {
		norm := NormalizeVaccineName(rec.VaccineName)
		c := courseMap[norm]
		if startedCourse == nil {
			startedCourse = c
		} else if c != nil && c != startedCourse {
			mixed = true
		}
	}

	if startedCourse == nil {
		return nil
	}

	// Cảnh báo nếu trộn phác đồ
	warning := ""
	if mixed {
		warning = " [Cảnh báo: Đã tiêm trộn các loại vắc xin khác nhau trong nhóm]"
	}

	// Kiểm tra hoàn thành
	if numDoses >= startedCourse.DosesRequired {
		// Kiểm tra booster cho Jevax (nếu có định nghĩa)
		if startedCourse.BoosterIntervalYears > 0 && numDoses == startedCourse.BoosterAfterDoseNumber {
			yearsNow := monthsNow / 12
			if startedCourse.BoosterMaxAgeYears == 0 || yearsNow < startedCourse.BoosterMaxAgeYears {
				lastDoseDate := allRecords[numDoses-1].Date
				earliestDate := AddYears(lastDoseDate, startedCourse.BoosterIntervalYears)
				
				tags := []string{"due"}
				if e.AnalysisDate.Before(earliestDate) {
					tags = []string{"too_young"}
				}

				return &AnalysisResult{
					VaccineNameForPopup:  startedCourse.Display,
					Description:          fmt.Sprintf("%s (Tiêm mũi nhắc lại mỗi %d năm)%s", startedCourse.Display, startedCourse.BoosterIntervalYears, warning),
					EarliestNextDoseDate: &earliestDate,
					StatusTags:           tags,
				}
			}
		}
		return nil
	}

	// Tiếp tục phác đồ đã bắt đầu
	lastDoseDate := allRecords[numDoses-1].Date
	intervalDays := 0
	if numDoses < len(startedCourse.MinIntervalDays) && startedCourse.MinIntervalDays[numDoses] != nil {
		intervalDays = *startedCourse.MinIntervalDays[numDoses]
	}
	earliestDate := lastDoseDate.AddDate(0, 0, intervalDays)

	tags := []string{"due"}
	if e.AnalysisDate.Before(earliestDate) {
		tags = []string{"too_young"}
	}

	return &AnalysisResult{
		VaccineNameForPopup:  startedCourse.Display,
		Description:          fmt.Sprintf("%s (Tiêm mũi %d)%s", startedCourse.Display, numDoses+1, warning),
		EarliestNextDoseDate: &earliestDate,
		StatusTags:           tags,
	}
}
