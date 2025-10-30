package repository

import (
	"context"
	"database/sql"
	"fmt"
	"go-api-streaming/domain/entity"
	"go-api-streaming/domain/repository"

	"github.com/google/uuid"
)

type transactionRepositoryImpl struct {
	db *sql.DB
}

func NewTransactionRepository(db *sql.DB) repository.TransactionRepository {
	return &transactionRepositoryImpl{
		db: db,
	}
}

func (r *transactionRepositoryImpl) Create(ctx context.Context, transaction *entity.Transaction) error {
	query := `
		INSERT INTO transactions (id, user_id, amount, currency, transaction_type, status, description, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`

	_, err := r.db.ExecContext(
		ctx,
		query,
		transaction.ID,
		transaction.UserID,
		transaction.Amount,
		transaction.Currency,
		transaction.TransactionType,
		transaction.Status,
		transaction.Description,
		transaction.CreatedAt,
		transaction.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to create transaction: %w", err)
	}

	return nil
}

func (r *transactionRepositoryImpl) GetByID(ctx context.Context, id uuid.UUID) (*entity.Transaction, error) {
	query := `
		SELECT id, user_id, amount, currency, transaction_type, status, description, created_at, updated_at
		FROM transactions
		WHERE id = $1
	`

	transaction := &entity.Transaction{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&transaction.ID,
		&transaction.UserID,
		&transaction.Amount,
		&transaction.Currency,
		&transaction.TransactionType,
		&transaction.Status,
		&transaction.Description,
		&transaction.CreatedAt,
		&transaction.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("transaction not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get transaction: %w", err)
	}

	return transaction, nil
}

func (r *transactionRepositoryImpl) GetByUserID(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*entity.Transaction, error) {
	query := `
		SELECT id, user_id, amount, currency, transaction_type, status, description, created_at, updated_at
		FROM transactions
		WHERE user_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.db.QueryContext(ctx, query, userID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get transactions: %w", err)
	}
	defer rows.Close()

	return r.scanTransactions(rows)
}

func (r *transactionRepositoryImpl) Update(ctx context.Context, transaction *entity.Transaction) error {
	query := `
		UPDATE transactions
		SET amount = $1, currency = $2, transaction_type = $3, status = $4, description = $5, updated_at = $6
		WHERE id = $7
	`

	result, err := r.db.ExecContext(
		ctx,
		query,
		transaction.Amount,
		transaction.Currency,
		transaction.TransactionType,
		transaction.Status,
		transaction.Description,
		transaction.UpdatedAt,
		transaction.ID,
	)

	if err != nil {
		return fmt.Errorf("failed to update transaction: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("transaction not found")
	}

	return nil
}

func (r *transactionRepositoryImpl) GetAll(ctx context.Context, limit, offset int) ([]*entity.Transaction, error) {
	query := `
		SELECT id, user_id, amount, currency, transaction_type, status, description, created_at, updated_at
		FROM transactions
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := r.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get transactions: %w", err)
	}
	defer rows.Close()

	return r.scanTransactions(rows)
}

func (r *transactionRepositoryImpl) GetByStatus(ctx context.Context, status string, limit, offset int) ([]*entity.Transaction, error) {
	query := `
		SELECT id, user_id, amount, currency, transaction_type, status, description, created_at, updated_at
		FROM transactions
		WHERE status = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.db.QueryContext(ctx, query, status, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get transactions: %w", err)
	}
	defer rows.Close()

	return r.scanTransactions(rows)
}

func (r *transactionRepositoryImpl) UserExists(ctx context.Context, userID uuid.UUID) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM users WHERE id = $1)`
	var exists bool
	
	err := r.db.QueryRowContext(ctx, query, userID).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check user existence: %w", err)
	}
	
	return exists, nil
}

func (r *transactionRepositoryImpl) scanTransactions(rows *sql.Rows) ([]*entity.Transaction, error) {
	var transactions []*entity.Transaction

	for rows.Next() {
		transaction := &entity.Transaction{}
		err := rows.Scan(
			&transaction.ID,
			&transaction.UserID,
			&transaction.Amount,
			&transaction.Currency,
			&transaction.TransactionType,
			&transaction.Status,
			&transaction.Description,
			&transaction.CreatedAt,
			&transaction.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan transaction: %w", err)
		}
		transactions = append(transactions, transaction)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %w", err)
	}

	return transactions, nil
}
