package handler

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"vercel-backend/pkg/analyzer"
	"vercel-backend/pkg/config"
	"vercel-backend/pkg/portal"
)

var (
	cfg *config.Config
	pc  *portal.PortalClient
)

func init() {
	cfg = config.LoadConfig()
	pc = portal.NewPortalClient(cfg.PORTAL_USERNAME, cfg.PORTAL_PASSWORD)
	// Vercel environment usually doesn't need to set Gin mode, but it's good practice
	if os.Getenv("VERCEL") == "1" {
		gin.SetMode(gin.ReleaseMode)
	}
}

// AuthRequired middleware checks for X-API-KEY header
func AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		apiKey := c.GetHeader("X-API-KEY")
		if apiKey == "" {
			// Fallback to query param if needed (some clients might use it)
			apiKey = c.Query("api_key")
		}

		if apiKey != cfg.X_API_KEY {
			c.JSON(http.StatusForbidden, gin.H{
				"status":  "error",
				"message": "Invalid API Key",
			})
			c.Abort()
			return
		}
		c.Next()
	}
}

// Handler is the entry point for Vercel Go runtime
func Handler(w http.ResponseWriter, r *http.Request) {
	router := gin.New()
	router.Use(gin.Recovery())

	// Health check
	router.GET("/api/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":      "ok",
			"message":     "Backend is running",
			"environment": os.Getenv("VERCEL_ENV"),
		})
	})

	// API Group with Auth
	api := router.Group("/api")
	api.Use(AuthRequired())
	{
		api.POST("/lookup", handleLookup)
		api.POST("/analyze", handleAnalyze)
	}

	router.ServeHTTP(w, r)
}

type LookupRequest struct {
	Phone string `json:"phone" binding:"required"`
}

func handleLookup(c *gin.Context) {
	var req LookupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Phone number is required"})
		return
	}

	results, err := pc.LookupPatients(req.Phone)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": fmt.Sprintf("Portal lookup error: %v", err),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   results,
	})
}

type Recommendation struct {
	VaccineName string   `json:"vaccine_name"`
	RuleType    string   `json:"rule_type"`
	Status      string   `json:"status"`
	NextDose    string   `json:"next_dose"`
	Message     string   `json:"message"`
	StatusTags  []string `json:"status_tags"`
}

type AdministeredVaccine struct {
	VaccineName string `json:"vaccine_name"`
	Date        string `json:"date"`
	Dose        string `json:"dose"`
	Provider    string `json:"provider"`
}

type AnalyzeRequest struct {
	PatientID string `json:"patient_id" binding:"required"`
}

type AnalyzeResponse struct {
	PatientName          string                `json:"patient_name"`
	DOB                  string                `json:"dob"`
	AnalysisDate         string                `json:"analysis_date"`
	MissingVaccines      []Recommendation      `json:"missing_vaccines"`
	AdministeredVaccines []AdministeredVaccine `json:"administered_vaccines"`
}

func handleAnalyze(c *gin.Context) {
	var req AnalyzeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Patient ID is required"})
		return
	}

	// 1. Fetch details from portal
	detail, err := pc.GetVaccinationHistory(req.PatientID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": fmt.Sprintf("Portal history error: %v", err),
		})
		return
	}

	// 2. Parse dates
	dob, err := time.Parse("02/01/2006", detail.PatientInfo.Birth)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": fmt.Sprintf("Failed to parse DOB: %v", err),
		})
		return
	}

	analysisDate := time.Now()
	if detail.PatientInfo.SystemDate != "" {
		if sd, err := time.Parse("02/01/2006", detail.PatientInfo.SystemDate); err == nil {
			analysisDate = sd
		}
	}

	// 3. Initialize engine and analyze
	rulesPath := filepath.Join("assets", "vaccine_rules.json")
	if _, err := os.Stat(rulesPath); os.IsNotExist(err) {
		rulesPath = "../assets/vaccine_rules.json"
	}

	engine, err := analyzer.NewEngine(rulesPath, dob, analysisDate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": fmt.Sprintf("Failed to initialize analyzer: %v", err),
		})
		return
	}

	rawResults := engine.Analyze(detail.History)

	// 4. Format for Android
	recommendations := make([]Recommendation, 0, len(rawResults))
	for _, res := range rawResults {
		status := "DUE_LATER"
		hasDue := false
		hasOverdue := false
		hasWarning := false
		hasCompleted := false

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

	history := make([]AdministeredVaccine, 0, len(detail.History))
	for _, rec := range detail.History {
		history = append(history, AdministeredVaccine{
			VaccineName: rec.VaccineName,
			Date:        rec.Date.Format("02/01/2006"),
			Dose:        rec.Dose,
			Provider:    "VNCDC",
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data": AnalyzeResponse{
			PatientName:          detail.PatientInfo.Name,
			DOB:                  detail.PatientInfo.Birth,
			AnalysisDate:         detail.PatientInfo.SystemDate,
			MissingVaccines:      recommendations,
			AdministeredVaccines: history,
		},
	})
}
