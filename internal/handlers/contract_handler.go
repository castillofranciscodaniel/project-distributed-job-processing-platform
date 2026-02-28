package handlers

import (
	"log"
	"net/http"

	"github.com/francisco/distributed-job-platform/internal/domain/contract"
	"github.com/francisco/distributed-job-platform/internal/handlers/helpers"
)

// ContractResponse maps to the ContractResponse schema defined in OpenAPI.
type ContractResponse struct {
	ContractID string `json:"contract_id"`
	Status     string `json:"status"`
	Message    string `json:"message"`
}

type ContractHandler struct {
	service *contract.ContractService
}

// NewContractHandler is the constructor for the handler struct.
func NewContractHandler(service *contract.ContractService) *ContractHandler {
	return &ContractHandler{
		service: service,
	}
}

// Upload handles the POST /api/v1/contracts endpoint.
func (h *ContractHandler) Upload(w http.ResponseWriter, r *http.Request) {
	log.Printf("Received request to upload a contract file")

	// Pre-flight setup: Parse the multipart form data (limit to 10 MB for this example)
	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		log.Printf("Error parsing form: %v", err)
		helpers.RespondWithError(w, http.StatusBadRequest, "Invalid input or file too large")
		return
	}

	// Retrieve the file from form data
	file, handler, err := r.FormFile("file")
	if err != nil {
		log.Printf("Error retrieving file 'file': %v", err)
		helpers.RespondWithError(w, http.StatusBadRequest, "Missing 'file' field in multipart form")
		return
	}
	defer file.Close()

	// Extract client_id from form (Simulating receiving the ClientID)
	clientID := r.Header.Get("client_id")
	if clientID == "" {
		// Just a fallback for this phase if client_id is not provided
		clientID = "unknown-client"
	}

	// Log basic metadata of the file
	log.Printf("File received: Filename: %s, Size: %d bytes, Header: %v. ClientID: %s", handler.Filename, handler.Size, handler.Header, clientID)

	// Create theoretical response
	response := ContractResponse{
		ContractID: "dummy-uuid-1234", // We'll delegate UUID and DB creation to the Service layer soon!
		Status:     "accepted",
		Message:    "File uploaded and processing started",
	}

	// Respond with 202 Accepted using our clean helper
	helpers.RespondWithJSON(w, http.StatusAccepted, response)
}
