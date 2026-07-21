package routes

import (
	"user-service/internal/controller"

	"github.com/gin-gonic/gin"
)

func UserRoutes(router *gin.Engine, userContoller *controller.UserContoller) {
	userRouter := router.Group("/user")
	{
		userRouter.GET("/:id", userContoller.GetByID)
		userRouter.POST("/register", userContoller.SignUp)
		userRouter.POST("/login", userContoller.Login)
	}
}
