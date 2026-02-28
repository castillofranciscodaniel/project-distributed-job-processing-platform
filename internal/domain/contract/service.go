package contract

import "context"

type ContractService struct {
	repository Repository
}

func NewContractService(repository Repository) *ContractService {
	return &ContractService{
		repository: repository,
	}
}

func (s *ContractService) CreateContract(ctx context.Context, c *Contract) error {
	return s.repository.Save(ctx, c)
}

func (s *ContractService) GetContractByID(ctx context.Context, id string) (*Contract, error) {
	return s.repository.GetByID(ctx, id)
}
