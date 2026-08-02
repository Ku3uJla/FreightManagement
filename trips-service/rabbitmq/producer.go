package rabbitmq

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"trips-service/dto"

	"github.com/google/uuid"
	amqp "github.com/rabbitmq/amqp091-go"
)

const (
	TripsExchange = "trips.events" // Единый Exchange для всех событий рейсов
)

type TripEventProducer interface {
	PublishStatusChanged(ctx context.Context, payload dto.TripStatusUpdatedPayload) error
	PublishTripCompleted(ctx context.Context, payload dto.TripCompletedPayload) error
	PublishTripCancelled(ctx context.Context, payload dto.TripCancelledPayload) error
}

type producer struct {
	channel *amqp.Channel
}

// NewProducer создает экземпляр Producer
func NewProducer(ch *amqp.Channel) TripEventProducer {
	return &producer{channel: ch}
}

// PublishStatusChanged отправляет событие о смене промежуточного статуса
func (p *producer) PublishStatusChanged(ctx context.Context, payload dto.TripStatusUpdatedPayload) error {
	event := dto.EventEnvelope{
		EventID:   uuid.NewString(),
		EventType: dto.EventTypeTripStatusUpdated,
		Timestamp: time.Now().UTC(),
		Payload:   payload,
	}

	routingKey := "trip.status_updated"
	return p.publish(ctx, TripsExchange, routingKey, event)
}

// PublishTripCompleted отправляет событие об успешном завершении рейса
func (p *producer) PublishTripCompleted(ctx context.Context, payload dto.TripCompletedPayload) error {
	event := dto.EventEnvelope{
		EventID:   uuid.NewString(),
		EventType: dto.EventTypeTripCompleted,
		Timestamp: time.Now().UTC(),
		Payload:   payload,
	}

	routingKey := "trip.completed"
	return p.publish(ctx, TripsExchange, routingKey, event)
}

// PublishTripCancelled отправляет событие об отмене рейса
func (p *producer) PublishTripCancelled(ctx context.Context, payload dto.TripCancelledPayload) error {
	event := dto.EventEnvelope{
		EventID:   uuid.NewString(),
		EventType: dto.EventTypeTripCancelled,
		Timestamp: time.Now().UTC(),
		Payload:   payload,
	}

	routingKey := "trip.cancelled"
	return p.publish(ctx, TripsExchange, routingKey, event)
}

// Внутренний приватный метод отправки
func (p *producer) publish(ctx context.Context, exchange, routingKey string, event dto.EventEnvelope) error {
	// 1. Декларируем Exchange
	err := p.channel.ExchangeDeclare(
		exchange,
		"topic",
		true,  // durable
		false, // auto-deleted
		false, // internal
		false, // no-wait
		nil,
	)
	if err != nil {
		return fmt.Errorf("failed to declare exchange: %w", err)
	}

	// 2. Сериализуем EventEnvelope в JSON
	body, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to marshal event: %w", err)
	}

	// 3. Устанавливаем таймаут отправки
	pubCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	// 4. Публикация
	err = p.channel.PublishWithContext(
		pubCtx,
		exchange,
		routingKey,
		false,
		false,
		amqp.Publishing{
			ContentType:  "application/json",
			DeliveryMode: amqp.Persistent,
			Timestamp:    event.Timestamp,
			MessageId:    event.EventID,
			Type:         event.EventType,
			Body:         body,
		},
	)
	if err != nil {
		return fmt.Errorf("failed to publish message [%s]: %w", routingKey, err)
	}

	log.Printf("📤 [Producer] Published [%s] | EventID: %s", routingKey, event.EventID)
	return nil
}
