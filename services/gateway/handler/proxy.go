package handler

import (
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// ServiceConfig holds configuration for a backend service
type ServiceConfig struct {
	Name    string
	BaseURL string
}

// ProxyHandler creates a reverse proxy handler for a backend service
type ProxyHandler struct {
	httpClient *http.Client
	services   map[string]*ServiceConfig
}

// NewProxyHandler creates a new proxy handler
func NewProxyHandler(services map[string]*ServiceConfig) *ProxyHandler {
	return &ProxyHandler{
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		services: services,
	}
}

// Proxy forwards the request to the appropriate backend service
func (p *ProxyHandler) Proxy(serviceName string) gin.HandlerFunc {
	return func(c *gin.Context) {
		service, exists := p.services[serviceName]
		if !exists {
			c.JSON(http.StatusBadGateway, gin.H{
				"success": false,
				"message": "Unknown service: " + serviceName,
			})
			return
		}

		// Build target URL
		targetURL := service.BaseURL + c.Request.URL.Path
		if c.Request.URL.RawQuery != "" {
			targetURL += "?" + c.Request.URL.RawQuery
		}

		// Create proxy request
		proxyReq, err := http.NewRequestWithContext(
			c.Request.Context(),
			c.Request.Method,
			targetURL,
			c.Request.Body,
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "Failed to create proxy request",
				"error":   err.Error(),
			})
			return
		}

		// Copy headers
		for key, values := range c.Request.Header {
			for _, value := range values {
				proxyReq.Header.Add(key, value)
			}
		}

		// Add user info headers if authenticated
		if userID, exists := c.Get("user_id"); exists {
			proxyReq.Header.Set("X-User-ID", formatUint(userID.(uint)))
		}
		if userEmail, exists := c.Get("user_email"); exists {
			proxyReq.Header.Set("X-User-Email", userEmail.(string))
		}

		// Add request ID
		if requestID, exists := c.Get("request_id"); exists {
			proxyReq.Header.Set("X-Request-ID", requestID.(string))
		}

		// Execute request
		resp, err := p.httpClient.Do(proxyReq)
		if err != nil {
			c.JSON(http.StatusBadGateway, gin.H{
				"success": false,
				"message": "Failed to reach backend service",
				"error":   err.Error(),
			})
			return
		}
		defer resp.Body.Close()

		// Copy response headers
		for key, values := range resp.Header {
			for _, value := range values {
				c.Header(key, value)
			}
		}

		// Copy response body
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "Failed to read response",
			})
			return
		}

		c.Data(resp.StatusCode, resp.Header.Get("Content-Type"), body)
	}
}

func formatUint(n uint) string {
	return strconv.FormatUint(uint64(n), 10)
}
