package publisher

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"order-service/events"

	amqp "github.com/rabbitmq/amqp091-go"
)

type OrderPublisher interface {
	PublishOrderCreated(ctx context.Context, payload events.OrderPayload) error
	PublishOrderUpdated(ctx context.Context, payload events.OrderPayload) error
	PublishOrderCanceled(ctx context.Context, payload events.OrderPayload) error
	Close()
}

type orderPublisher struct {
	conn *amqp.Connection
	ch   *amqp.Channel
}

// NewOrderPublisher создает подключение к RabbitMQ и объявляет Exchange
func NewOrderPublisher(amqpURL string) (OrderPublisher, error) {
	conn, err := amqp.Dial(amqpURL)
	if err != nil {
		return nil, fmt.Errorf("ошибка подключения к RabbitMQ: %w", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		_ = conn.Close()
		return nil, fmt.Errorf("ошибка открытия канала RabbitMQ: %w", err)
	}

	// Гарантируем наличие Topic Exchange
	err = ch.ExchangeDeclare(
		events.ExchangeName, // name
		"topic",             // type
		true,                // durable (сохраняется при перезапуске брокера)
		false,               // auto-deleted
		false,               // internal
		false,               // no-wait
		nil,                 // args
	)
	if err != nil {
		_ = ch.Close()
		_ = conn.Close()
		return nil, fmt.Errorf("ошибка объявления exchange [%s]: %w", events.ExchangeName, err)
	}

	log.Println("✅ [Order Service] Publisher успешно подключен к RabbitMQ")
	return &orderPublisher{conn: conn, ch: ch}, nil
}

// publish — универсальный внутренний метод отправки
func (p *orderPublisher) publish(ctx context.Context, routingKey string, payload events.OrderPayload) error {
	// Добавляем текущее время, если оно не было проставлено
	if payload.Timestamp.IsZero() {
		payload.Timestamp = time.Now()
	}

	data, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("ошибка маршалинга JSON: %w", err)
	}

	// Ограничиваем таймаут отправки по сети 5 секундами
	pubCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	err = p.ch.PublishWithContext(
		pubCtx,
		events.ExchangeName,
		routingKey,
		false, // mandatory
		false, // immediate
		amqp.Publishing{
			ContentType:  "application/json",
			DeliveryMode: amqp.Persistent, // Сохраняем сообщение на диск на случай сбоя RabbitMQ
			Timestamp:    time.Now(),
			Body:         data,
		},
	)
	if err != nil {
		return fmt.Errorf("ошибка публикации сообщения [%s]: %w", routingKey, err)
	}

	log.Printf("📤 [Order Service] Опубликовано событие [%s] для OrderID=%s", routingKey, payload.OrderID)
	return nil
}

func (p *orderPublisher) PublishOrderCreated(ctx context.Context, payload events.OrderPayload) error {
	return p.publish(ctx, events.OrderCreated, payload)
}

func (p *orderPublisher) PublishOrderUpdated(ctx context.Context, payload events.OrderPayload) error {
	return p.publish(ctx, events.OrderUpdated, payload)
}

func (p *orderPublisher) PublishOrderCanceled(ctx context.Context, payload events.OrderPayload) error {
	return p.publish(ctx, events.OrderCanceled, payload)
}

func (p *orderPublisher) Close() {
	if p.ch != nil {
		_ = p.ch.Close()
	}
	if p.conn != nil {
		_ = p.conn.Close()
	}
	log.Println("🔌 [Order Service] Подключение к RabbitMQ закрыто")
}
