package domain

import (
	"time"

	"github.com/google/uuid"
)

type MovementEvent struct {
	ID        uuid.UUID            `json:"id"`
	Products  []ProductPerMovement `json:"products"`
	RequestID uuid.UUID            `json:"request_id"`
}

type ProductPerMovement struct {
	ProductID    uuid.UUID `json:"product_id"`
	Count        int       `json:"count"`
	MovementID   uuid.UUID `json:"movement_id"`
	MovementType int       `json:"movement_type"`
	CreatedAt    time.Time `json:"created_at"`
}
