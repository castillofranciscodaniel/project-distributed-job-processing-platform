package contract

import "context"

// Repository defines the contract for persisting Contract entities.
type Repository interface {
	Save(ctx context.Context, c *Contract) error
	GetByID(ctx context.Context, id string) (*Contract, error)
}
