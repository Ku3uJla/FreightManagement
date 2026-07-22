package controller

import (
	"user-service/internal/features"
	"user-service/internal/repository/model"
	"user-service/internal/service"

	"github.com/gin-gonic/gin"
)

type AuthContoller struct {
	AuthService *service.AuthService
}

func NewAuthController(authService *service.AuthService) *AuthContoller {
	return &AuthContoller{AuthService: authService}
}

func (s *AuthContoller) Login(ctx *gin.Context) {
	var user model.User

	ctx.ShouldBindJSON(&user)
	ID, role, err := s.AuthService.Login(ctx.Request.Context(), user.Email, user.Password)
	if err != nil {
		ctx.JSON(403, gin.H{"error": err})
		return
	}

	token, err := features.CreateToken(ID, role)
	if err != nil {
		ctx.JSON(500, gin.H{"error": err})
		return
	}

	ctx.SetCookie("JWT", token, 3600, "/", "localhost", false, true)
	ctx.JSON(200, gin.H{"message": "logged"})
}

func (s *AuthContoller) SignUp(ctx *gin.Context) {
	var user model.User
	ctx.ShouldBindJSON(&user)
	if err := s.AuthService.SignUp(ctx.Request.Context(), user.Login, user.Email, user.Password); err != nil {
		ctx.JSON(403, gin.H{"error": err})
	}
	ctx.JSON(200, gin.H{"message": "Signed UP"})
}
