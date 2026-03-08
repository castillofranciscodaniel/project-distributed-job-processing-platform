package contract

import (
	"context"
	"io"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type FileStorage interface {
	UploadFile(ctx context.Context, clientID string, fileName string, fileContent io.Reader, metadata map[string]string) (string, error)
	GetPresignedURL(ctx context.Context, key string) (string, error)
	DownloadFile(ctx context.Context, key string) (io.ReadCloser, error)
}

type Repository interface {
	Save(ctx context.Context, c *Contract) error
	GetByID(ctx context.Context, id primitive.ObjectID) (*Contract, error)
	GetAllByClientID(ctx context.Context, clientID primitive.ObjectID) ([]Contract, error)
}

type EventPublisher interface {
	Publish(ctx context.Context, message interface{}) error
}
