package routes

import (
	"auth-service/internal/controller"

	"github.com/gin-gonic/gin"
)

func AuthRoutes(router *gin.Engine, authController *controller.AuthController) {
	authRouter := router.Group("/auth")
	{
		authRouter.POST("/register", authController.SignUp)
		authRouter.POST("/login", authController.Login)
	}
}
