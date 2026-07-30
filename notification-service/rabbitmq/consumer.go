package rabbitmq

import (
	"encoding/json"
	"fmt"
	"log"

	"notification-service/events"
	provider "notification-service/providers"
)

type SingleConsumer struct {
	client   *Client
	provider provider.Provider
}

func NewSingleConsumer(client *Client, p provider.Provider) *SingleConsumer {
	return &SingleConsumer{
		client:   client,
		provider: p,
	}
}

// StartListening подписывает один Consumer на все переданные routingKeys
func (c *SingleConsumer) StartListening(queueName string, routingKeys []string) error {
	q, err := c.client.Channel.QueueDeclare(
		queueName,
		true,  // durable
		false, // auto-delete
		false, // exclusive
		false, // no-wait
		nil,
	)
	if err != nil {
		return err
	}

	// Привязываем все необходимые ключи маршрутизации к одной очереди
	for _, key := range routingKeys {
		err = c.client.Channel.QueueBind(
			q.Name,
			key,
			ExchangeName,
			false,
			nil,
		)
		if err != nil {
			return err
		}
	}

	msgs, err := c.client.Channel.Consume(
		q.Name,
		"",    // consumer tag
		true,  // auto-ack
		false, // exclusive
		false, // no-local
		false, // no-wait
		nil,
	)
	if err != nil {
		return err
	}

	log.Printf("📥 [SingleConsumer] Единый обработчик запущен на очереди [%s]", q.Name)

	go func() {
		for msg := range msgs {
			c.processEvent(msg.RoutingKey, msg.Body)
		}
	}()

	return nil
}

// processEvent разбирает события и использует provider.Send(...)
func (c *SingleConsumer) processEvent(routingKey string, body []byte) {
	switch routingKey {

	case events.UserCreated:
		var payload events.UserCreatedPayload
		if err := json.Unmarshal(body, &payload); err == nil {
			msg := fmt.Sprintf("Добро пожаловать, %s!", payload.Name)
			_ = c.provider.Send(payload.Email, msg)
		}

	case events.OrderCreated:
		var payload events.OrderPayload
		if err := json.Unmarshal(body, &payload); err == nil {
			msg := fmt.Sprintf("Ваш заказ #%s успешно создан на сумму %.2f$", payload.OrderID, payload.Amount)
			_ = c.provider.Send(payload.UserEmail, msg)
		}

	case events.OrderUpdated:
		var payload events.OrderPayload
		if err := json.Unmarshal(body, &payload); err == nil {
			msg := fmt.Sprintf("Статус вашего заказа #%s изменен на: %s", payload.OrderID, payload.Status)
			_ = c.provider.Send(payload.UserEmail, msg)
		}

	case events.DriverOrder:
		var payload events.DriverOrderPayload
		if err := json.Unmarshal(body, &payload); err == nil {
			msg := fmt.Sprintf("Вы привязаны к заказу #%s! %s", payload.OrderID, payload.Message)
			_ = c.provider.Send(payload.DriverEmail, msg)
		}

	default:
		log.Printf("⚠️ [SingleConsumer] Неизвестный routing key: %s", routingKey)
	}
}
