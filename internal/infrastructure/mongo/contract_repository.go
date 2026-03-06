package mongo

import (
	"context"
	"fmt"
	"time"

	"github.com/francisco/distributed-job-platform/internal/domain/contract"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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
	c.UpdatedAt = time.Now()
	res, err := r.collection.InsertOne(ctx, c)
	if err != nil {
		return fmt.Errorf("failed to save contract to mongo: %w", err)
	}

	if oid, ok := res.InsertedID.(primitive.ObjectID); ok {
		c.ID = oid
	}
	return nil
}

func (r *contractRepository) GetByID(ctx context.Context, id primitive.ObjectID) (*contract.Contract, error) {
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

func (r *contractRepository) GetAllByClientID(ctx context.Context, clientID primitive.ObjectID) ([]contract.Contract, error) {
	filter := bson.M{"client_id": clientID}

	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to get contracts from mongo: %w", err)
	}
	defer cursor.Close(ctx)

	var result []contract.Contract
	for cursor.Next(ctx) {
		var c contract.Contract
		if err := cursor.Decode(&c); err != nil {
			return nil, fmt.Errorf("failed to decode contract from mongo: %w", err)
		}
		result = append(result, c)
	}

	if err := cursor.Err(); err != nil {
		return nil, fmt.Errorf("cursor error: %w", err)
	}

	return result, nil
}
