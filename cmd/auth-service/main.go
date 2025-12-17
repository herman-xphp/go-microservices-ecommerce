package main

import (
	"log"
	"net"
	"os"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"

	"github.com/herman-xphp/go-microservices-ecommerce/pkg/database"
	pb "github.com/herman-xphp/go-microservices-ecommerce/proto/auth"
	"github.com/herman-xphp/go-microservices-ecommerce/services/auth/domain"
	authgrpc "github.com/herman-xphp/go-microservices-ecommerce/services/auth/grpc"
	"github.com/herman-xphp/go-microservices-ecommerce/services/auth/handler"
	"github.com/herman-xphp/go-microservices-ecommerce/services/auth/repository"
	"github.com/herman-xphp/go-microservices-ecommerce/services/auth/service"
)

func main() {
	// Load configuration from environment variables
	httpPort := getEnv("HTTP_PORT", "8081")
	grpcPort := getEnv("GRPC_PORT", "9091")
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

	// Start gRPC server in a goroutine
	go startGRPCServer(grpcPort, authService)

	// Setup Gin router
	router := gin.Default()

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":    "ok",
			"service":   "auth-service",
			"http_port": httpPort,
			"grpc_port": grpcPort,
		})
	})

	// Register API routes
	api := router.Group("/api/v1")
	authHandler.RegisterRoutes(api)
	authHandler.RegisterProtectedRoutes(api)

	// Start HTTP server
	log.Printf("🚀 Auth Service HTTP starting on port %s", httpPort)
	if err := router.Run(":" + httpPort); err != nil {
		log.Fatalf("❌ Failed to start HTTP server: %v", err)
	}
}

func startGRPCServer(port string, authService service.AuthService) {
	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("❌ Failed to listen on gRPC port %s: %v", port, err)
	}

	grpcServer := grpc.NewServer()
	authGRPCServer := authgrpc.NewAuthGRPCServer(authService)
	pb.RegisterAuthServiceServer(grpcServer, authGRPCServer)

	log.Printf("🚀 Auth Service gRPC starting on port %s", port)
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
