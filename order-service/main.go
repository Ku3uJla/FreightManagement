package main

import (
	"log"
	"order-service/internal/handler"
	"order-service/internal/repository"
	"order-service/internal/repository/db"
	"order-service/internal/routes"
	"order-service/internal/service"

	"github.com/gin-gonic/gin"
)

func main() {
	db := db.InitDB()
	router := gin.Default()

	orderRepo := repository.NewOrderRepository(db)
	orderService := service.NewOrderService(orderRepo)
	orderHandler := handler.NewOrderHandler(orderService)

	routes.OrderRoutes(router, orderHandler)

	if err := router.Run(":8082"); err != nil {
		log.Fatal(err)
	}
}
