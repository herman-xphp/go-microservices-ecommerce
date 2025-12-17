package main

import (
	"log"
	"net"
	"os"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"

	"github.com/herman-xphp/go-microservices-ecommerce/pkg/database"
	pb "github.com/herman-xphp/go-microservices-ecommerce/proto/product"
	"github.com/herman-xphp/go-microservices-ecommerce/services/product/domain"
	productgrpc "github.com/herman-xphp/go-microservices-ecommerce/services/product/grpc"
	"github.com/herman-xphp/go-microservices-ecommerce/services/product/handler"
	"github.com/herman-xphp/go-microservices-ecommerce/services/product/repository"
	"github.com/herman-xphp/go-microservices-ecommerce/services/product/service"
)

func main() {
	// Load configuration from environment variables
	httpPort := getEnv("HTTP_PORT", "8082")
	grpcPort := getEnv("GRPC_PORT", "9092")

	dbConfig := database.Config{
		Host:     getEnv("DB_HOST", "localhost"),
		Port:     getEnv("DB_PORT", "5432"),
		User:     getEnv("DB_USER", "user"),
		Password: getEnv("DB_PASSWORD", "password"),
		DBName:   getEnv("DB_NAME", "goshop_product"),
		SSLMode:  getEnv("DB_SSLMODE", "disable"),
	}

	// Initialize database connection
	db, err := database.NewPostgresConnection(dbConfig)
	if err != nil {
		log.Fatalf("❌ Failed to connect to database: %v", err)
	}

	// Auto-migrate database schema
	if err := db.AutoMigrate(&domain.Category{}, &domain.Product{}); err != nil {
		log.Fatalf("❌ Failed to migrate database: %v", err)
	}
	log.Println("✅ Database migrated successfully")

	// Initialize layers (Dependency Injection)
	productRepo := repository.NewProductRepository(db)
	categoryRepo := repository.NewCategoryRepository(db)
	productService := service.NewProductService(productRepo, categoryRepo)
	productHandler := handler.NewProductHandler(productService)

	// Start gRPC server in a goroutine
	go startGRPCServer(grpcPort, productService)

	// Setup Gin router
	router := gin.Default()

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":    "ok",
			"service":   "product-service",
			"http_port": httpPort,
			"grpc_port": grpcPort,
		})
	})

	// Register API routes
	api := router.Group("/api/v1")
	productHandler.RegisterRoutes(api)

	// Start HTTP server
	log.Printf("🚀 Product Service HTTP starting on port %s", httpPort)
	if err := router.Run(":" + httpPort); err != nil {
		log.Fatalf("❌ Failed to start HTTP server: %v", err)
	}
}

func startGRPCServer(port string, productService service.ProductService) {
	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("❌ Failed to listen on gRPC port %s: %v", port, err)
	}

	grpcServer := grpc.NewServer()
	productGRPCServer := productgrpc.NewProductGRPCServer(productService)
	pb.RegisterProductServiceServer(grpcServer, productGRPCServer)

	log.Printf("🚀 Product Service gRPC starting on port %s", port)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("❌ Failed to start gRPC server: %v", err)
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
