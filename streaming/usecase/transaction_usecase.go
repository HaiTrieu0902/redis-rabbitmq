package usecase

import (
	"context"
	"fmt"
	"go-api-streaming/domain/entity"
	"go-api-streaming/domain/repository"
	"go-api-streaming/infrastructure/messaging"
	"time"

	"github.com/google/uuid"
)

type TransactionUseCase interface {
	CreateTransaction(ctx context.Context, req *CreateTransactionRequest) (*entity.Transaction, error)
	GetTransaction(ctx context.Context, id uuid.UUID) (*entity.Transaction, error)
	GetUserTransactions(ctx context.Context, userID uuid.UUID, page, pageSize int) ([]*entity.Transaction, error)
	UpdateTransactionStatus(ctx context.Context, id uuid.UUID, status string) (*entity.Transaction, error)
	GetAllTransactions(ctx context.Context, page, pageSize int) ([]*entity.Transaction, error)
	GetTransactionsByStatus(ctx context.Context, status string, page, pageSize int) ([]*entity.Transaction, error)
}

type transactionUseCase struct {
	repo      repository.TransactionRepository
	rabbitmq  *messaging.RabbitMQClient
	queueName string
}

type CreateTransactionRequest struct {
	UserID          uuid.UUID `json:"user_id"`
	Amount          float64   `json:"amount"`
	Currency        string    `json:"currency"`
	TransactionType string    `json:"transaction_type"`
	Description     *string   `json:"description,omitempty"`
}

func NewTransactionUseCase(repo repository.TransactionRepository, rabbitmq *messaging.RabbitMQClient) TransactionUseCase {
	return &transactionUseCase{
		repo:      repo,
		rabbitmq:  rabbitmq,
		queueName: "transaction_events",
	}
}

func (u *transactionUseCase) CreateTransaction(ctx context.Context, req *CreateTransactionRequest) (*entity.Transaction, error) {
	// Validate request
	if err := u.validateCreateRequest(req); err != nil {
		return nil, err
	}

	// Validate user exists
	if err := u.validateUserExists(ctx, req.UserID); err != nil {
		return nil, err
	}

	transaction := &entity.Transaction{
		ID:              uuid.New(),
		UserID:          req.UserID,
		Amount:          req.Amount,
		Currency:        req.Currency,
		TransactionType: req.TransactionType,
		Status:          entity.TransactionStatusPending,
		Description:     req.Description,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	// Save to database
	if err := u.repo.Create(ctx, transaction); err != nil {
		return nil, fmt.Errorf("failed to create transaction: %w", err)
	}

	// Publish event to RabbitMQ
	u.publishTransactionEvent(transaction, "transaction.created")

	return transaction, nil
}

func (u *transactionUseCase) GetTransaction(ctx context.Context, id uuid.UUID) (*entity.Transaction, error) {
	return u.repo.GetByID(ctx, id)
}

func (u *transactionUseCase) GetUserTransactions(ctx context.Context, userID uuid.UUID, page, pageSize int) ([]*entity.Transaction, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	offset := (page - 1) * pageSize
	return u.repo.GetByUserID(ctx, userID, pageSize, offset)
}

func (u *transactionUseCase) UpdateTransactionStatus(ctx context.Context, id uuid.UUID, status string) (*entity.Transaction, error) {
	// Validate status
	if !u.isValidStatus(status) {
		return nil, fmt.Errorf("invalid status: %s", status)
	}

	// Get existing transaction
	transaction, err := u.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Update status
	transaction.Status = status
	transaction.UpdatedAt = time.Now()

	if err := u.repo.Update(ctx, transaction); err != nil {
		return nil, fmt.Errorf("failed to update transaction: %w", err)
	}

	// Publish event to RabbitMQ
	u.publishTransactionEvent(transaction, "transaction.updated")

	return transaction, nil
}

func (u *transactionUseCase) GetAllTransactions(ctx context.Context, page, pageSize int) ([]*entity.Transaction, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	offset := (page - 1) * pageSize
	return u.repo.GetAll(ctx, pageSize, offset)
}

func (u *transactionUseCase) GetTransactionsByStatus(ctx context.Context, status string, page, pageSize int) ([]*entity.Transaction, error) {
	if !u.isValidStatus(status) {
		return nil, fmt.Errorf("invalid status: %s", status)
	}

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	offset := (page - 1) * pageSize
	return u.repo.GetByStatus(ctx, status, pageSize, offset)
}

func (u *transactionUseCase) validateCreateRequest(req *CreateTransactionRequest) error {
	if req.Amount <= 0 {
		return fmt.Errorf("amount must be greater than 0")
	}

	if req.Currency == "" {
		return fmt.Errorf("currency is required")
	}

	if !u.isValidTransactionType(req.TransactionType) {
		return fmt.Errorf("invalid transaction type: %s", req.TransactionType)
	}

	return nil
}

func (u *transactionUseCase) validateUserExists(ctx context.Context, userID uuid.UUID) error {
	exists, err := u.repo.UserExists(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to validate user: %w", err)
	}
	
	if !exists {
		return fmt.Errorf("user with ID %s does not exist, please ensure the user is created before creating a transaction", userID)
	}
	
	return nil
}

func (u *transactionUseCase) isValidTransactionType(transactionType string) bool {
	return transactionType == entity.TransactionTypeDeposit ||
		transactionType == entity.TransactionTypeWithdraw ||
		transactionType == entity.TransactionTypePurchase
}

func (u *transactionUseCase) isValidStatus(status string) bool {
	return status == entity.TransactionStatusPending ||
		status == entity.TransactionStatusSuccess ||
		status == entity.TransactionStatusFailed
}

func (u *transactionUseCase) publishTransactionEvent(transaction *entity.Transaction, eventType string) {
	event := map[string]interface{}{
		"event_type":  eventType,
		"transaction": transaction,
		"timestamp":   time.Now(),
	}

	if err := u.rabbitmq.PublishMessage(u.queueName, event); err != nil {
		// Log error but don't fail the operation
		fmt.Printf("Warning: failed to publish event: %v\n", err)
	}
}
