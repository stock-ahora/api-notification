package messaging

import (
	"fmt"

	"github.com/streadway/amqp"
)

type RabbitConfig struct {
	URL string
}

func Connect(config RabbitConfig) (*amqp.Connection, *amqp.Channel, error) {
	conn, err := amqp.Dial(config.URL)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to connect: %w", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, nil, fmt.Errorf("failed to open channel: %w", err)
	}

	return conn, ch, nil
}
