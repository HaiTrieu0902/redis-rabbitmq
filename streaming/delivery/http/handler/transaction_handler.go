package handler

import (
	"go-api-streaming/delivery/http/middleware"
	"go-api-streaming/usecase"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type TransactionHandler struct {
	useCase usecase.TransactionUseCase
}

func NewTransactionHandler(useCase usecase.TransactionUseCase) *TransactionHandler {
	return &TransactionHandler{
		useCase: useCase,
	}
}

// CreateTransaction godoc
// @Summary Create a new transaction
// @Tags transactions
// @Accept json
// @Produce json
// @Param transaction body usecase.CreateTransactionRequest true "Transaction data"
// @Success 201 {object} entity.Transaction
// @Router /transactions [post]
func (h *TransactionHandler) CreateTransaction(c *gin.Context) {
	var req usecase.CreateTransactionRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get user ID from context (set by auth middleware)
	userID, err := middleware.GetUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "user not authenticated",
			"details": err.Error(),
		})
		return
	}

	// Validate user ID is not nil/empty
	if userID == uuid.Nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid user ID in token",
			"details": "JWT token must contain a valid user_id claim. Please ensure your authentication service generates tokens with the user_id field.",
		})
		return
	}

	req.UserID = userID

	transaction, err := h.useCase.CreateTransaction(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "transaction created successfully",
		"data":    transaction,
	})
}

// GetTransaction godoc
// @Summary Get a transaction by ID
// @Tags transactions
// @Produce json
// @Param id path string true "Transaction ID"
// @Success 200 {object} entity.Transaction
// @Router /transactions/{id} [get]
func (h *TransactionHandler) GetTransaction(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid transaction ID"})
		return
	}

	transaction, err := h.useCase.GetTransaction(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": transaction})
}

// GetUserTransactions godoc
// @Summary Get user's transactions
// @Tags transactions
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Page size" default(10)
// @Success 200 {array} entity.Transaction
// @Router /transactions/my [get]
func (h *TransactionHandler) GetUserTransactions(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	transactions, err := h.useCase.GetUserTransactions(c.Request.Context(), userID, page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": transactions,
		"page": page,
		"page_size": pageSize,
	})
}

// UpdateTransactionStatus godoc
// @Summary Update transaction status
// @Tags transactions
// @Accept json
// @Produce json
// @Param id path string true "Transaction ID"
// @Param status body object true "Status update"
// @Success 200 {object} entity.Transaction
// @Router /transactions/{id}/status [patch]
func (h *TransactionHandler) UpdateTransactionStatus(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid transaction ID"})
		return
	}

	var req struct {
		Status string `json:"status" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	transaction, err := h.useCase.UpdateTransactionStatus(c.Request.Context(), id, req.Status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "transaction status updated successfully",
		"data":    transaction,
	})
}

// GetAllTransactions godoc
// @Summary Get all transactions (admin)
// @Tags transactions
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Page size" default(10)
// @Success 200 {array} entity.Transaction
// @Router /transactions [get]
func (h *TransactionHandler) GetAllTransactions(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	transactions, err := h.useCase.GetAllTransactions(c.Request.Context(), page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": transactions,
		"page": page,
		"page_size": pageSize,
	})
}

// GetTransactionsByStatus godoc
// @Summary Get transactions by status
// @Tags transactions
// @Produce json
// @Param status query string true "Transaction status"
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Page size" default(10)
// @Success 200 {array} entity.Transaction
// @Router /transactions/status [get]
func (h *TransactionHandler) GetTransactionsByStatus(c *gin.Context) {
	status := c.Query("status")
	if status == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "status parameter is required"})
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	transactions, err := h.useCase.GetTransactionsByStatus(c.Request.Context(), status, page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": transactions,
		"status": status,
		"page": page,
		"page_size": pageSize,
	})
}
