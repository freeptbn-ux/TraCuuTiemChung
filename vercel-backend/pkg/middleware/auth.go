package middleware

import (
	"net/http"
	"vercel-backend/pkg/config"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// RequestID adds a unique ID to each request
func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := c.GetHeader("X-Request-Id")
		if requestID == "" {
			requestID = uuid.New().String()
		}
		c.Set("request_id", requestID)
		c.Writer.Header().Set("X-Request-Id", requestID)
		c.Next()
	}
}

// AuthRequired middleware checks for X-API-KEY header
func AuthRequired(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		apiKey := c.GetHeader("X-API-KEY")
		if apiKey == "" {
			apiKey = c.Query("api_key")
		}

		if apiKey != cfg.X_API_KEY {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"status":     "error",
				"error_code": "UNAUTHORIZED",
				"message":    "Invalid API Key",
				"request_id": c.GetString("request_id"),
			})
			return
		}
		c.Next()
	}
}
