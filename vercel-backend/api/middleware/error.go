package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// ErrorHandler middleware handles errors collected by c.Error()
func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) > 0 {
			// If response already written, we can't do much about the status code
			// but we can log the errors
			if c.Writer.Written() {
				return
			}

			err := c.Errors.Last()
			
			statusCode := http.StatusInternalServerError
			if c.Writer.Status() >= 400 {
				statusCode = c.Writer.Status()
			}

			c.JSON(statusCode, gin.H{
				"status":     "error",
				"error_code": "INTERNAL_SERVER_ERROR",
				"message":    err.Error(),
				"request_id": c.GetString("request_id"),
			})
		}
	}
}
