package mongo

import (
	"context"
	"fmt"

	"github.com/francisco/distributed-job-platform/internal/domain/contract"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type contractRepository struct {
	collection *mongo.Collection
}

func NewContractRepository(db *mongo.Database) contract.Repository {
	return &contractRepository{
		collection: db.Collection("contracts"),
	}
}

func (r *contractRepository) Save(ctx context.Context, c *contract.Contract) error {
	_, err := r.collection.InsertOne(ctx, c)
	if err != nil {
		return fmt.Errorf("failed to save contract to mongo: %w", err)
	}
	return nil
}

func (r *contractRepository) GetByID(ctx context.Context, id string) (*contract.Contract, error) {
	var result contract.Contract
	filter := bson.M{"_id": id}

	err := r.collection.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("contract not found in database")
		}
		return nil, fmt.Errorf("failed to get contract from mongo: %w", err)
	}
	return &result, nil
}
