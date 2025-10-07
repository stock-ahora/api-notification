package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/stock-ahora/api-notification/internal/domain"
	"github.com/stock-ahora/api-notification/internal/service"

	"github.com/google/uuid"
)

type Handlers struct {
	service *service.NotificationService
}

func NewHandlers(service *service.NotificationService) *Handlers {
	return &Handlers{
		service: service,
	}
}

func (h *Handlers) GetTemplates(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(h.service.GetTemplates())
}

func (h *Handlers) CreateTestNotification(w http.ResponseWriter, r *http.Request) {
	var input struct {
		ProductName string `json:"product_name"`
		Count       int    `json:"count"`
		Type        string `json:"type"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	notification := domain.Notification{
		ID:        uuid.New().String(),
		Type:      "test",
		Title:     "🧪 Notificación de prueba",
		Message:   fmt.Sprintf("Prueba: %d unidades de %s (%s)", input.Count, input.ProductName, input.Type),
		Timestamp: time.Now(),
		RequestID: uuid.New().String(),
		Metadata: map[string]interface{}{
			"product_name": input.ProductName,
			"count":        input.Count,
			"type":         input.Type,
			"test":         true,
		},
	}

	// Aquí podrías publicar la notificación si lo deseas

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(notification)
}
