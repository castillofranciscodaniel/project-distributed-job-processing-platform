package s3

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/s3/transfermanager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type Storage struct {
	client            *s3.Client
	transferManager   *transfermanager.Client
	presignClient     *s3.PresignClient
	bucket            string
	expirationMinutes int
}

func NewS3Storage(ctx context.Context, bucketName, region string, expirationMinutes int) (*Storage, error) {
	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion(region))
	if err != nil {
		return nil, fmt.Errorf("unable to load SDK config: %w", err)
	}

	client := s3.NewFromConfig(cfg)
	tmClient := transfermanager.New(client)
	presignClient := s3.NewPresignClient(client)

	return &Storage{
		client:            client,
		transferManager:   tmClient,
		presignClient:     presignClient,
		bucket:            bucketName,
		expirationMinutes: expirationMinutes,
	}, nil
}

func (s *Storage) UploadFile(ctx context.Context, clientID string, fileName string, fileContent io.Reader) (string, error) {
	key := fmt.Sprintf("%s/%s", clientID, fileName)

	result, err := s.transferManager.UploadObject(ctx, &transfermanager.UploadObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(key),
		Body:   fileContent,
	})

	if err != nil {
		return "", fmt.Errorf("failed to upload file to S3: %w", err)
	}

	log.Printf("Successfully uploaded file to S3: %s", *result.Location)
	return *result.Key, nil
}

func (s *Storage) GetPresignedURL(ctx context.Context, key string) (string, error) {
	request, err := s.presignClient.PresignGetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(key),
	}, func(o *s3.PresignOptions) {
		o.Expires = time.Duration(s.expirationMinutes) * time.Minute
	})
	if err != nil {
		return "", fmt.Errorf("failed to generate presigned URL: %w", err)
	}

	return request.URL, nil
}
