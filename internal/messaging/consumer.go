package messaging

import (
	"encoding/json"
	"log"

	"github.com/stock-ahora/api-notification/internal/domain"
	"github.com/stock-ahora/api-notification/internal/service"

	"github.com/streadway/amqp"
)

type Consumer struct {
	channel *amqp.Channel
	service *service.NotificationService
	done    chan bool
}

func NewConsumer(channel *amqp.Channel, service *service.NotificationService) *Consumer {
	return &Consumer{
		channel: channel,
		service: service,
		done:    make(chan bool),
	}
}

func (c *Consumer) Start() error {
	if err := c.setupQueues(); err != nil {
		return err
	}

	msgs, err := c.channel.Consume(
		"notifications.movements",
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	go c.processMessages(msgs)

	log.Println("🎧 Consumer started. Waiting for messages...")
	<-c.done

	return nil
}

func (c *Consumer) setupQueues() error {
	// Declarar exchange
	err := c.channel.ExchangeDeclare(
		"movements",
		"topic",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	// Declarar cola
	q, err := c.channel.QueueDeclare(
		"notifications.movements",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	// Bindings
	bindings := []string{"movement.created", "movement.updated"}
	for _, key := range bindings {
		err = c.channel.QueueBind(
			q.Name,
			key,
			"movements",
			false,
			nil,
		)
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *Consumer) processMessages(msgs <-chan amqp.Delivery) {
	for msg := range msgs {
		go c.handleMessage(msg)
	}
}

func (c *Consumer) handleMessage(msg amqp.Delivery) {
	log.Printf("📥 Message received: %s", msg.RoutingKey)

	var event domain.MovementEvent
	if err := json.Unmarshal(msg.Body, &event); err != nil {
		log.Printf("❌ Error deserializing: %v", err)
		msg.Nack(false, false)
		return
	}

	var err error
	switch msg.RoutingKey {
	case "movement.created":
		err = c.service.ProcessMovementCreated(event)
	case "movement.updated":
		err = c.service.ProcessMovementUpdated(event)
	default:
		log.Printf("⚠️ Unhandled routing key: %s", msg.RoutingKey)
	}

	if err != nil {
		log.Printf("❌ Error processing: %v", err)
		msg.Nack(false, true)
		return
	}

	msg.Ack(false)
}

func (c *Consumer) Stop() {
	close(c.done)
}
