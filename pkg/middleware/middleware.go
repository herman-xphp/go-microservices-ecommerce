package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/herman-xphp/go-microservices-ecommerce/pkg/logger"
)

const (
	// RequestIDHeader is the header name for request ID
	RequestIDHeader = "X-Request-ID"
	// RequestIDKey is the context key for request ID
	RequestIDKey = "request_id"
)

// RequestID adds a unique request ID to each request
func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Check if request ID already exists in header
		requestID := c.GetHeader(RequestIDHeader)
		if requestID == "" {
			requestID = uuid.New().String()
		}

		// Set request ID in context and response header
		c.Set(RequestIDKey, requestID)
		c.Header(RequestIDHeader, requestID)

		c.Next()
	}
}

// Logger is a middleware that logs HTTP requests
func Logger(serviceName string) gin.HandlerFunc {
	log := logger.WithService(serviceName)

	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery

		// Get request ID
		requestID, _ := c.Get(RequestIDKey)

		c.Next()

		latency := time.Since(start)
		status := c.Writer.Status()

		// Log the request
		log.Info().
			Str("request_id", requestID.(string)).
			Str("method", c.Request.Method).
			Str("path", path).
			Str("query", query).
			Int("status", status).
			Dur("latency", latency).
			Str("ip", c.ClientIP()).
			Msg("HTTP Request")
	}
}

// Recovery is a middleware that recovers from panics
func Recovery() gin.HandlerFunc {
	log := logger.Get()

	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				requestID, _ := c.Get(RequestIDKey)
				log.Error().
					Interface("error", err).
					Str("request_id", requestID.(string)).
					Msg("Panic recovered")

				c.AbortWithStatusJSON(500, gin.H{
					"success":    false,
					"message":    "Internal server error",
					"request_id": requestID,
				})
			}
		}()
		c.Next()
	}
}
