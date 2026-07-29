package routes

import (
	"user-service/internal/controller"

	"github.com/gin-gonic/gin"
)

func UserRoutes(
	router *gin.Engine,
	userController *controller.UserController,
) {

	userRouter := router.Group("/user")
	{
		userRouter.GET("/:id", userController.GetByID)
		userRouter.PUT("/:id/role", userController.UpdateRole)
	}
}
