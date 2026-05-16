package handler

import (
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"vercel-backend/assets"
	"vercel-backend/pkg/analyzer"
	"vercel-backend/pkg/config"
	"vercel-backend/pkg/logger"
	"vercel-backend/pkg/portal"
	"log/slog"
	"vercel-backend/pkg/middleware"

	"github.com/redis/go-redis/v9"
)

var (
	cfg    *config.Config
	pc     *portal.PortalClient
	router *gin.Engine
)

func init() {
	// 1. Core Services
	logger.InitLogger()
	cfg = config.LoadConfig()

	var redisClient *redis.Client
	if cfg.Redis.URL != "" {
		opts, err := redis.ParseURL(cfg.Redis.URL)
		if err == nil {
			redisClient = redis.NewClient(opts)
			slog.Info("redis connected", "url", cfg.Redis.URL)
		} else {
			slog.Error("failed to parse redis url", "error", err)
		}
	}

	pc = portal.NewPortalClient(cfg.PORTAL_USERNAME, cfg.PORTAL_PASSWORD, redisClient)

	// 2. Gin Engine Setup
	if os.Getenv("VERCEL") == "1" {
		gin.SetMode(gin.ReleaseMode)
	}

	router = gin.New()

	// 3. Global Middleware
	router.Use(
		middleware.RequestID(),
		middleware.LoggerMiddleware(),
		middleware.ErrorHandler(),
		gin.Recovery(),
	)

	// 4. Routes Definition
	// Root & Favicon to avoid 404 logs
	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "online",
			"service": "Tra Cuu Tiem Chung API",
			"message": "Welcome! The API is running correctly.",
			"docs":    "/api/health",
		})
	})

	router.GET("/favicon.ico", func(c *gin.Context) {
		c.Status(http.StatusNoContent)
	})

	router.GET("/favicon.png", func(c *gin.Context) {
		c.Status(http.StatusNoContent)
	})

	// Health check
	router.GET("/api/health", func(c *gin.Context) {
		sendSuccess(c, gin.H{
			"environment": os.Getenv("VERCEL_ENV"),
			"message":     "Backend is running",
		})
	})

	// 404 Handler
	router.NoRoute(func(c *gin.Context) {
		slog.Warn("route not found", "path", c.Request.URL.Path)
		sendError(c, http.StatusNotFound, "Endpoint not found: "+c.Request.URL.Path)
	})

	// API Group with Auth & Rate Limit
	api := router.Group("/api")
	api.Use(
		middleware.AuthRequired(cfg),
		middleware.RateLimit(pc.RedisClient(), 50, time.Minute),
	)
	{
		api.POST("/lookup", handleLookup)
		api.POST("/analyze", handleAnalyze)
		api.GET("/debug/portal", handleDebugPortal)
	}
}

func handleDebugPortal(c *gin.Context) {
	status := pc.CheckPortalConnectivity()
	sendSuccess(c, status)
}


// Handler is the entry point for Vercel Go runtime
func Handler(w http.ResponseWriter, r *http.Request) {
	router.ServeHTTP(w, r)
}

// StandardResponse is the base for all API responses
type StandardResponse struct {
	Status    string      `json:"status"`
	Data      interface{} `json:"data,omitempty"`
	Message   string      `json:"message,omitempty"`
	Detail    string      `json:"detail,omitempty"` // For compatibility with Android app
	RequestID string      `json:"request_id"`
}

func sendSuccess(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, StandardResponse{
		Status:    "success",
		Data:      data,
		RequestID: c.GetString("request_id"),
	})
}

func sendError(c *gin.Context, code int, message string) {
	c.AbortWithStatusJSON(code, StandardResponse{
		Status:    "error",
		Message:   message,
		Detail:    message,
		RequestID: c.GetString("request_id"),
	})
}

type LookupRequest struct {
	Phone string `json:"phone" binding:"required"`
}

func handleLookup(c *gin.Context) {
	var req LookupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		slog.Error("lookup binding failed", "error", err)
		sendError(c, http.StatusBadRequest, "Phone number is required: "+err.Error())
		return
	}

	results, err := pc.LookupPatients(req.Phone)
	if err != nil {
		slog.Error("lookup failed", "phone", req.Phone, "error", err)
		c.Error(err) // Centralized error handling
		return
	}

	if len(results) == 0 {
		slog.Warn("no patients found for phone", "phone", req.Phone)
	} else {
		slog.Info("found patients", "phone", req.Phone, "count", len(results))
	}

	sendSuccess(c, results)
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
		sendError(c, http.StatusBadRequest, "Patient ID is required")
		return
	}

	// 1. Fetch details from portal
	detail, err := pc.GetVaccinationHistory(req.PatientID)
	if err != nil {
		slog.Error("history fetch failed", "patient_id", req.PatientID, "error", err)
		c.Error(err)
		return
	}

	// 2. Parse dates
	dob, err := time.Parse("02/01/2006", detail.PatientInfo.Birth)
	if err != nil {
		c.Error(fmt.Errorf("failed to parse DOB: %w", err))
		return
	}

	analysisDate := time.Now()
	if detail.PatientInfo.SystemDate != "" {
		if sd, err := time.Parse("02/01/2006", detail.PatientInfo.SystemDate); err == nil {
			analysisDate = sd
		}
	}

	// 3. Initialize engine and analyze using embedded rules
	engine, err := analyzer.NewEngineFromBytes(assets.VaccineRulesJSON, dob, analysisDate)
	if err != nil {
		c.Error(fmt.Errorf("failed to initialize analyzer: %w", err))
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

	sendSuccess(c, AnalyzeResponse{
		PatientName:          detail.PatientInfo.Name,
		DOB:                  detail.PatientInfo.Birth,
		AnalysisDate:         detail.PatientInfo.SystemDate,
		MissingVaccines:      recommendations,
		AdministeredVaccines: history,
	})
}
