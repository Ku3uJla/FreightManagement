package main

import (
	"context"
	"log"
	"os"
	"time"

	"notification-service/events"
	provider "notification-service/providers"
	"notification-service/rabbitmq"
)

func main() {
	rabbitURL := os.Getenv("RABBITMQ_URL")
	// 1. Инициализация клиета RabbitMQ
	client, err := rabbitmq.NewClient(rabbitURL)
	if err != nil {
		log.Fatalf("Ошибка подключения: %v", err)
	}
	defer client.Close()

	mockProvider := provider.NewMockNotificationProvider()

	userPub := rabbitmq.NewUserPublisher(client)
	orderPub := rabbitmq.NewOrderPublisher(client)
	driverPub := rabbitmq.NewDriverPublisher(client)

	consumer := rabbitmq.NewSingleConsumer(client, mockProvider)

	allKeys := []string{
		events.UserCreated,
		events.OrderCreated,
		events.OrderUpdated,
		events.DriverOrder,
	}

	err = consumer.StartListening("notifications_queue", allKeys)
	if err != nil {
		log.Fatalf("Ошибка запуска Consumer: %v", err)
	}

	time.Sleep(300 * time.Millisecond) // Даем временя на обработку биндингов
	ctx := context.Background()

	// -------------------------------------------------------------
	// Эмуляция отправки событий из разных издателей
	// -------------------------------------------------------------

	log.Println("\n--- 📤 UserPublisher отправляет событие ---")
	_ = userPub.PublishUserCreated(ctx, events.UserCreatedPayload{
		UserID:    "usr_101",
		Email:     "john@example.com",
		Name:      "Джон Доу",
		CreatedAt: time.Now(),
	})

	time.Sleep(1 * time.Second)

	log.Println("\n--- 📤 OrderPublisher отправляет событие создания ---")
	_ = orderPub.PublishOrderCreated(ctx, events.OrderPayload{
		OrderID:   "ord_555",
		UserID:    "usr_101",
		UserEmail: "john@example.com",
		Status:    "created",
		Amount:    45.00,
		UpdatedAt: time.Now(),
	})

	time.Sleep(1 * time.Second)

	log.Println("\n--- 📤 DriverPublisher отправляет событие вызова водителя ---")
	_ = driverPub.PublishDriverAssigned(ctx, events.DriverOrderPayload{
		DriverID:    "drv_777",
		DriverEmail: "driver_alex@example.com",
		OrderID:     "ord_555",
		Message:     "Заберите пассажира на ул. Пушкина, д. 10",
		AssignedAt:  time.Now(),
	})

	time.Sleep(1 * time.Second)

	log.Println("\n--- 📤 OrderPublisher отправляет событие обновления ---")
	_ = orderPub.PublishOrderUpdated(ctx, events.OrderPayload{
		OrderID:   "ord_555",
		UserID:    "usr_101",
		UserEmail: "john@example.com",
		DriverID:  "drv_777",
		Status:    "in_progress",
		Amount:    45.00,
		UpdatedAt: time.Now(),
	})

	// Ожидание завершения работы воркера
	time.Sleep(2 * time.Second)
	log.Println("\n✅ Демонстрация успешно завершена.")
}
