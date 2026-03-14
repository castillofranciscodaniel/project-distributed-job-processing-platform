package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/francisco/distributed-job-platform/internal/config"
	"github.com/francisco/distributed-job-platform/internal/domain/contract"
	"github.com/francisco/distributed-job-platform/internal/infrastructure/aws"
	"github.com/francisco/distributed-job-platform/internal/infrastructure/mongo"
	"github.com/francisco/distributed-job-platform/internal/infrastructure/s3"
)

func main() {
	log.Println("Starting Distributed Job Worker...")

	cfg := config.Load()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Handle signals for graceful shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	// Dependency Injection
	client, err := mongo.ConnectDB(cfg.MongoURI)
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	db := client.Database("distributed_jobs_db")

	s3Storage, err := s3.NewS3Storage(ctx, cfg.S3.ContractBucket, cfg.S3.Region, cfg.S3.PresignedURLExpirationMinutes)
	if err != nil {
		log.Fatalf("Failed to initialize S3 storage: %v", err)
	}

	sqsConsumer, err := aws.NewSQSConsumer(ctx, cfg.S3.Region, cfg.SQS.QueueURL)
	if err != nil {
		log.Fatalf("Failed to initialize SQS consumer: %v", err)
	}

	clientRepo := mongo.NewClientRepository(db)
	contractRepo := mongo.NewContractRepository(db)
	workerService := contract.NewContractWorkerService(contractRepo, clientRepo, s3Storage)

	// New Contract Consumer component
	contractConsumer := aws.NewContractZipConsumer(sqsConsumer, workerService)

	// Run consumer in a separate goroutine
	go contractConsumer.Start(ctx)

	log.Println("Worker is running. Press Ctrl+C to stop.")
	<-stop
	log.Println("Shutting down worker...")
	cancel()
}
