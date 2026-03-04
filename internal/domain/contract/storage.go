package contract

import (
	"context"
	"io"

	"github.com/google/uuid"
)

// FileStorage defines the behavior for uploading files to a storage system.
type FileStorage interface {
	UploadFile(ctx context.Context, clientID uuid.UUID, fileName string, fileContent io.Reader) (string, error)
	GetPresignedURL(ctx context.Context, key string) (string, error)
}
