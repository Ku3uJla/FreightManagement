package main

import (
	"log"
	"user-service/internal/controller"
	"user-service/internal/repository"
	"user-service/internal/repository/db"
	"user-service/internal/routes"
	"user-service/internal/service"

	"github.com/gin-gonic/gin"
)

func main() {
	db := db.InitDB()
	router := gin.Default()

	userStore := repository.NewUserRepository(db)
	userService := service.NewUserService(userStore)
	userController := controller.NewUserController(userService)

	authService := service.NewAuthService(userStore)
	authController := controller.NewAuthController(authService)
	routes.UserRoutes(router, userController)
	routes.AuthRoutes(router, authController)

	router.GET("/", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{"message": "user-service"})
	})

	if err := router.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}
