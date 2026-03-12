package handlers

import (
	"context"
	"net/http"
	"time"

	"github.com/francisco/distributed-job-platform/internal/handlers/helpers"
	"go.mongodb.org/mongo-driver/mongo"
)

type HealthHandler struct {
	client *mongo.Client
}

func NewHealthHandler(client *mongo.Client) *HealthHandler {
	return &HealthHandler{
		client: client,
	}
}

func (h *HealthHandler) Health(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()

	err := h.client.Ping(ctx, nil)
	if err != nil {
		helpers.RespondWithError(w, http.StatusServiceUnavailable, "Database is unreachable")
		return
	}

	helpers.RespondWithJSON(w, http.StatusOK, map[string]string{"status": "UP"})
}
