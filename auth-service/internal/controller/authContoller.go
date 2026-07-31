package controller

import (
	"auth-service/dto"
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

func (h *AuthController) SignUp(c *gin.Context) {
	var user dto.CreateAuthRequest
	c.ShouldBindJSON(&user)
	err := h.authService.SignUp(c.Request.Context(), user)
	if err != nil {
		c.JSON(403, gin.H{"error": err})
		return
	}
	c.JSON(200, gin.H{"message": "Signed UP"})
}
