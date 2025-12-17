package main

import (
	"net"
	"os"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"

	"github.com/herman-xphp/go-microservices-ecommerce/pkg/database"
	"github.com/herman-xphp/go-microservices-ecommerce/pkg/logger"
	"github.com/herman-xphp/go-microservices-ecommerce/pkg/middleware"
	pb "github.com/herman-xphp/go-microservices-ecommerce/proto/product"
	"github.com/herman-xphp/go-microservices-ecommerce/services/product/domain"
	productgrpc "github.com/herman-xphp/go-microservices-ecommerce/services/product/grpc"
	"github.com/herman-xphp/go-microservices-ecommerce/services/product/handler"
	"github.com/herman-xphp/go-microservices-ecommerce/services/product/repository"
	"github.com/herman-xphp/go-microservices-ecommerce/services/product/service"
)

const serviceName = "product-service"

func main() {
	log := logger.WithService(serviceName)

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
		log.Fatal().Err(err).Msg("Failed to connect to database")
	}

	// Auto-migrate database schema
	if err := db.AutoMigrate(&domain.Category{}, &domain.Product{}); err != nil {
		log.Fatal().Err(err).Msg("Failed to migrate database")
	}
	log.Info().Msg("Database migrated successfully")

	// Initialize layers (Dependency Injection)
	productRepo := repository.NewProductRepository(db)
	categoryRepo := repository.NewCategoryRepository(db)
	productService := service.NewProductService(productRepo, categoryRepo)
	productHandler := handler.NewProductHandler(productService)

	// Start gRPC server in a goroutine
	go startGRPCServer(grpcPort, productService)

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
			"status":    "ok",
			"service":   serviceName,
			"http_port": httpPort,
			"grpc_port": grpcPort,
		})
	})

	// Register API routes
	api := router.Group("/api/v1")
	productHandler.RegisterRoutes(api)

	// Start HTTP server
	log.Info().Str("port", httpPort).Msg("Product Service HTTP starting")
	if err := router.Run(":" + httpPort); err != nil {
		log.Fatal().Err(err).Msg("Failed to start HTTP server")
	}
}

func startGRPCServer(port string, productService service.ProductService) {
	log := logger.WithService(serviceName)

	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatal().Err(err).Str("port", port).Msg("Failed to listen on gRPC port")
	}

	grpcServer := grpc.NewServer()
	productGRPCServer := productgrpc.NewProductGRPCServer(productService)
	pb.RegisterProductServiceServer(grpcServer, productGRPCServer)

	log.Info().Str("port", port).Msg("Product Service gRPC starting")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatal().Err(err).Msg("Failed to start gRPC server")
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
