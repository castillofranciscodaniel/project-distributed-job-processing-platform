package client

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Client struct {
	ID         primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Name       string             `json:"name" bson:"name"`
	LastName   string             `json:"lastName" bson:"last_name"`
	Identifier string             `json:"identifier" bson:"identifier"` // e.g., TAX ID, DNI, etc.
	CreatedAt  time.Time          `json:"createdAt" bson:"created_at"`
	UpdatedAt  time.Time          `json:"updatedAt" bson:"updated_at"`
}

func NewClient(name, lastName, identifier string) *Client {
	now := time.Now()
	return &Client{
		ID:         primitive.NewObjectID(),
		Name:       name,
		LastName:   lastName,
		Identifier: identifier,
		CreatedAt:  now,
		UpdatedAt:  now,
	}
}
