package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/francisco/distributed-job-platform/internal/domain/contract"
	"github.com/francisco/distributed-job-platform/internal/handlers"
	"github.com/francisco/distributed-job-platform/internal/infrastructure/mongo"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	// Initialize MongoDB connection
	mongoURI := "mongodb://admin:password@localhost:27017"
	client, err := mongo.ConnectDB(mongoURI)
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

	// =========================================================================
	// Dependency Injection (Manual DI)
	// =========================================================================
	db := client.Database("distributed_jobs_db")
	contractRepo := mongo.NewContractRepository(db)
	contractService := contract.NewContractService(contractRepo)
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

	// Register the POST /api/v1/contracts endpoint.
	r.Post("/api/v1/contracts", contractHandler.Upload)

	// Start the server on port 8080
	port := ":8080"
	log.Printf("Starting HTTP server on port %s", port)

	err = http.ListenAndServe(port, r)
	if err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
