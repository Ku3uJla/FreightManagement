package rabbitmq

import (
	"fmt"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

const ExchangeName = "app.events"

type Client struct {
	Conn    *amqp.Connection
	Channel *amqp.Channel
}

func NewClient(amqpURL string) (*Client, error) {
	conn, err := amqp.Dial(amqpURL)
	if err != nil {
		return nil, fmt.Errorf("ошибка подключения к RabbitMQ: %w", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("ошибка открытия канала: %w", err)
	}

	err = ch.ExchangeDeclare(
		ExchangeName, // name
		"topic",      // type
		true,         // durable
		false,        // auto-deleted
		false,        // internal
		false,        // no-wait
		nil,          // args
	)
	if err != nil {
		ch.Close()
		conn.Close()
		return nil, fmt.Errorf("ошибка объявления exchange: %w", err)
	}

	log.Println("✅ Успешное подключение к RabbitMQ и Exchange:", ExchangeName)
	return &Client{Conn: conn, Channel: ch}, nil
}

func (c *Client) Close() {
	if c.Channel != nil {
		c.Channel.Close()
	}
	if c.Conn != nil {
		c.Conn.Close()
	}
}
