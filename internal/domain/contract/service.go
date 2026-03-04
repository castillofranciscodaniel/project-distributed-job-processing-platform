package contract

import (
	"context"
	"io"

	"github.com/google/uuid"
)

type ContractService struct {
	repository  Repository
	fileStorage FileStorage
	bucket      string
}

func NewContractService(repository Repository, fileStorage FileStorage, bucket string) *ContractService {
	return &ContractService{
		repository:  repository,
		fileStorage: fileStorage,
		bucket:      bucket,
	}
}

func (s *ContractService) CreateContract(ctx context.Context, clientID uuid.UUID, fileName string, file io.Reader) (*Contract, error) {
	key, err := s.fileStorage.UploadFile(ctx, clientID, fileName, file)
	if err != nil {
		return nil, err
	}

	c := NewContract(clientID, key, s.bucket, StatusAccepted)

	if err := s.repository.Save(ctx, c); err != nil {
		return nil, err
	}

	return c, nil
}

func (s *ContractService) GetContractByID(ctx context.Context, id uuid.UUID) (*Contract, error) {
	contract, err := s.repository.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	url, err := s.fileStorage.GetPresignedURL(ctx, contract.Key)
	if err != nil {
		return nil, err
	}

	contract.URL = url
	return contract, nil
}
