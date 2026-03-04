package contract

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// ContractStatus defines the allowed states for a Contract.
type ContractStatus string

const (
	StatusPending    ContractStatus = "PENDING"
	StatusAccepted   ContractStatus = "ACCEPTED"
	StatusProcessing ContractStatus = "PROCESSING"
	StatusCompleted  ContractStatus = "COMPLETED"
	StatusFailed     ContractStatus = "FAILED"
)

type Contract struct {
	ID        primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	ClientID  primitive.ObjectID `json:"clientId" bson:"client_id"`
	Status    ContractStatus     `json:"status" bson:"status"`
	Key       string             `json:"key" bson:"key"`
	Bucket    string             `json:"bucket" bson:"bucket"`
	CreatedAt time.Time          `json:"createdAt" bson:"created_at"`
	UpdatedAt time.Time          `json:"updatedAt" bson:"updated_at"`
	URL       string             `json:"url" bson:"-"`
}

func NewContract(clientID primitive.ObjectID, key string, bucket string, status ContractStatus) *Contract {
	now := time.Now()
	return &Contract{
		ID:        primitive.NewObjectID(),
		ClientID:  clientID,
		Status:    status,
		Key:       key,
		Bucket:    bucket,
		CreatedAt: now,
		UpdatedAt: now,
	}
}
