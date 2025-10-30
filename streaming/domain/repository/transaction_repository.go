package repository

import (
	"context"
	"go-api-streaming/domain/entity"

	"github.com/google/uuid"
)

type TransactionRepository interface {
	Create(ctx context.Context, transaction *entity.Transaction) error
	GetByID(ctx context.Context, id uuid.UUID) (*entity.Transaction, error)
	GetByUserID(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*entity.Transaction, error)
	Update(ctx context.Context, transaction *entity.Transaction) error
	GetAll(ctx context.Context, limit, offset int) ([]*entity.Transaction, error)
	GetByStatus(ctx context.Context, status string, limit, offset int) ([]*entity.Transaction, error)
	UserExists(ctx context.Context, userID uuid.UUID) (bool, error)
}
