package main

import (
	"os"

	"github.com/gin-gonic/gin"

	"github.com/herman-xphp/go-microservices-ecommerce/pkg/database"
	"github.com/herman-xphp/go-microservices-ecommerce/pkg/logger"
	"github.com/herman-xphp/go-microservices-ecommerce/pkg/middleware"
	"github.com/herman-xphp/go-microservices-ecommerce/services/notification/domain"
	"github.com/herman-xphp/go-microservices-ecommerce/services/notification/handler"
	"github.com/herman-xphp/go-microservices-ecommerce/services/notification/service"
)

const serviceName = "notification-service"

func main() {
	log := logger.WithService(serviceName)

	// Load configuration
	httpPort := getEnv("HTTP_PORT", "8086")

	dbConfig := database.Config{
		Host:     getEnv("DB_HOST", "localhost"),
		Port:     getEnv("DB_PORT", "5432"),
		User:     getEnv("DB_USER", "user"),
		Password: getEnv("DB_PASSWORD", "password"),
		DBName:   getEnv("DB_NAME", "goshop_notification"),
		SSLMode:  getEnv("DB_SSLMODE", "disable"),
	}

	// Initialize database connection
	db, err := database.NewPostgresConnection(dbConfig)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to connect to database")
	}

	// Auto-migrate database schema
	if err := db.AutoMigrate(&domain.Notification{}, &domain.NotificationTemplate{}); err != nil {
		log.Fatal().Err(err).Msg("Failed to migrate database")
	}
	log.Info().Msg("Database migrated successfully")

	// Initialize service (no real email/sms/push senders for now, using mock)
	notificationService := service.NewNotificationService(db, nil, nil, nil)
	notificationHandler := handler.NewNotificationHandler(notificationService)

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

	// Health check endpoints
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"service": serviceName,
			"port":    httpPort,
		})
	})

	router.GET("/health/live", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "up"})
	})

	router.GET("/health/ready", func(c *gin.Context) {
		sqlDB, err := db.DB()
		if err != nil || sqlDB.Ping() != nil {
			c.JSON(503, gin.H{"status": "down", "message": "Database unavailable"})
			return
		}
		c.JSON(200, gin.H{"status": "up"})
	})

	// Register API routes
	api := router.Group("/api/v1")
	notificationHandler.RegisterRoutes(api)

	// Start HTTP server
	log.Info().Str("port", httpPort).Msg("Notification Service HTTP starting")
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
