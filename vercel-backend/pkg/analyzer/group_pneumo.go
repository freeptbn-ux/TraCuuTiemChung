package analyzer

import (
	"fmt"
	"strings"

	"vercel-backend/pkg/models"
)

// processPneumoRules xử lý logic phức hợp cho nhóm vắc xin Phế cầu (Pneumococcal).
// Bao gồm: Prevenar 13, Synflorix, Vaxneuvance và Pneumovax 23.
func (e *Engine) processPneumoRules(administeredMap map[string][]models.VaccineRecord) []AnalysisResult {
	var results []AnalysisResult

	// 1. Thu thập bản ghi
	prevenarRule := e.Rules["Prevenar13"]
	vaxneuvanceRule := e.Rules["Vaxneuvance"]
	synflorixRule := e.Rules["Synflorix"]
	pneumovaxRule := e.Rules["Pneumovax23"]

	prevenarRecs := e.getMatchingRecords(prevenarRule.NamesNorm, administeredMap)
	vaxneuvanceRecs := e.getMatchingRecords(vaxneuvanceRule.NamesNorm, administeredMap)
	synflorixRecs := e.getMatchingRecords(synflorixRule.NamesNorm, administeredMap)
	pneumovaxRecs := e.getMatchingRecords(pneumovaxRule.NamesNorm, administeredMap)

	hasPneumovax := len(pneumovaxRecs) > 0

	// 2. Kiểm tra việc tiêm lẫn (Mixed series)
	activeSeriesKeys := []string{}
	activeSeriesNames := []string{}
	if len(prevenarRecs) > 0 {
		activeSeriesKeys = append(activeSeriesKeys, "Prevenar13")
		activeSeriesNames = append(activeSeriesNames, prevenarRule.DisplayName)
	}
	if len(vaxneuvanceRecs) > 0 {
		activeSeriesKeys = append(activeSeriesKeys, "Vaxneuvance")
		activeSeriesNames = append(activeSeriesNames, vaxneuvanceRule.DisplayName)
	}
	if len(synflorixRecs) > 0 {
		activeSeriesKeys = append(activeSeriesKeys, "Synflorix")
		activeSeriesNames = append(activeSeriesNames, synflorixRule.DisplayName)
	}

	if len(activeSeriesKeys) > 1 {
		mixedNames := []string{}
		for _, k := range activeSeriesKeys {
			mixedNames = append(mixedNames, e.Rules[k].DisplayName)
		}
		results = append(results, AnalysisResult{
			VaccineNameForPopup: "Phế cầu (nhiều loại)",
			Description:         fmt.Sprintf("Cảnh báo: Đã ghi nhận tiêm xen kẽ các loại phế cầu (%s). Không nên sử dụng xen kẽ.", strings.Join(mixedNames, " và ")),
			StatusTags:          []string{"error_interchange", "pneumo_mixed"},
		})
		return results
	}

	// 3. Skip logic: Nếu đã tiêm Pneumovax23 thì skip tất cả PCV
	if hasPneumovax {
		return results
	}

	// 4. Quyết định phác đồ PCV nào để recommend
	_, _, patientAgeYears := GetAgeAtDate(e.DOB, e.AnalysisDate)
	
	// Determine which PCV to show
	var pcvToProcess []string
	if len(activeSeriesKeys) > 0 {
		// Nếu đã tiêm, chỉ hiện loại đang tiêm (hoặc loại đầu tiên nếu mixed)
		pcvToProcess = []string{activeSeriesKeys[0]}
	} else {
		// Chưa tiêm loại nào, hiện cả 3
		pcvToProcess = []string{"Prevenar13", "Vaxneuvance", "Synflorix"}
	}

	for _, pKey := range pcvToProcess {
		pRule := e.Rules[pKey]
		resList := e.checkAgeDependentSeries(pKey, pRule, administeredMap)
		
		if patientAgeYears >= 2 && len(activeSeriesKeys) > 0 {
			// Logic đặc biệt cho trẻ trên 2 tuổi ĐÃ TIÊM PCV
			// Count ALL PCV doses (Synflorix + Prevenar 13 + Vaxneuvance)
			allPcvNames := []string{"synflorix", "prevenar13", "vaxneuvance", "prevenar 13"}
			pcvCount := len(e.getMatchingRecords(allPcvNames, administeredMap))
			
			if pcvCount < 3 {
				results = append(results, AnalysisResult{
					VaccineNameForPopup: pneumovaxRule.DisplayName,
					Description:         fmt.Sprintf("%s: Có thể tiêm 1 mũi để hoàn thành phác đồ phế cầu (do đã trên 2 tuổi và đã tiêm < 3 mũi %s).", pneumovaxRule.DisplayName, pRule.DisplayName),
					EarliestNextDoseDate: &e.AnalysisDate,
					StatusTags:          []string{"info", "alternative_completion"},
				})
				return results // Skip others like in Python
			} else if pcvCount == 3 {
				if len(resList) > 0 {
					results = append(results, AnalysisResult{
						VaccineNameForPopup: pneumovaxRule.DisplayName,
						Description:         fmt.Sprintf("%s: Có thể tiêm 1 mũi thay thế cho mũi 4 của %s (do đã trên 2 tuổi).", pneumovaxRule.DisplayName, pRule.DisplayName),
						EarliestNextDoseDate: &e.AnalysisDate,
						StatusTags:          []string{"info", "alternative_booster"},
					})
					// In Python, it skips the PCV recommendation if Pneumovax is suggested as alternative booster
					continue 
				}
			}
		}
		results = append(results, resList...)
	}

	// 5. Logic cho Pneumovax 23 (nếu chưa tiêm và chưa được xử lý ở trên)
	if !hasPneumovax {
		// Check if already suggested as alternative
		alreadySuggested := false
		for _, r := range results {
			// Robust matching for already suggested Pneumovax 23
			normR := NormalizeForMatch(r.VaccineNameForPopup)
			if strings.Contains(normR, "pneumovax23") || strings.Contains(normR, "pneumo23") {
				alreadySuggested = true
				break
			}
		}

		if !alreadySuggested {
			resPvList := e.checkSingleSeries("Pneumovax23", pneumovaxRule, administeredMap)
			results = append(results, resPvList...)
		}
	}

	return results
}
