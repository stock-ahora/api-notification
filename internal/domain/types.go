package domain

const (
	MovementTypeEntry = 1 // Entrada
	MovementTypeExit  = 2 // Salida
)

const (
	NotificationTypeMovementCreated = "movement.created"
	NotificationTypeMovementUpdated = "movement.updated"
)

func GetMovementTypeString(movementType int) string {
	if movementType == MovementTypeEntry {
		return "entrada"
	}
	return "salida"
}
