package contract

import (
	"context"
	"io"

	"go.mongodb.org/mongo-driver/bson/primitive"
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

func (s *ContractService) CreateContract(ctx context.Context, clientID primitive.ObjectID, fileName string, file io.Reader) (*Contract, error) {
	key, err := s.fileStorage.UploadFile(ctx, clientID.Hex(), fileName, file)
	if err != nil {
		return nil, err
	}

	c := NewContract(clientID, key, s.bucket, StatusAccepted)

	if err := s.repository.Save(ctx, c); err != nil {
		return nil, err
	}

	return c, nil
}

func (s *ContractService) GetContractByID(ctx context.Context, id primitive.ObjectID) (*Contract, error) {
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

func (s *ContractService) GetAllContractsByClientID(ctx context.Context, clientID primitive.ObjectID) ([]Contract, error) {
	contracts, err := s.repository.GetAllByClientID(ctx, clientID)
	if err != nil {
		return nil, err
	}

	for i := range contracts {
		url, err := s.fileStorage.GetPresignedURL(ctx, contracts[i].Key)
		if err == nil {
			contracts[i].URL = url
		}
	}

	return contracts, nil
}
