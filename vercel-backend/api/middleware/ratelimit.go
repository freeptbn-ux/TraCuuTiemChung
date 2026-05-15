package middleware

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

// RateLimit middleware uses Redis Fixed Window algorithm
func RateLimit(rdb *redis.Client, limit int, window time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		if rdb == nil {
			c.Next()
			return
		}

		ctx := context.Background()
		key := fmt.Sprintf("ratelimit:%s", c.ClientIP())

		// Fixed window implementation
		count, err := rdb.Get(ctx, key).Int()
		if err != nil && err != redis.Nil {
			// On Redis error, we allow the request to proceed (fail open)
			// or we could fail closed. Standard is usually fail open for ratelimiting.
			c.Next()
			return
		}

		if count >= limit {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"status":     "error",
				"error_code": "RATE_LIMIT_EXCEEDED",
				"message":    "Too many requests, please try again later",
				"request_id": c.GetString("request_id"),
			})
			return
		}

		// Increment and set TTL if new key
		pipe := rdb.Pipeline()
		pipe.Incr(ctx, key)
		if err == redis.Nil {
			pipe.Expire(ctx, key, window)
		}
		_, err = pipe.Exec(ctx)
		if err != nil {
			// Log error but proceed
		}

		c.Next()
	}
}
