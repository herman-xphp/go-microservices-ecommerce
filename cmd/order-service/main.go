package main

import (
	"os"

	"github.com/gin-gonic/gin"

	"github.com/herman-xphp/go-microservices-ecommerce/pkg/database"
	"github.com/herman-xphp/go-microservices-ecommerce/pkg/logger"
	"github.com/herman-xphp/go-microservices-ecommerce/pkg/middleware"
	"github.com/herman-xphp/go-microservices-ecommerce/services/order/client"
	"github.com/herman-xphp/go-microservices-ecommerce/services/order/domain"
	"github.com/herman-xphp/go-microservices-ecommerce/services/order/handler"
	"github.com/herman-xphp/go-microservices-ecommerce/services/order/repository"
	"github.com/herman-xphp/go-microservices-ecommerce/services/order/service"
)

const serviceName = "order-service"

func main() {
	log := logger.WithService(serviceName)

	// Load configuration from environment variables
	httpPort := getEnv("HTTP_PORT", "8083")
	productServiceAddr := getEnv("PRODUCT_SERVICE_ADDR", "localhost:9092")

	dbConfig := database.Config{
		Host:     getEnv("DB_HOST", "localhost"),
		Port:     getEnv("DB_PORT", "5432"),
		User:     getEnv("DB_USER", "user"),
		Password: getEnv("DB_PASSWORD", "password"),
		DBName:   getEnv("DB_NAME", "goshop_order"),
		SSLMode:  getEnv("DB_SSLMODE", "disable"),
	}

	// Initialize database connection
	db, err := database.NewPostgresConnection(dbConfig)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to connect to database")
	}

	// Auto-migrate database schema
	if err := db.AutoMigrate(&domain.Order{}, &domain.OrderItem{}); err != nil {
		log.Fatal().Err(err).Msg("Failed to migrate database")
	}
	log.Info().Msg("Database migrated successfully")

	// Initialize gRPC client to Product Service
	productClient, err := client.NewProductClient(productServiceAddr)
	if err != nil {
		log.Fatal().Err(err).Str("addr", productServiceAddr).Msg("Failed to connect to Product Service")
	}
	defer productClient.Close()
	log.Info().Str("addr", productServiceAddr).Msg("Connected to Product Service")

	// Initialize layers (Dependency Injection)
	orderRepo := repository.NewOrderRepository(db)
	orderService := service.NewOrderService(orderRepo, productClient)
	orderHandler := handler.NewOrderHandler(orderService)

	// Setup Gin router with middleware
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
			"status":               "ok",
			"service":              serviceName,
			"http_port":            httpPort,
			"product_service_addr": productServiceAddr,
		})
	})

	// Register API routes
	api := router.Group("/api/v1")
	orderHandler.RegisterRoutes(api)

	// Start HTTP server
	log.Info().Str("port", httpPort).Msg("Order Service HTTP starting")
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
