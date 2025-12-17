package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/username/go-microservices-ecommerce/pkg/database"
	"github.com/username/go-microservices-ecommerce/services/auth/domain"
	"github.com/username/go-microservices-ecommerce/services/auth/handler"
	"github.com/username/go-microservices-ecommerce/services/auth/repository"
	"github.com/username/go-microservices-ecommerce/services/auth/service"
)

func main() {
	// Load configuration from environment variables
	port := getEnv("PORT", "8081")
	jwtSecret := getEnv("JWT_SECRET", "your-super-secret-key-change-in-production")

	dbConfig := database.Config{
		Host:     getEnv("DB_HOST", "localhost"),
		Port:     getEnv("DB_PORT", "5432"),
		User:     getEnv("DB_USER", "user"),
		Password: getEnv("DB_PASSWORD", "password"),
		DBName:   getEnv("DB_NAME", "goshop_auth"),
		SSLMode:  getEnv("DB_SSLMODE", "disable"),
	}

	// Initialize database connection
	db, err := database.NewPostgresConnection(dbConfig)
	if err != nil {
		log.Fatalf("❌ Failed to connect to database: %v", err)
	}

	// Auto-migrate database schema
	if err := db.AutoMigrate(&domain.User{}); err != nil {
		log.Fatalf("❌ Failed to migrate database: %v", err)
	}
	log.Println("✅ Database migrated successfully")

	// Initialize layers (Dependency Injection)
	userRepo := repository.NewUserRepository(db)
	authService := service.NewAuthService(userRepo, jwtSecret)
	authHandler := handler.NewAuthHandler(authService)

	// Setup Gin router
	router := gin.Default()

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok", "service": "auth-service"})
	})

	// Register API routes
	api := router.Group("/api/v1")
	authHandler.RegisterRoutes(api)
	authHandler.RegisterProtectedRoutes(api)

	// Start server
	log.Printf("🚀 Auth Service starting on port %s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("❌ Failed to start server: %v", err)
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
