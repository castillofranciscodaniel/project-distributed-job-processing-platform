package contract

import (
	"context"
	"io"

	"github.com/francisco/distributed-job-platform/internal/domain/client"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ContractService struct {
	repository       Repository
	clientRepository client.Repository
	fileStorage      FileStorage
	eventPublisher   EventPublisher
	bucket           string
}

func NewContractService(repository Repository, clientRepo client.Repository, fileStorage FileStorage, eventPublisher EventPublisher, bucket string) *ContractService {
	return &ContractService{
		repository:       repository,
		clientRepository: clientRepo,
		fileStorage:      fileStorage,
		eventPublisher:   eventPublisher,
		bucket:           bucket,
	}
}

func (s *ContractService) CreateContract(ctx context.Context, clientID primitive.ObjectID, fileName string, file io.Reader) (*Contract, error) {
	key, err := s.fileStorage.UploadFile(ctx, clientID.Hex(), fileName, file, nil)
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

type ZippingRequestedEvent struct {
	ClientID string `json:"client_id"`
}

func (s *ContractService) GetAllContractsZippedByClientID(ctx context.Context, clientID primitive.ObjectID) error {
	event := ZippingRequestedEvent{
		ClientID: clientID.Hex(),
	}

	return s.eventPublisher.Publish(ctx, event)
}
