package contract

import (
	"context"

	"github.com/google/uuid"
)

type Repository interface {
	Save(ctx context.Context, c *Contract) error
	GetByID(ctx context.Context, id uuid.UUID) (*Contract, error)
}
