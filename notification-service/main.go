package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"notification-service/events"
	provider "notification-service/providers"
	"notification-service/rabbitmq"
)

func main() {
	rabbitURL := os.Getenv("RABBITMQ_URL")

	log.Printf("Connecting to RabbitMQ at: %s", rabbitURL)

	// 2. Инициализация клиента RabbitMQ
	client, err := rabbitmq.NewClient(rabbitURL)
	if err != nil {
		log.Fatalf("Ошибка подключения к RabbitMQ: %v", err)
	}
	defer client.Close()

	// 3. Инициализация провайдера уведомлений
	mockProvider := provider.NewMockNotificationProvider()

	// 4. Инициализация и запуск Consumer
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

	log.Println("🚀 Notification Service успешно запущен и ожидает сообщений...")

	// 5. БЛОКИРОВКА main: Graceful Shutdown
	// Ждем сигнал от ОС (Ctrl+C локально или SIGTERM от Docker/Kubernetes)
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	<-stop // Поток заблокирован здесь, пока сервис работает

	log.Println("🛑 Завершение работы Notification Service...")
}
