package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"

	"github.com/herman-xphp/go-microservices-ecommerce/pkg/database"
	"github.com/herman-xphp/go-microservices-ecommerce/services/order/client"
	"github.com/herman-xphp/go-microservices-ecommerce/services/order/domain"
	"github.com/herman-xphp/go-microservices-ecommerce/services/order/handler"
	"github.com/herman-xphp/go-microservices-ecommerce/services/order/repository"
	"github.com/herman-xphp/go-microservices-ecommerce/services/order/service"
)

func main() {
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
		log.Fatalf("❌ Failed to connect to database: %v", err)
	}

	// Auto-migrate database schema
	if err := db.AutoMigrate(&domain.Order{}, &domain.OrderItem{}); err != nil {
		log.Fatalf("❌ Failed to migrate database: %v", err)
	}
	log.Println("✅ Database migrated successfully")

	// Initialize gRPC client to Product Service
	productClient, err := client.NewProductClient(productServiceAddr)
	if err != nil {
		log.Fatalf("❌ Failed to connect to Product Service: %v", err)
	}
	defer productClient.Close()

	// Initialize layers (Dependency Injection)
	orderRepo := repository.NewOrderRepository(db)
	orderService := service.NewOrderService(orderRepo, productClient)
	orderHandler := handler.NewOrderHandler(orderService)

	// Setup Gin router
	router := gin.Default()

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":               "ok",
			"service":              "order-service",
			"http_port":            httpPort,
			"product_service_addr": productServiceAddr,
		})
	})

	// Register API routes
	api := router.Group("/api/v1")
	orderHandler.RegisterRoutes(api)

	// Start HTTP server
	log.Printf("🚀 Order Service HTTP starting on port %s", httpPort)
	log.Printf("📡 Connected to Product Service at %s", productServiceAddr)
	if err := router.Run(":" + httpPort); err != nil {
		log.Fatalf("❌ Failed to start HTTP server: %v", err)
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
