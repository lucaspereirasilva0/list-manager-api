package handlers

import (
	"encoding/json"
	"net/http"
)

// VersionResponse represents the structure of the version JSON response.
type VersionResponse struct {
	Version string `json:"version"`
}

// GetVersion handles requests for the application version.
// It returns a JSON object with the current application version.
func GetVersion(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// TODO: Get the actual version dynamically (e.g., from a build flag, env var, or config)
	// For now, a hardcoded version is used.
	version := "1.0.0"

	response := VersionResponse{Version: version}
	err := json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
