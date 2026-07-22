package routes

import (
	"user-service/internal/controller"
	"user-service/internal/middleware"

	"github.com/gin-gonic/gin"
)

func UserRoutes(router *gin.Engine, userController *controller.UserContoller) {
	userRouter := router.Group("/user")
	userRouter.Use(middleware.AuthMiddleware())
	{
		userRouter.GET("/:id", userController.GetByID)
	}
}

func AuthRoutes(router *gin.Engine, authController *controller.AuthContoller) {
	authRouter := router.Group("/auth")
	{
		authRouter.POST("/register", authController.SignUp)
		authRouter.POST("/login", authController.Login)
	}
}
