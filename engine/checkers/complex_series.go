package checkers

import (
	"tracuutiemchung-engine/models"
)

// ApplyComplexInteractions handles special cases like Measles vs MMR and Meningococcal interactions.
// It follows Phase 05 plan to handle "MVVAC vs MMR" and "VA-MENGOC-BC vs MenQuadfi".
func ApplyComplexInteractions(results []models.MissingItem, adminMap map[string][]models.AdministeredDose) []models.MissingItem {
	filteredResults := []models.MissingItem{}

	// Check for MMR (using the display name as expected by checkers)
	hasMMR := len(adminMap["Vắc xin Sởi-Quai bị-Rubella (MMR-II/Priorix)"]) > 0

	for _, item := range results {
		// 1. MVVAC vs MMR: Measles is covered by MMR
		if item.VaccineName == "MVVAC (Sởi đơn)" && hasMMR {
			continue
		}

		// 2. VA-MENGOC-BC vs MenQuadfi: If MenQuadfi (better) is already injected, do NOT recommend VA-MENGOC-BC (inferior).
		if item.VaccineName == "VA - MENGOC - BC (Não mô cầu BC)" {
			hasACYW := len(adminMap["Vắc xin Não mô cầu ACYW-135 (Menactra/MenQuadfi)"]) > 0
			if hasACYW {
				continue
			}
		}

		filteredResults = append(filteredResults, item)
	}

	return filteredResults
}
