package main

import (
	"log"
	"resource-service/internal/handler"
	"resource-service/internal/repository"
	"resource-service/internal/repository/db"
	"resource-service/internal/routes"
	"resource-service/internal/service"

	"github.com/gin-gonic/gin"
)

func main() {
	db := db.InitDB()
	router := gin.Default()
	DriverRepo := repository.NewDriverRepository(db)
	AutoRepo := repository.NewAutoRepository(db)

	DriverService := service.NewDriverService(DriverRepo)
	AutoService := service.NewAutoService(AutoRepo)

	DriverController := handler.NewDriverHandler(DriverService)
	AutoController := handler.NewAutoHandler(AutoService)

	routes.AutoRoutes(router, AutoController)
	routes.DriverRoutes(router, DriverController)

	err := router.Run(":8081")
	if err != nil {
		log.Fatal(err)
	}
}
