package contract

import (
	"time"
)

// Contract represents the core business entity for a background task involving contracts.
type Contract struct {
	ID        string    `json:"id" bson:"_id,omitempty"`
	ClientID  string    `json:"client_id" bson:"client_id"`
	Status    string    `json:"status" bson:"status"`
	FileURL   string    `json:"file_url" bson:"file_url"`
	CreatedAt time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time `json:"updated_at" bson:"updated_at"`
}
