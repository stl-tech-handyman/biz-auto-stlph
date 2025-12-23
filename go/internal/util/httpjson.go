package util

import (
	"encoding/json"
	"net/http"
)

// WriteJSON writes a JSON response
func WriteJSON(w http.ResponseWriter, status int, data any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(data)
}

// WriteError writes an error JSON response
func WriteError(w http.ResponseWriter, status int, message string) error {
	return WriteJSON(w, status, map[string]string{
		"error": message,
	})
}

// ReadJSON reads JSON from request body
func ReadJSON(r *http.Request, dst any) error {
	return json.NewDecoder(r.Body).Decode(dst)
}

