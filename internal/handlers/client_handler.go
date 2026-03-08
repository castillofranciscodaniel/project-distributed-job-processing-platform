package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/francisco/distributed-job-platform/internal/domain/client"
	"github.com/francisco/distributed-job-platform/internal/handlers/helpers"
	"github.com/go-chi/chi/v5"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CreateClientRequest struct {
	Name       string `json:"name"`
	LastName   string `json:"lastName"`
	Email      string `json:"email"`
	Identifier string `json:"identifier"`
}

type ClientHandler struct {
	service *client.ClientService
}

func NewClientHandler(service *client.ClientService) *ClientHandler {
	return &ClientHandler{
		service: service,
	}
}

func (h *ClientHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req CreateClientRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		helpers.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if req.Name == "" || req.LastName == "" || req.Email == "" || req.Identifier == "" {
		helpers.RespondWithError(w, http.StatusBadRequest, "Name, lastName, email and identifier are required")
		return
	}

	c, err := h.service.CreateClient(r.Context(), req.Name, req.LastName, req.Email, req.Identifier)
	if err != nil {
		helpers.HandleError(w, err)
		return
	}

	helpers.RespondWithJSON(w, http.StatusCreated, c)
}

func (h *ClientHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		helpers.RespondWithError(w, http.StatusBadRequest, "Invalid ID format")
		return
	}

	c, err := h.service.GetClientByID(r.Context(), id)
	if err != nil {
		helpers.HandleError(w, err)
		return
	}

	helpers.RespondWithJSON(w, http.StatusOK, c)
}
