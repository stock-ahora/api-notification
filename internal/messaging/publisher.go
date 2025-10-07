package messaging

import (
	"encoding/json"
	"time"

	"github.com/stock-ahora/api-notification/internal/domain"

	"github.com/streadway/amqp"
)

type Publisher struct {
	channel *amqp.Channel
}

func NewPublisher(channel *amqp.Channel) *Publisher {
	// Declarar exchange de notificaciones
	channel.ExchangeDeclare(
		"notifications",
		"topic",
		true,
		false,
		false,
		false,
		nil,
	)

	return &Publisher{
		channel: channel,
	}
}

func (p *Publisher) PublishNotification(notification domain.Notification) error {
	body, err := json.Marshal(notification)
	if err != nil {
		return err
	}

	return p.channel.Publish(
		"notifications",
		"notification.created",
		false,
		false,
		amqp.Publishing{
			ContentType:  "application/json",
			Body:         body,
			DeliveryMode: amqp.Persistent,
			Timestamp:    time.Now(),
		},
	)
}
