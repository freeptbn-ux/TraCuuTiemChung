package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"vercel-backend/internal/analyzer"
	"vercel-backend/internal/portal"
)

func main() {
	htmlPath := "../test/Gia-Han.html"
	absPath, _ := filepath.Abs(htmlPath)
	content, err := os.ReadFile(absPath)
	if err != nil {
		fmt.Printf("Error reading HTML: %v\n", err)
		os.Exit(1)
	}

	pc := portal.NewPortalClient("", "")
	detail, err := pc.ParsePatientDetail(string(content))
	if err != nil {
		fmt.Printf("Error parsing HTML: %v\n", err)
		os.Exit(1)
	}

	dob, _ := analyzer.ParseDateDDMMYYYY(detail.PatientInfo.Birth)
	sysDate, _ := analyzer.ParseDateDDMMYYYY(detail.PatientInfo.SystemDate)
	if sysDate.IsZero() {
		sysDate = time.Now()
	}

	rulesPath := "assets/vaccine_rules.json"
	engine, err := analyzer.NewEngine(rulesPath, dob, sysDate)
	if err != nil {
		fmt.Printf("Error creating engine: %v\n", err)
		os.Exit(1)
	}

	results := engine.Analyze(detail.History)

	output := map[string]interface{}{
		"patient_name": detail.PatientInfo.Name,
		"dob":          detail.PatientInfo.Birth,
		"analysis_date": detail.PatientInfo.SystemDate,
		"missing_vaccines": results,
		"administered_vaccines": detail.History,
	}

	jsonData, _ := json.MarshalIndent(output, "", "  ")
	os.WriteFile("../testdata/go_giahan_output.json", jsonData, 0644)
	fmt.Println("Generated testdata/go_giahan_output.json")
}
