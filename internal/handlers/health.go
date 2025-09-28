// internal/handlers/health.go
package handlers

import (
	"encoding/json"
	"net/http"
)

type HealthResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

func HealthCheck(w http.ResponseWriter, r *http.Request) {
	// Solo permitir GET
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	// Configurar header
	w.Header().Set("Content-Type", "application/json")

	// Crear respuesta
	response := HealthResponse{
		Status:  "ok",
		Message: "API funcionando correctamente",
	}

	// Enviar respuesta
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
