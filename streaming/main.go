package main

import (
	"fmt"
	"go-api-streaming/delivery/http/handler"
	"go-api-streaming/delivery/http/middleware"
	"go-api-streaming/delivery/http/router"
	"go-api-streaming/infrastructure/config"
	"go-api-streaming/infrastructure/database"
	"go-api-streaming/infrastructure/messaging"
	"go-api-streaming/infrastructure/repository"
	"go-api-streaming/usecase"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Set Gin mode
	gin.SetMode(cfg.Server.GinMode)

	// Initialize database
	db, err := database.NewPostgresDB(&cfg.Database)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Initialize RabbitMQ
	rabbitmq, err := messaging.NewRabbitMQClient(&cfg.RabbitMQ)
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}
	defer rabbitmq.Close()

	// Declare transaction events queue
	if err := rabbitmq.DeclareQueue("transaction_events"); err != nil {
		log.Fatalf("Failed to declare queue: %v", err)
	}

	// Initialize repositories
	transactionRepo := repository.NewTransactionRepository(db)

	// Initialize use cases
	transactionUseCase := usecase.NewTransactionUseCase(transactionRepo, rabbitmq)

	// Initialize handlers
	transactionHandler := handler.NewTransactionHandler(transactionUseCase)

	// Initialize middleware
	authMiddleware := middleware.NewAuthMiddleware(cfg.JWT.Secret)

	// Setup router
	r := router.SetupRouter(transactionHandler, authMiddleware)

	// Start server
	addr := fmt.Sprintf(":%s", cfg.Server.Port)
	fmt.Printf("🚀 Streaming service is running on port %s\n", cfg.Server.Port)
	if err := r.Run(addr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
