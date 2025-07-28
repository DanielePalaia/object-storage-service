package api

import (
	"encoding/json"
	"net/http"
	"time"
)

// HealthResponse represents the health check response
type HealthResponse struct {
	Status    string    `json:"status"`
	Timestamp time.Time `json:"timestamp"`
	Service   string    `json:"service"`
	Version   string    `json:"version,omitempty"`
}

// HealthHandler handles health check requests
// @Summary Health check endpoint
// @Description Returns the health status of the service
// @Tags health
// @Produce json
// @Success 200 {object} HealthResponse
// @Router /health [get]
func HealthHandler(w http.ResponseWriter, r *http.Request) {
	response := HealthResponse{
		Status:    "healthy",
		Timestamp: time.Now().UTC(),
		Service:   "object-storage-service",
		Version:   "1.0.0", // You can make this configurable
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to encode health response", http.StatusInternalServerError)
		return
	}
}
