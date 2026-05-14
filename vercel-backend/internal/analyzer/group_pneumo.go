package analyzer

import (
	"fmt"
	"strings"

	"vercel-backend/internal/models"
)

// processPneumoRules xử lý logic phức hợp cho nhóm vắc xin Phế cầu (Pneumococcal).
// Bao gồm: Prevenar 13, Synflorix, Vaxneuvance và Pneumovax 23.
func (e *Engine) processPneumoRules(administeredMap map[string][]models.VaccineRecord) []AnalysisResult {
	var results []AnalysisResult

	// 1. Thu thập bản ghi của từng loại
	prevenarRule := e.Rules["Prevenar13"]
	vaxneuvanceRule := e.Rules["Vaxneuvance"]
	synflorixRule := e.Rules["Synflorix"]
	pneumovaxRule := e.Rules["Pneumovax23"]

	prevenarRecs := e.getMatchingRecords(prevenarRule.NamesNorm, administeredMap)
	vaxneuvanceRecs := e.getMatchingRecords(vaxneuvanceRule.NamesNorm, administeredMap)
	synflorixRecs := e.getMatchingRecords(synflorixRule.NamesNorm, administeredMap)
	pneumovaxRecs := e.getMatchingRecords(pneumovaxRule.NamesNorm, administeredMap)

	// 2. Kiểm tra việc tiêm lẫn (Mixed series)
	var activeSeries []string
	if len(prevenarRecs) > 0 {
		activeSeries = append(activeSeries, "Prevenar 13")
	}
	if len(vaxneuvanceRecs) > 0 {
		activeSeries = append(activeSeries, "Vaxneuvance")
	}
	if len(synflorixRecs) > 0 {
		activeSeries = append(activeSeries, "Synflorix")
	}

	if len(activeSeries) > 1 {
		results = append(results, AnalysisResult{
			VaccineNameForPopup: "Phế cầu (Mixed)",
			Description:         fmt.Sprintf("Cảnh báo: Tiêm lẫn các loại phế cầu (%s). Nên tiêm cùng một loại để đạt hiệu quả tốt nhất.", strings.Join(activeSeries, " + ")),
			StatusTags:          []string{"error_interchange"},
		})
	}

	// 3. Quyết định phác đồ nào để recommend
	// Ưu tiên: Nếu đã bắt đầu loại nào thì tiếp tục loại đó.
	// Nếu chưa bắt đầu, hoặc đã bắt đầu nhưng yêu cầu parity là không hiện cái khác.
	
	hasSynflorix := len(synflorixRecs) > 0
	hasPrevenar := len(prevenarRecs) > 0
	hasVaxneuvance := len(vaxneuvanceRecs) > 0
	
	pcvCompleted := false

	if hasSynflorix {
		res := e.checkAgeDependentSeries("Synflorix", synflorixRule, administeredMap)
		if res != nil {
			results = append(results, *res)
		} else {
			pcvCompleted = true
		}
	} else if hasPrevenar {
		res := e.checkAgeDependentSeries("Prevenar13", prevenarRule, administeredMap)
		if res != nil {
			results = append(results, *res)
		} else {
			pcvCompleted = true
		}
	} else if hasVaxneuvance {
		res := e.checkAgeDependentSeries("Vaxneuvance", vaxneuvanceRule, administeredMap)
		if res != nil {
			results = append(results, *res)
		} else {
			pcvCompleted = true
		}
	} else {
		// Chưa tiêm loại nào, hiện cả 3 lựa chọn PCV
		resP := e.checkAgeDependentSeries("Prevenar13", prevenarRule, administeredMap)
		if resP != nil {
			results = append(results, *resP)
		}
		resV := e.checkAgeDependentSeries("Vaxneuvance", vaxneuvanceRule, administeredMap)
		if resV != nil {
			results = append(results, *resV)
		}
		resS := e.checkAgeDependentSeries("Synflorix", synflorixRule, administeredMap)
		if resS != nil {
			results = append(results, *resS)
		}
	}

	monthsNow, _, _ := GetAgeAtDate(e.DOB, e.AnalysisDate)

	// 4. Logic cho Pneumovax 23 (Phế cầu 23)
	if len(pneumovaxRecs) == 0 {
		resPv := e.checkSingleSeries("Pneumovax23", pneumovaxRule, administeredMap)
		if resPv != nil {
			if monthsNow >= 24 {
				if !pcvCompleted {
					resPv.Description = fmt.Sprintf("%s (Có thể tiêm bổ sung hoặc thay thế để mở rộng bảo vệ)", resPv.VaccineNameForPopup)
				} else {
					resPv.Description = fmt.Sprintf("%s (Tiêm mũi nhắc mở rộng bảo vệ)", resPv.VaccineNameForPopup)
				}
			} else {
				// Vẫn hiện nhưng là too_young (theo parity)
				resPv.StatusTags = []string{"too_young"}
			}
			results = append(results, *resPv)
		}
	}

	return results
}
