package publisher

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"user-service/events"

	amqp "github.com/rabbitmq/amqp091-go"
)

type UserPublisher struct {
	conn *amqp.Connection
	ch   *amqp.Channel
}

// NewUserPublisher инициализирует подключение и объявляет Exchange
func NewUserPublisher(amqpURL string) (*UserPublisher, error) {
	conn, err := amqp.Dial(amqpURL)
	if err != nil {
		return nil, fmt.Errorf("ошибка подключения к RabbitMQ: %w", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("ошибка открытия канала: %w", err)
	}

	// Объявляем Topic Exchange
	err = ch.ExchangeDeclare(
		events.ExchangeName, // name
		"topic",             // type
		true,                // durable
		false,               // auto-deleted
		false,               // internal
		false,               // no-wait
		nil,                 // args
	)
	if err != nil {
		ch.Close()
		conn.Close()
		return nil, fmt.Errorf("ошибка объявления exchange: %w", err)
	}

	log.Println(" [User Service] Publisher успешно подключен к RabbitMQ")
	return &UserPublisher{conn: conn, ch: ch}, nil
}

// PublishUserCreated отправляет событие регистрации нового пользователя
func (p *UserPublisher) PublishUserCreated(ctx context.Context, payload events.UserCreatedPayload) error {
	data, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("ошибка сериализации JSON: %w", err)
	}

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	err = p.ch.PublishWithContext(
		ctx,
		events.ExchangeName,
		events.UserCreated, // routing key: "user.created"
		false,
		false,
		amqp.Publishing{
			ContentType:  "application/json",
			DeliveryMode: amqp.Persistent, // Сохраняем на диск
			Timestamp:    time.Now(),
			Body:         data,
		},
	)
	if err != nil {
		return fmt.Errorf("ошибка публикации события [user.created]: %w", err)
	}

	log.Printf("📤 [User Service] Опубликовано событие [user.created] для ID=%s", payload.UserID)
	return nil
}

// Close корректно закрывает соединения
func (p *UserPublisher) Close() {
	if p.ch != nil {
		_ = p.ch.Close()
	}
	if p.conn != nil {
		_ = p.conn.Close()
	}
	log.Println("[User Service] Подключение к RabbitMQ закрыто")
}
