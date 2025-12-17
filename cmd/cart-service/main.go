package main

import (
	"context"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"

	"github.com/herman-xphp/go-microservices-ecommerce/pkg/logger"
	"github.com/herman-xphp/go-microservices-ecommerce/pkg/middleware"
	"github.com/herman-xphp/go-microservices-ecommerce/services/cart/client"
	"github.com/herman-xphp/go-microservices-ecommerce/services/cart/handler"
	"github.com/herman-xphp/go-microservices-ecommerce/services/cart/repository"
	"github.com/herman-xphp/go-microservices-ecommerce/services/cart/service"
)

const serviceName = "cart-service"

func main() {
	log := logger.WithService(serviceName)

	// Load configuration
	httpPort := getEnv("HTTP_PORT", "8085")
	redisAddr := getEnv("REDIS_ADDR", "localhost:6379")
	redisPassword := getEnv("REDIS_PASSWORD", "")
	productServiceAddr := getEnv("PRODUCT_SERVICE_ADDR", "localhost:9092")

	// Initialize Redis client
	redisClient := redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: redisPassword,
		DB:       0,
	})

	// Test Redis connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := redisClient.Ping(ctx).Err(); err != nil {
		log.Fatal().Err(err).Msg("Failed to connect to Redis")
	}
	log.Info().Str("addr", redisAddr).Msg("Connected to Redis")

	// Initialize Product Service gRPC client
	productClient, err := client.NewProductClient(productServiceAddr)
	if err != nil {
		log.Fatal().Err(err).Str("addr", productServiceAddr).Msg("Failed to connect to Product Service")
	}
	defer productClient.Close()
	log.Info().Str("addr", productServiceAddr).Msg("Connected to Product Service gRPC")

	// Initialize layers (Dependency Injection)
	cartRepo := repository.NewRedisCartRepository(redisClient)
	cartService := service.NewCartService(cartRepo, productClient)
	cartHandler := handler.NewCartHandler(cartService)

	// Setup Gin router
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()

	// Apply middleware
	router.Use(middleware.Recovery())
	router.Use(middleware.RequestID())
	router.Use(middleware.Logger(serviceName))
	router.Use(middleware.CORS())
	router.Use(middleware.SecureHeaders())
	router.Use(middleware.RateLimiter(100, 10))

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"service": serviceName,
			"port":    httpPort,
		})
	})

	// Liveness probe
	router.GET("/health/live", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "up"})
	})

	// Readiness probe
	router.GET("/health/ready", func(c *gin.Context) {
		ctx := c.Request.Context()
		if err := redisClient.Ping(ctx).Err(); err != nil {
			c.JSON(503, gin.H{
				"status":  "down",
				"message": "Redis unavailable",
			})
			return
		}
		c.JSON(200, gin.H{"status": "up"})
	})

	// Register API routes
	api := router.Group("/api/v1")
	cartHandler.RegisterRoutes(api)

	// Start HTTP server
	log.Info().Str("port", httpPort).Msg("Cart Service HTTP starting")
	if err := router.Run(":" + httpPort); err != nil {
		log.Fatal().Err(err).Msg("Failed to start HTTP server")
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
