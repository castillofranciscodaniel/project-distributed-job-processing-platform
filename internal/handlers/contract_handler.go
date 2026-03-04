package handlers

import (
	"log"
	"net/http"

	"github.com/francisco/distributed-job-platform/internal/domain/contract"
	"github.com/francisco/distributed-job-platform/internal/handlers/helpers"
	"github.com/go-chi/chi/v5"
	"go.mongodb.org/mongo-driver/bson/primitive"
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

	clientID, err := primitive.ObjectIDFromHex(clientIDStr)
	if err != nil {
		helpers.RespondWithError(w, http.StatusBadRequest, "Invalid 'client_id' format. Must be a 24-character hex string (ObjectID)")
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
		ContractID: c.ID.Hex(),
		URL:        "",
		Status:     string(c.Status),
		Message:    "File uploaded and processing started",
	}

	helpers.RespondWithJSON(w, http.StatusAccepted, response)
}

func (h *ContractHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	contractIDStr := chi.URLParam(r, "id")
	contractID, err := primitive.ObjectIDFromHex(contractIDStr)
	if err != nil {
		helpers.RespondWithError(w, http.StatusBadRequest, "Invalid 'id' format")
		return
	}

	contract, err := h.service.GetContractByID(r.Context(), contractID)
	if err != nil {
		helpers.RespondWithError(w, http.StatusNotFound, "Contract not found")
		return
	}

	response := ContractResponse{
		ContractID: contract.ID.Hex(),
		URL:        contract.URL,
		Status:     string(contract.Status),
		Message:    "File retrieved successfully",
	}

	helpers.RespondWithJSON(w, http.StatusOK, response)
}

func (h *ContractHandler) ListByClientID(w http.ResponseWriter, r *http.Request) {
	clientIDStr := chi.URLParam(r, "clientID")
	clientID, err := primitive.ObjectIDFromHex(clientIDStr)
	if err != nil {
		helpers.RespondWithError(w, http.StatusBadRequest, "Invalid 'clientID' format. Must be a 24-character hex string (ObjectID)")
		return
	}

	contracts, err := h.service.GetAllContractsByClientID(r.Context(), clientID)
	if err != nil {
		log.Printf("Error listing contracts: %v", err)
		helpers.RespondWithError(w, http.StatusInternalServerError, "Failed to fetch contracts")
		return
	}

	helpers.RespondWithJSON(w, http.StatusOK, contracts)
}
