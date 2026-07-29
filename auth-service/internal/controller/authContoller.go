package controller

import (
	"auth-service/internal/features"
	"auth-service/internal/repository/model"
	"auth-service/internal/service"

	"github.com/gin-gonic/gin"
)

type AuthController struct {
	authService service.AuthService
}

func NewAuthController(as service.AuthService) *AuthController {
	return &AuthController{authService: as}
}

func (s *AuthController) Login(c *gin.Context) {
	var user model.Auth

	ctx := c.Request.Context()
	err := c.ShouldBindJSON(&user)
	if err != nil {
		c.JSON(401, gin.H{"error": err})
		return
	}
	ID, err := s.authService.LoginByLogin(ctx, user.Login, user.Password)
	if err != nil {
		c.JSON(403, gin.H{"error": err})
		return
	}

	token, err := features.CreateToken(ID)
	if err != nil {
		c.JSON(500, gin.H{"error": err})
		return
	}

	c.SetCookie("JWT", token, 3600, "/", "localhost", false, true)
	c.JSON(200, gin.H{"message": "logged"})
}

func (s *AuthController) SignUp(ctx *gin.Context) {
	var user model.Auth
	ctx.ShouldBindJSON(&user)
	err := s.authService.SignUp(ctx.Request.Context(), user.Login, user.Email, user.Password)
	if err != nil {
		ctx.JSON(403, gin.H{"error": err})
		return
	}
	ctx.JSON(200, gin.H{"message": "Signed UP"})
}
