package controller

import (
	"user-service/internal/service"

	"github.com/gin-gonic/gin"
)

type UserContoller struct {
	UserService *service.UserService
}

func NewUserController(userService *service.UserService) *UserContoller {
	return &UserContoller{UserService: userService}
}

func (s *UserContoller) GetByID(ctx *gin.Context) {
	id := ctx.Param("id")
	user, err := s.UserService.GetUser(ctx.Request.Context(), id)
	if err != nil {
		ctx.JSON(500, gin.H{"error": err})
	}
	ctx.JSON(200, gin.H{"message": user})
}
