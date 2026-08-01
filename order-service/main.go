package main

import (
	"log"
	"order-service/internal/handler"
	"order-service/internal/repository"
	"order-service/internal/repository/db"
	"order-service/internal/routes"
	"order-service/internal/service"
	"order-service/publisher"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {
	db := db.InitDB()
	router := gin.Default()
	rabbitURL := os.Getenv("RABBITMQ_URL")
	orderPub, err := publisher.NewOrderPublisher(rabbitURL)
	if err != nil {
		log.Fatalf("RabbitMQ connection error: %v", err)
	}
	defer orderPub.Close()

	orderRepo := repository.NewOrderRepository(db)

	orderSrv := service.NewOrderService(orderRepo, orderPub)
	orderHandler := handler.NewOrderHandler(orderSrv)

	routes.OrderRoutes(router, orderHandler)

	if err := router.Run(":8082"); err != nil {
		log.Fatal(err)
	}
}
