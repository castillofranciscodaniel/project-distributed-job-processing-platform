package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/francisco/distributed-job-platform/internal/config"
	clientDomain "github.com/francisco/distributed-job-platform/internal/domain/client"
	"github.com/francisco/distributed-job-platform/internal/domain/contract"
	"github.com/francisco/distributed-job-platform/internal/handlers"
	"github.com/francisco/distributed-job-platform/internal/infrastructure/aws"
	"github.com/francisco/distributed-job-platform/internal/infrastructure/mongo"
	"github.com/francisco/distributed-job-platform/internal/infrastructure/s3"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"
)

func main() {
	cwd, _ := os.Getwd()
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: .env file not found in %s. Using system environment variables.", cwd)
	}

	cfg := config.Load()

	client, err := mongo.ConnectDB(cfg.MongoURI)
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}

	// Ensure the connection is cleaned up when the server terminates
	defer func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := client.Disconnect(ctx); err != nil {
			log.Printf("Error disconnecting from MongoDB: %v", err)
		}
	}()

	log.Printf("Internal Config Loaded - S3 Bucket: %s, Region: %s", cfg.S3.ContractBucket, cfg.S3.Region)

	s3Storage, err := s3.NewS3Storage(context.Background(), cfg.S3.ContractBucket, cfg.S3.Region, cfg.S3.PresignedURLExpirationMinutes)
	if err != nil {
		log.Fatalf("Failed to initialize S3 storage: %v", err)
	}

	// =========================================================================
	// Dependency Injection (Manual DI)
	// =========================================================================
	db := client.Database("distributed_jobs_db")

	clientRepo := mongo.NewClientRepository(db)
	clientService := clientDomain.NewClientService(clientRepo)
	clientHandler := handlers.NewClientHandler(clientService)

	snsPublisher, err := aws.NewSNSPublisher(context.Background(), cfg.S3.Region, cfg.SNS.TopicArn)
	if err != nil {
		log.Fatalf("Failed to initialize SNS publisher: %v", err)
	}

	contractRepo := mongo.NewContractRepository(db)
	contractService := contract.NewContractService(contractRepo, clientRepo, s3Storage, snsPublisher, cfg.S3.ContractBucket)
	contractHandler := handlers.NewContractHandler(contractService)

	// =========================================================================
	// Setup the Chi router
	// =========================================================================
	r := chi.NewRouter()

	// Add basic middlewares provided by Chi for standard robust behavior
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// Client routes
	r.Post("/api/v1/clients", clientHandler.Create)
	r.Get("/api/v1/clients/{id}", clientHandler.GetByID)

	// Contract routes
	r.Post("/api/v1/contracts", contractHandler.Upload)
	r.Get("/api/v1/contracts/client", contractHandler.ListByClientID)
	r.Post("/api/v1/contracts/client/zip", contractHandler.DownloadZipped)
	r.Get("/api/v1/contracts/{id}", contractHandler.GetByID)

	// Start the server
	log.Printf("Starting HTTP server on port %s", cfg.Port)

	err = http.ListenAndServe(":"+cfg.Port, r)
	if err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
