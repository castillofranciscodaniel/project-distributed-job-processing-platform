package client

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Client struct {
	ID         primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Name       string             `json:"name" bson:"name"`
	LastName   string             `json:"lastName" bson:"last_name"`
	Email      string             `json:"email" bson:"email"`
	Identifier string             `json:"identifier" bson:"identifier"` // e.g., TAX ID, DNI, etc.
	CreatedAt  time.Time          `json:"createdAt" bson:"created_at"`
	UpdatedAt  time.Time          `json:"updatedAt" bson:"updated_at"`
}

func NewClient(name, lastName, email, identifier string) *Client {
	now := time.Now()
	return &Client{
		Name:       name,
		LastName:   lastName,
		Email:      email,
		Identifier: identifier,
		CreatedAt:  now,
		UpdatedAt:  now,
	}
}
