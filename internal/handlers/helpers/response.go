package helpers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"go.mongodb.org/mongo-driver/mongo"
)

// RespondWithError is a helper to format HTTP error responses consistently.
func RespondWithError(w http.ResponseWriter, code int, message string) {
	RespondWithJSON(w, code, map[string]any{
		"error":   http.StatusText(code),
		"message": message,
		"status":  code,
	})
}

// HandleError detects the error type and responds with the appropriate status code.
func HandleError(w http.ResponseWriter, err error) {
	if err == nil {
		return
	}

	code := http.StatusInternalServerError
	message := err.Error()

	// Detect Not Found
	if errors.Is(err, mongo.ErrNoDocuments) || strings.Contains(strings.ToLower(message), "not found") {
		code = http.StatusNotFound
	}

	RespondWithError(w, code, message)
}

// RespondWithJSON is a helper to write JSON responses with the proper Content-Type.
func RespondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, err := json.Marshal(payload)
	if err != nil {
		// Fallback error if JSON marshalling fails
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error": "Internal Server Error", "message": "Failed to marshal JSON response"}`))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}
