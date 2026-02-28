package helpers

import (
	"encoding/json"
	"net/http"
)

// RespondWithError is a helper to format HTTP error responses consistently.
func RespondWithError(w http.ResponseWriter, code int, message string) {
	RespondWithJSON(w, code, map[string]string{
		"error":   http.StatusText(code),
		"message": message,
	})
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
