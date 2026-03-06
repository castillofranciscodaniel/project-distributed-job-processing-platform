package contract

import (
	"archive/zip"
	"bytes"
	"context"
	"fmt"
	"io"
	"path"
	"sync"

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

type fileResult struct {
	name    string
	content []byte
	err     error
}

func (s *ContractService) GetAllContractsZippedByClientID(ctx context.Context, clientID primitive.ObjectID) ([]byte, error) {
	contracts, err := s.repository.GetAllByClientID(ctx, clientID)
	if err != nil {
		return nil, err
	}

	buf := new(bytes.Buffer)
	zw := zip.NewWriter(buf)

	results := make(chan fileResult, len(contracts))
	var wg sync.WaitGroup

	// Phase 1: Launch Concurrent Downloads
	for _, c := range contracts {
		wg.Add(1)
		go s.downloadContractAsync(ctx, c, results, &wg)
	}

	// Phase 2: Orchestrate closure
	go func() {
		wg.Wait()
		close(results)
	}()

	// Phase 3: Sequential Zip Writing
	if err := s.writeZipArchive(zw, results); err != nil {
		return nil, err
	}

	if err := zw.Close(); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func (s *ContractService) downloadContractAsync(ctx context.Context, c Contract, results chan<- fileResult, wg *sync.WaitGroup) {
	defer wg.Done()

	rc, err := s.fileStorage.DownloadFile(ctx, c.Key)
	if err != nil {
		results <- fileResult{err: fmt.Errorf("failed to download %s: %w", c.Key, err)}
		return
	}
	defer rc.Close()

	content, err := io.ReadAll(rc)
	if err != nil {
		results <- fileResult{err: fmt.Errorf("failed to read %s: %w", c.Key, err)}
		return
	}

	results <- fileResult{
		name:    path.Base(c.Key),
		content: content,
	}
}

func (s *ContractService) writeZipArchive(zw *zip.Writer, results <-chan fileResult) error {
	for res := range results {
		if res.err != nil {
			return res.err
		}

		f, err := zw.Create(res.name)
		if err != nil {
			return fmt.Errorf("failed to create zip entry for %s: %w", res.name, err)
		}

		if _, err := f.Write(res.content); err != nil {
			return fmt.Errorf("failed to write content for %s: %w", res.name, err)
		}
	}
	return nil
}
