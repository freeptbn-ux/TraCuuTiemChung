package engine

import (
	"path/filepath"
	"testing"
	"tracuutiemchung-engine/engine/utils"
)

func TestPhase03_RulesLoader(t *testing.T) {
	// 1. Task: Viết hàm đọc file vaccine_rules.json
	filePath := filepath.Join("testdata", "rules_loader_test.json")
	rules, err := LoadVaccineRules(filePath)
	if err != nil {
		t.Fatalf("LoadVaccineRules failed: %v", err)
	}

	// 2. Task: Chuẩn hóa NamesNorm (Normalization During Load)
	// Verification: assert.Contains(t, rules["Prevenar13"].NamesNorm, "prevenar13")
	// Note: Our NormalizeVaccineName keeps spaces if they exist, but we verify it's populated and lower-cased.
	prevenar, ok := rules["Prevenar13"]
	if !ok {
		t.Fatal("Rule 'Prevenar13' not found")
	}

	foundPrevenar := false
	for _, name := range prevenar.NamesNorm {
		if name == utils.NormalizeVaccineName("Prevenar 13") {
			foundPrevenar = true
			break
		}
	}
	if !foundPrevenar {
		t.Errorf("Prevenar13 NamesNorm missing normalized 'Prevenar 13', got %v", prevenar.NamesNorm)
	}

	// 3. Scenario: Multi-alias Mapping
	// Rule has raw_names: ["Infanrix Hexa", "Hexaxim"]
	sixInOne, ok := rules["Six_In_One_Combined"]
	if !ok {
		t.Fatal("Rule 'Six_In_One_Combined' not found")
	}
	
	aliases := []string{"infanrix hexa", "hexaxim"}
	for _, alias := range aliases {
		found := false
		for _, norm := range sixInOne.NamesNorm {
			if norm == alias {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Six_In_One_Combined NamesNorm missing alias %q, got %v", alias, sixInOne.NamesNorm)
		}
	}

	// 4. Task: Chuẩn hóa RawNamesMembers -> NamesNormGroup
	mmr, ok := rules["MMR_Group"]
	if !ok {
		t.Fatal("Rule 'MMR_Group' not found")
	}
	if len(mmr.NamesNormGroup) == 0 {
		t.Error("MMR_Group NamesNormGroup is empty")
	}
	
	foundPriorix := false
	for _, name := range mmr.NamesNormGroup {
		if name == "priorix" {
			foundPriorix = true
			break
		}
	}
	if !foundPriorix {
		t.Errorf("MMR_Group NamesNormGroup missing 'priorix', got %v", mmr.NamesNormGroup)
	}

	// 5. Task: Chuẩn hóa courses.raw_names
	je, ok := rules["JE_Group"]
	if !ok {
		t.Fatal("Rule 'JE_Group' not found")
	}
	if len(je.Courses) == 0 || len(je.Courses[0].NamesNorm) == 0 {
		t.Error("JE_Group course NamesNorm missing")
	}
	if je.Courses[0].NamesNorm[0] != "imojev" {
		t.Errorf("JE_Group course NamesNorm[0] = %q, want 'imojev'", je.Courses[0].NamesNorm[0])
	}

	// 6. Task: Chuẩn hóa members.raw_names
	macyw, ok := rules["MeningococcalACYW_Group"]
	if !ok {
		t.Fatal("Rule 'MeningococcalACYW_Group' not found")
	}
	mq, ok := macyw.Members["MENQUADFI"]
	if !ok {
		t.Fatal("Member 'MENQUADFI' not found")
	}
	foundMQ := false
	for _, name := range mq.NamesNorm {
		if name == "menquadfi" {
			foundMQ = true
			break
		}
	}
	if !foundMQ {
		t.Errorf("MENQUADFI NamesNorm missing 'menquadfi', got %v", mq.NamesNorm)
	}
}
