package client

import (
	"context"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Repository interface {
	Save(ctx context.Context, c *Client) error
	GetByID(ctx context.Context, id primitive.ObjectID) (*Client, error)
	GetByIdentifier(ctx context.Context, identifier string) (*Client, error)
}
