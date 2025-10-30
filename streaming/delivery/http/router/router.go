package router

import (
	"go-api-streaming/delivery/http/handler"
	"go-api-streaming/delivery/http/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRouter(
	transactionHandler *handler.TransactionHandler,
	authMiddleware *middleware.AuthMiddleware,
) *gin.Engine {
	router := gin.Default()

	// CORS middleware
	router.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, PATCH, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"service": "streaming",
		})
	})

	// API routes
	api := router.Group("/api/v1")
	{
		// Transaction routes (protected)
		transactions := api.Group("/transactions")
		transactions.Use(authMiddleware.Authenticate())
		{
			transactions.POST("", transactionHandler.CreateTransaction)
			transactions.GET("/:id", transactionHandler.GetTransaction)
			transactions.GET("/my", transactionHandler.GetUserTransactions)
			transactions.PATCH("/:id/status", transactionHandler.UpdateTransactionStatus)
			transactions.GET("", transactionHandler.GetAllTransactions)
			transactions.GET("/status", transactionHandler.GetTransactionsByStatus)
		}
	}

	return router
}
