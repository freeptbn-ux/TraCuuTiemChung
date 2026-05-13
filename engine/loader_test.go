package engine

import (
	"testing"
	"path/filepath"
)

func TestLoadVaccineRules(t *testing.T) {
	// Use absolute path or relative path to find the JSON file
	// Since we are in engine/, the assets are in ../vercel-backend/assets/
	filePath := filepath.Join("..", "vercel-backend", "assets", "vaccine_rules.json")

	rules, err := LoadVaccineRules(filePath)
	if err != nil {
		t.Fatalf("Failed to load vaccine rules: %v", err)
	}

	if len(rules) == 0 {
		t.Error("Loaded rules are empty")
	}

	// Assert rule "Prevenar13" có trường NamesNorm đã được lower-case
	rule, ok := rules["Prevenar13"]
	if !ok {
		t.Fatal("Rule 'Prevenar13' not found in loaded rules")
	}

	if len(rule.NamesNorm) == 0 {
		t.Error("Prevenar13 NamesNorm is empty")
	}

	found := false
	for _, name := range rule.NamesNorm {
		if name == "prevenar 13" {
			found = true
			break
		}
	}

	if !found {
		t.Errorf("Prevenar13 NamesNorm does not contain 'prevenar 13', got %v", rule.NamesNorm)
	}

	// Verify deep normalization (Members)
	macyw, ok := rules["MeningococcalACYW_Group"]
	if ok {
		member, ok := macyw.Members["MENQUADFI"]
		if ok {
			if len(member.NamesNorm) == 0 {
				t.Error("MENQUADFI NamesNorm is empty")
			}
			foundMQ := false
			for _, name := range member.NamesNorm {
				if name == "menquadfi" {
					foundMQ = true
					break
				}
			}
			if !foundMQ {
				t.Errorf("MENQUADFI NamesNorm does not contain 'menquadfi', got %v", member.NamesNorm)
			}
		} else {
			t.Error("Member 'MENQUADFI' not found in MeningococcalACYW_Group")
		}
	}

	// Verify deep normalization (Courses)
	je, ok := rules["JE_Group"]
	if ok {
		if len(je.Courses) == 0 {
			t.Error("JE_Group Courses is empty")
		} else {
			course := je.Courses[0] // Imojev
			if len(course.NamesNorm) == 0 {
				t.Error("JE_Group course 0 NamesNorm is empty")
			}
			if course.NamesNorm[0] != "imojev" {
				t.Errorf("JE_Group course 0 NamesNorm[0] = %q, want 'imojev'", course.NamesNorm[0])
			}
		}
	}
}
