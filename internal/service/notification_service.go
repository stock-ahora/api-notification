package service

import (
	"fmt"
	"log"
	"time"

	"github.com/stock-ahora/api-notification/internal/domain"
	"github.com/stock-ahora/api-notification/internal/messaging"

	"github.com/google/uuid"
)

type NotificationService struct {
	publisher *messaging.Publisher
	templates *TemplateManager
}

func NewNotificationService(publisher *messaging.Publisher) *NotificationService {
	return &NotificationService{
		publisher: publisher,
		templates: NewTemplateManager(),
	}
}

func (s *NotificationService) ProcessMovementCreated(event domain.MovementEvent) error {
	log.Printf("📨 Processing movement created: %s", event.ID)

	for _, product := range event.Products {
		movementType := domain.GetMovementTypeString(product.MovementType)
		templateKey := fmt.Sprintf("movement.created.%s", movementType)

		template, exists := s.templates.Get(templateKey)
		if !exists {
			log.Printf("⚠️ Template not found: %s", templateKey)
			continue
		}

		notification := s.createNotification(
			domain.NotificationTypeMovementCreated,
			template,
			product,
			event.RequestID.String(),
		)

		if err := s.publisher.PublishNotification(notification); err != nil {
			log.Printf("❌ Error publishing notification: %v", err)
			return err
		}

		log.Printf("✅ Notification published: %s", notification.ID)
	}

	return nil
}

func (s *NotificationService) ProcessMovementUpdated(event domain.MovementEvent) error {
	log.Printf("📨 Processing movement updated: %s", event.ID)

	for _, product := range event.Products {
		movementType := domain.GetMovementTypeString(product.MovementType)
		templateKey := fmt.Sprintf("movement.updated.%s", movementType)

		template, exists := s.templates.Get(templateKey)
		if !exists {
			continue
		}

		notification := s.createNotification(
			domain.NotificationTypeMovementUpdated,
			template,
			product,
			event.RequestID.String(),
		)

		if err := s.publisher.PublishNotification(notification); err != nil {
			return err
		}
	}

	return nil
}

func (s *NotificationService) createNotification(
	notificationType string,
	template MessageTemplate,
	product domain.ProductPerMovement,
	requestID string,
) domain.Notification {
	return domain.Notification{
		ID:        uuid.New().String(),
		Type:      notificationType,
		Title:     template.Title,
		Message:   fmt.Sprintf(template.Message, product.Count, product.ProductID, product.Count),
		Timestamp: time.Now(),
		RequestID: requestID,
		Metadata: map[string]interface{}{
			"product_id":    product.ProductID.String(),
			"movement_id":   product.MovementID.String(),
			"movement_type": domain.GetMovementTypeString(product.MovementType),
			"count":         product.Count,
		},
	}
}

func (s *NotificationService) GetTemplates() map[string]MessageTemplate {
	return s.templates.GetAll()
}
