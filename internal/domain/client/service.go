package client

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ClientService struct {
	repository Repository
}

func NewClientService(repo Repository) *ClientService {
	return &ClientService{
		repository: repo,
	}
}

func (s *ClientService) CreateClient(ctx context.Context, name, lastName, identifier string) (*Client, error) {
	// Check if client with identifier already exists
	existing, _ := s.repository.GetByIdentifier(ctx, identifier)
	if existing != nil {
		return nil, fmt.Errorf("client with identifier %s already exists", identifier)
	}

	c := NewClient(name, lastName, identifier)
	if err := s.repository.Save(ctx, c); err != nil {
		return nil, err
	}

	return c, nil
}

func (s *ClientService) GetClientByID(ctx context.Context, id primitive.ObjectID) (*Client, error) {
	return s.repository.GetByID(ctx, id)
}
