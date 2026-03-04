package handlers

import (
	"log"
	"net/http"

	"github.com/francisco/distributed-job-platform/internal/domain/contract"
	"github.com/francisco/distributed-job-platform/internal/handlers/helpers"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type ContractResponse struct {
	ContractID string `json:"contract_id"`
	URL        string `json:"url"`
	Status     string `json:"status"`
	Message    string `json:"message"`
}

type ContractHandler struct {
	service *contract.ContractService
}

func NewContractHandler(service *contract.ContractService) *ContractHandler {
	return &ContractHandler{
		service: service,
	}
}

func (h *ContractHandler) Upload(w http.ResponseWriter, r *http.Request) {
	log.Printf("Received request to upload a contract file")

	// Pre-flight setup: Parse the multipart form data (limit to 10 MB for this example)
	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		log.Printf("Error parsing form: %v", err)
		helpers.RespondWithError(w, http.StatusBadRequest, "Invalid input or file too large")
		return
	}

	file, handler, err := r.FormFile("file")
	if err != nil {
		log.Printf("Error retrieving file 'file': %v", err)
		helpers.RespondWithError(w, http.StatusBadRequest, "Missing 'file' field in multipart form")
		return
	}
	defer file.Close()

	clientIDStr := r.Header.Get("client_id")
	if clientIDStr == "" {
		helpers.RespondWithError(w, http.StatusBadRequest, "Missing 'client_id' header")
		return
	}

	clientID, err := uuid.Parse(clientIDStr)
	if err != nil {
		helpers.RespondWithError(w, http.StatusBadRequest, "Invalid 'client_id' format. Must be a UUID")
		return
	}

	log.Printf("File received: Filename: %s, Size: %d bytes. ClientID: %v", handler.Filename, handler.Size, clientID)

	c, err := h.service.CreateContract(r.Context(), clientID, handler.Filename, file)
	if err != nil {
		log.Printf("Error creating contract: %v", err)
		helpers.RespondWithError(w, http.StatusInternalServerError, "Failed to process contract")
		return
	}

	response := ContractResponse{
		ContractID: c.ID.String(),
		Status:     string(c.Status),
		Message:    "File uploaded and processing started",
	}

	helpers.RespondWithJSON(w, http.StatusAccepted, response)
}

func (h *ContractHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	contractID := chi.URLParam(r, "id")
	contract, err := h.service.GetContractByID(r.Context(), uuid.Must(uuid.Parse(contractID)))

	if err != nil {
		helpers.RespondWithError(w, http.StatusNotFound, "Contract not found")
		return
	}

	response := ContractResponse{
		ContractID: contract.ID.String(),
		URL:        contract.URL,
		Status:     string(contract.Status),
		Message:    "File uploaded and processing started",
	}

	helpers.RespondWithJSON(w, http.StatusOK, response)
}
