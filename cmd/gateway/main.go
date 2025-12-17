package main

import (
	"fmt"
	"os"

	"github.com/gin-gonic/gin"

	"github.com/herman-xphp/go-microservices-ecommerce/pkg/logger"
	"github.com/herman-xphp/go-microservices-ecommerce/pkg/middleware"
	"github.com/herman-xphp/go-microservices-ecommerce/services/gateway/client"
	"github.com/herman-xphp/go-microservices-ecommerce/services/gateway/handler"
)

const serviceName = "api-gateway"

func main() {
	log := logger.WithService(serviceName)

	// Load configuration
	httpPort := getEnv("HTTP_PORT", "8080")
	authServiceAddr := getEnv("AUTH_SERVICE_ADDR", "localhost:9091")
	productServiceAddr := getEnv("PRODUCT_SERVICE_ADDR", "localhost:9092")
	authServiceURL := getEnv("AUTH_SERVICE_URL", "http://localhost:8081")
	productServiceURL := getEnv("PRODUCT_SERVICE_URL", "http://localhost:8082")
	orderServiceURL := getEnv("ORDER_SERVICE_URL", "http://localhost:8083")

	// Initialize gRPC clients
	authClient, err := client.NewAuthClient(authServiceAddr)
	if err != nil {
		log.Fatal().Err(err).Str("addr", authServiceAddr).Msg("Failed to connect to Auth Service gRPC")
	}
	defer authClient.Close()
	log.Info().Str("addr", authServiceAddr).Msg("Connected to Auth Service gRPC")

	productClient, err := client.NewProductClient(productServiceAddr)
	if err != nil {
		log.Fatal().Err(err).Str("addr", productServiceAddr).Msg("Failed to connect to Product Service gRPC")
	}
	defer productClient.Close()
	log.Info().Str("addr", productServiceAddr).Msg("Connected to Product Service gRPC")

	// Initialize handlers
	gatewayHandler := handler.NewGatewayHandler(authClient, productClient)

	// Define backend services for proxy
	services := map[string]*handler.ServiceConfig{
		"auth": {
			Name:    "auth-service",
			BaseURL: authServiceURL,
		},
		"product": {
			Name:    "product-service",
			BaseURL: productServiceURL,
		},
		"order": {
			Name:    "order-service",
			BaseURL: orderServiceURL,
		},
	}
	proxyHandler := handler.NewProxyHandler(services)

	// Setup Gin router
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()

	// Apply global middleware
	router.Use(middleware.Recovery())
	router.Use(middleware.RequestID())
	router.Use(middleware.Logger(serviceName))
	router.Use(middleware.CORS())
	router.Use(middleware.SecureHeaders())
	router.Use(middleware.RateLimiter(200, 20)) // Higher rate limit for gateway

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"service": serviceName,
			"port":    httpPort,
			"backends": gin.H{
				"auth_grpc":    authServiceAddr,
				"product_grpc": productServiceAddr,
				"auth_http":    authServiceURL,
				"product_http": productServiceURL,
				"order_http":   orderServiceURL,
			},
		})
	})

	// API v1 routes
	api := router.Group("/api/v1")
	{
		// ==================== PUBLIC ROUTES ====================
		// Auth routes (proxy to auth-service)
		auth := api.Group("/auth")
		{
			auth.POST("/register", proxyHandler.Proxy("auth"))
			auth.POST("/login", proxyHandler.Proxy("auth"))
		}

		// Public product routes (optional auth)
		products := api.Group("/products")
		products.Use(handler.OptionalAuthMiddleware(authClient))
		{
			products.GET("", proxyHandler.Proxy("product"))
			products.GET("/:id", proxyHandler.Proxy("product"))
			products.GET("/:id/stock", gatewayHandler.GetProductWithStock)
		}

		// Public category routes
		categories := api.Group("/categories")
		{
			categories.GET("", proxyHandler.Proxy("product"))
		}

		// ==================== PROTECTED ROUTES ====================
		protected := api.Group("")
		protected.Use(handler.AuthMiddleware(authClient))
		{
			// User profile (gateway aggregation)
			protected.GET("/me", gatewayHandler.GetUserProfile)

			// Auth protected routes
			protected.GET("/auth/profile", proxyHandler.Proxy("auth"))

			// Product management (admin)
			protected.POST("/products", proxyHandler.Proxy("product"))
			protected.PUT("/products/:id", proxyHandler.Proxy("product"))
			protected.DELETE("/products/:id", proxyHandler.Proxy("product"))

			// Category management
			protected.POST("/categories", proxyHandler.Proxy("product"))

			// Order routes
			protected.POST("/orders", proxyHandler.Proxy("order"))
			protected.GET("/orders", proxyHandler.Proxy("order"))
			protected.GET("/orders/:id", proxyHandler.Proxy("order"))
			protected.PUT("/orders/:id/status", proxyHandler.Proxy("order"))
			protected.POST("/orders/:id/cancel", proxyHandler.Proxy("order"))
		}
	}

	// Start HTTP server
	log.Info().Str("port", httpPort).Msg("API Gateway starting")
	log.Info().Msg(fmt.Sprintf("Backends: auth=%s, product=%s, order=%s", authServiceURL, productServiceURL, orderServiceURL))

	if err := router.Run(":" + httpPort); err != nil {
		log.Fatal().Err(err).Msg("Failed to start API Gateway")
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
