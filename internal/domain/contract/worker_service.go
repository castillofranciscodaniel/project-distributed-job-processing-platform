package contract

import (
	"archive/zip"
	"bytes"
	"context"
	"fmt"
	"io"
	"path"
	"sync"

	"github.com/francisco/distributed-job-platform/internal/domain/client"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ContractWorkerService struct {
	repository       Repository
	clientRepository client.Repository
	fileStorage      FileStorage
}

func NewContractWorkerService(repository Repository, clientRepo client.Repository, fileStorage FileStorage) *ContractWorkerService {
	return &ContractWorkerService{
		repository:       repository,
		clientRepository: clientRepo,
		fileStorage:      fileStorage,
	}
}

type fileResult struct {
	name    string
	content []byte
	err     error
}

func (s *ContractWorkerService) ProcessZippingRequest(ctx context.Context, clientID primitive.ObjectID) error {
	contracts, err := s.repository.GetAllByClientID(ctx, clientID)
	if err != nil {
		return err
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
		return err
	}

	if err := zw.Close(); err != nil {
		return err
	}

	// Final Phase: Upload to S3 in the client folder
	// We add metadata so Lambda knows where to send the email
	metadata := make(map[string]string)
	metadata["client-id"] = clientID.Hex()

	// Fetch client to get email
	cl, err := s.clientRepository.GetByID(ctx, clientID)
	if err == nil && cl.Email != "" {
		metadata["client-email"] = cl.Email
	}

	zipFileName := fmt.Sprintf("contracts_%s.zip", clientID.Hex())
	_, err = s.fileStorage.UploadFile(ctx, clientID.Hex(), zipFileName, buf, metadata)
	if err != nil {
		return fmt.Errorf("failed to upload final zip to S3: %w", err)
	}

	return nil
}

func (s *ContractWorkerService) downloadContractAsync(ctx context.Context, c Contract, results chan<- fileResult, wg *sync.WaitGroup) {
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

func (s *ContractWorkerService) writeZipArchive(zw *zip.Writer, results <-chan fileResult) error {
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
