package handler

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	"vercel-backend/pkg/analyzer"
	"vercel-backend/pkg/portal"
)

func TestLiveVNVC(t *testing.T) {
	username := "bn_dv_tcdvquevo"
	password := "Tinh@2027"
	phone := "0388634123"

	pcLocal := portal.NewPortalClient(username, password, nil)

	fmt.Println("Logging in...")
	err := pcLocal.Login()
	if err != nil {
		t.Fatalf("Login failed: %v", err)
	}

	fmt.Println("Looking up phone:", phone)
	results, err := pcLocal.LookupPatients(phone)
	if err != nil {
		t.Fatalf("Lookup failed: %v", err)
	}

	if len(results) == 0 {
		t.Fatalf("No patients found in lookup!")
	}

	var targetPatientID string
	var targetName string
	for _, p := range results {
		fmt.Printf("Found patient: %s (ID: %s)\n", p.Name, p.ID)
		if strings.Contains(strings.ToLower(p.Name), "khôi") {
			targetPatientID = p.ID
			targetName = p.Name
		}
	}

	if targetPatientID == "" {
		targetPatientID = results[0].ID
		targetName = results[0].Name
	}

	fmt.Printf("\nAnalyzing patient: %s (ID: %s)\n", targetName, targetPatientID)

	detail, err := pcLocal.GetVaccinationHistory(targetPatientID)
	if err != nil {
		t.Fatalf("History fetch failed: %v", err)
	}

	fmt.Printf("History contains %d records.\n", len(detail.History))

	dob, _ := time.Parse("02/01/2006", detail.PatientInfo.Birth)
	analysisDate := time.Now()
	if detail.PatientInfo.SystemDate != "" {
		if sd, err := time.Parse("02/01/2006", detail.PatientInfo.SystemDate); err == nil {
			analysisDate = sd
		}
	}

	eng, err := analyzer.NewEngine("../assets/vaccine_rules.json", dob, analysisDate)
	if err != nil {
		t.Fatalf("Analyzer failed: %v", err)
	}

	rawResults := eng.Analyze(detail.History)
	fmt.Printf("Analysis complete: %d results.\n", len(rawResults))

	// Format output similar to handler
	type Recommendation struct {
		VaccineName string   `json:"vaccine_name"`
		RuleType    string   `json:"rule_type"`
		Status      string   `json:"status"`
		NextDose    string   `json:"next_dose"`
		Message     string   `json:"message"`
		StatusTags  []string `json:"status_tags"`
	}

	recommendations := make([]Recommendation, 0, len(rawResults))
	for _, res := range rawResults {
		status := "DUE_LATER"
		hasDue, hasOverdue, hasWarning, hasCompleted := false, false, false, false
		for _, tag := range res.StatusTags {
			if tag == "due" || tag == "eligible" {
				hasDue = true
			} else if tag == "overdue" {
				hasOverdue = true
			} else if strings.Contains(tag, "error") || tag == "warning" {
				hasWarning = true
			} else if tag == "completed" {
				hasCompleted = true
			}
		}

		if hasDue {
			status = "DUE_NOW"
		} else if hasOverdue {
			status = "OVERDUE"
		} else if hasWarning {
			status = "NEEDS_REVIEW"
		} else if hasCompleted {
			status = "COMPLETED"
		}

		nextDoseStr := ""
		if res.EarliestNextDoseDate != nil {
			nextDoseStr = res.EarliestNextDoseDate.Format("02/01/2006")
		}

		recommendations = append(recommendations, Recommendation{
			VaccineName: res.VaccineNameForPopup,
			RuleType:    "standard",
			Status:      status,
			NextDose:    nextDoseStr,
			Message:     res.Description,
			StatusTags:  res.StatusTags,
		})
	}

	out, _ := json.MarshalIndent(map[string]interface{}{
		"patient_name":     detail.PatientInfo.Name,
		"dob":              detail.PatientInfo.Birth,
		"missing_vaccines": recommendations,
	}, "", "  ")

	os.WriteFile("live_output.json", out, 0644)
	fmt.Println("Analysis saved to live_output.json")
}
