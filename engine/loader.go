package engine

import (
	"encoding/json"
	"os"
	"tracuutiemchung-engine/engine/utils"
	"tracuutiemchung-engine/models"
)

// LoadVaccineRules reads the JSON file and returns a map of normalized vaccine rules.
func LoadVaccineRules(filePath string) (map[string]models.VaccineRule, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	var rules map[string]models.VaccineRule
	if err := json.Unmarshal(data, &rules); err != nil {
		return nil, err
	}

	// Normalize names in-memory
	for id, rule := range rules {
		// Normalize rule level RawNames
		if len(rule.RawNames) > 0 {
			rule.NamesNorm = make([]string, len(rule.RawNames))
			for i, name := range rule.RawNames {
				rule.NamesNorm[i] = utils.NormalizeVaccineName(name)
			}
		}

		// Normalize RawNamesMembers into NamesNormGroup
		if len(rule.RawNamesMembers) > 0 {
			for _, names := range rule.RawNamesMembers {
				for _, name := range names {
					rule.NamesNormGroup = append(rule.NamesNormGroup, utils.NormalizeVaccineName(name))
				}
			}
		}

		// Normalize courses level RawNames
		for i := range rule.Courses {
			if len(rule.Courses[i].RawNames) > 0 {
				rule.Courses[i].NamesNorm = make([]string, len(rule.Courses[i].RawNames))
				for j, name := range rule.Courses[i].RawNames {
					rule.Courses[i].NamesNorm[j] = utils.NormalizeVaccineName(name)
				}
			}
		}

		// Normalize members level RawNames
		for memberID, member := range rule.Members {
			if len(member.RawNames) > 0 {
				member.NamesNorm = make([]string, len(member.RawNames))
				for i, name := range member.RawNames {
					member.NamesNorm[i] = utils.NormalizeVaccineName(name)
				}
				rule.Members[memberID] = member
			}
		}

		// Update the rule in the map
		rules[id] = rule
	}

	return rules, nil
}
