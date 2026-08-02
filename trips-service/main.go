package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	amqp "github.com/rabbitmq/amqp091-go"

	"trips-service/internal/handler"
	"trips-service/internal/repository"
	"trips-service/internal/repository/db"
	"trips-service/internal/routes"
	"trips-service/internal/service"
	"trips-service/rabbitmq"
)

func main() {
	// 1. Инициализация RabbitMQ
	rabbitURL := os.Getenv("RABBITMQ_URL")
	if rabbitURL == "" {
		rabbitURL = "amqp://guest:guest@localhost:5672/"
	}

	conn, err := amqp.Dial(rabbitURL)
	if err != nil {
		log.Fatalf("Критическая ошибка подключения к RabbitMQ: %v", err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("Критическая ошибка создания RabbitMQ Channel: %v", err)
	}
	defer ch.Close()

	// Инициализируем продюсер и сохраняем в переменную
	pub := rabbitmq.NewProducer(ch)

	// 2. Инициализация базы данных
	database := db.InitDB()

	// 3. Сборка слоев (Repository -> Service -> Handler)
	tripStore := repository.NewTripsRepository(database)
	tripService := service.NewTripsService(tripStore, pub)
	tripController := handler.NewTripHandler(tripService)

	// 5. Запуск HTTP сервера (Gin)
	router := gin.Default()

	routes.TripRoutes(
		router,
		tripController,
	)

	log.Println("HTTP сервер запущен на :8080")
	if err := router.Run(":8080"); err != nil {
		log.Fatalf("Ошибка запуска HTTP сервера: %v", err)
	}
}
