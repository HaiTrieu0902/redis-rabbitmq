package entity

import (
	"time"

	"github.com/google/uuid"
)

type Transaction struct {
	ID              uuid.UUID  `json:"id"`
	UserID          uuid.UUID  `json:"user_id"`
	Amount          float64    `json:"amount"`
	Currency        string     `json:"currency"`
	TransactionType string     `json:"transaction_type"`
	Status          string     `json:"status"`
	Description     *string    `json:"description,omitempty"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
}

// Transaction types
const (
	TransactionTypeDeposit  = "deposit"
	TransactionTypeWithdraw = "withdraw"
	TransactionTypePurchase = "purchase"
)

// Transaction statuses
const (
	TransactionStatusPending = "pending"
	TransactionStatusSuccess = "success"
	TransactionStatusFailed  = "failed"
)
