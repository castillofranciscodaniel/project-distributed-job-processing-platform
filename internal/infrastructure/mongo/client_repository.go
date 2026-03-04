package mongo

import (
	"context"
	"fmt"

	"github.com/francisco/distributed-job-platform/internal/domain/client"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type clientRepository struct {
	collection *mongo.Collection
}

func NewClientRepository(db *mongo.Database) client.Repository {
	return &clientRepository{
		collection: db.Collection("clients"),
	}
}

func (r *clientRepository) Save(ctx context.Context, c *client.Client) error {
	_, err := r.collection.InsertOne(ctx, c)
	if err != nil {
		return fmt.Errorf("failed to save client: %w", err)
	}
	return nil
}

func (r *clientRepository) GetByID(ctx context.Context, id primitive.ObjectID) (*client.Client, error) {
	var result client.Client
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("client not found")
		}
		return nil, err
	}
	return &result, nil
}

func (r *clientRepository) GetByIdentifier(ctx context.Context, identifier string) (*client.Client, error) {
	var result client.Client
	err := r.collection.FindOne(ctx, bson.M{"identifier": identifier}).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil // Return nil, nil if not found for the check
		}
		return nil, err
	}
	return &result, nil
}
