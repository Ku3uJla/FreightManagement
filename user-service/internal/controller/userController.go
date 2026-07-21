package controller

import (
	"user-service/internal/repository/model"
	"user-service/internal/service"

	"github.com/gin-gonic/gin"
)

type UserContoller struct {
	userService *service.UserService
}

func NewUserController(userService *service.UserService) *UserContoller {
	return &UserContoller{userService: userService}
}

func (s *UserContoller) Login(ctx *gin.Context) {
	var user model.User

	ctx.ShouldBindJSON(&user)
	token, err := s.userService.Login(ctx, user.Email, user.Password)
	if err != nil {
		ctx.JSON(403, gin.H{"error": err})
		return
	}
	ctx.SetCookie("JWT", token, 3600, "/", "localhost", false, true)
	ctx.JSON(200, gin.H{"message": "logged"})
}

func (s *UserContoller) SignUp(ctx *gin.Context) {
	var user model.User
	ctx.ShouldBindJSON(&user)
	if err := s.userService.SignUp(ctx, user.Email, user.Password); err != nil {
		ctx.JSON(403, gin.H{"error": err})
	}
	ctx.JSON(200, gin.H{"message": "Signed UP"})
}

func (s *UserContoller) GetByID(ctx *gin.Context) {
	id := ctx.Param("id")
	user, err := s.userService.GetUser(ctx, id)
	if err != nil {
		ctx.JSON(500, gin.H{"error": err})
	}
	ctx.JSON(200, gin.H{"message": user})
}
