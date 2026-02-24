package handler

import (
	"encoding/json"
	"net/http"
)

// HealthCheck godoc
// @Summary      Cek Status Server
// @Description  Endpoint untuk memastikan backend berjalan lancar
// @Tags         System
// @Produce      json
// @Success      200  {object}  map[string]string
// @Router       /health [get]
func HealthCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"status": "server is running",
	})
}