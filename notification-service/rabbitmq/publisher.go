package rabbitmq

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"notification-service/events"

	amqp "github.com/rabbitmq/amqp091-go"
)

// BasePublisher содержит общую логику отправки
type BasePublisher struct {
	client *Client
}

func (p *BasePublisher) publish(ctx context.Context, routingKey string, payload interface{}) error {
	data, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("ошибка сериализации JSON: %w", err)
	}

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	return p.client.Channel.PublishWithContext(
		ctx,
		ExchangeName,
		routingKey,
		false,
		false,
		amqp.Publishing{
			ContentType:  "application/json",
			DeliveryMode: amqp.Persistent,
			Timestamp:    time.Now(),
			Body:         data,
		},
	)
}

// 1. UserPublisher — отправка пользовательских событий
type UserPublisher struct {
	BasePublisher
}

func NewUserPublisher(client *Client) *UserPublisher {
	return &UserPublisher{BasePublisher{client: client}}
}

func (p *UserPublisher) PublishUserCreated(ctx context.Context, payload events.UserCreatedPayload) error {
	return p.publish(ctx, events.UserCreated, payload)
}

// 2. OrderPublisher — отправка событий заказов
type OrderPublisher struct {
	BasePublisher
}

func NewOrderPublisher(client *Client) *OrderPublisher {
	return &OrderPublisher{BasePublisher{client: client}}
}

func (p *OrderPublisher) PublishOrderCreated(ctx context.Context, payload events.OrderPayload) error {
	return p.publish(ctx, events.OrderCreated, payload)
}

func (p *OrderPublisher) PublishOrderUpdated(ctx context.Context, payload events.OrderPayload) error {
	return p.publish(ctx, events.OrderUpdated, payload)
}

// 3. DriverPublisher — отправка событий водителей
type DriverPublisher struct {
	BasePublisher
}

func NewDriverPublisher(client *Client) *DriverPublisher {
	return &DriverPublisher{BasePublisher{client: client}}
}

func (p *DriverPublisher) PublishDriverAssigned(ctx context.Context, payload events.DriverOrderPayload) error {
	return p.publish(ctx, events.DriverOrder, payload)
}
