package controller

import (
	"context"
	"log"
	"strconv"
	"user-service/internal/service"

	pb "user-service/proto/userpb"

	"github.com/gin-gonic/gin"
)

type UserController struct {
	UserService *service.UserService
	pb.UnimplementedUserServiceServer
}

func NewUserController(userService *service.UserService) *UserController {
	return &UserController{UserService: userService}
}

func (s *UserController) GetByID(c *gin.Context) {
	ctx := c.Request.Context()
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(500, gin.H{"error": err})
		return
	}
	user, err := s.UserService.GetUser(ctx, id)
	if err != nil {
		c.JSON(500, gin.H{"error": err})
		return
	}
	c.JSON(200, gin.H{"message": user})
}

func (h *UserController) UpdateRole(c *gin.Context) {
	ctx := c.Request.Context()
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(500, gin.H{"error": err})
		return
	}

	var role int
	err = c.ShouldBindJSON(&role)
	if err != nil {
		c.JSON(403, gin.H{"error": err})
		return
	}
	err = h.UserService.UpdateRole(ctx, id, role)
	if err != nil {
		c.JSON(500, gin.H{"error": err})
	}
}

func (h *UserController) CreateUser(
	ctx context.Context,
	req *pb.CreateUserRequest,
) (*pb.CreateUserResponse, error) {

	userID := int(req.Id)
	log.Println("CreateUser на serverе вызван")
	err := h.UserService.CreateUser(ctx, userID)
	if err != nil {
		return &pb.CreateUserResponse{
			Success: false,
		}, err
	}
	return &pb.CreateUserResponse{
		Success: true,
	}, nil
}

func (h *UserController) GetUser(
	ctx context.Context,
	req *pb.GetUserRequest,
) (*pb.GetUserResponse, error) {
	var fullname string
	var role int64
	userID := int(req.Id)
	model, err := h.UserService.GetUser(ctx, userID)

	if model.FullName != nil {
		fullname = *model.FullName
	}
	if model.Role != nil {
		role = int64(*model.Role)
	}
	if err != nil {
		return &pb.GetUserResponse{}, err
	}
	return &pb.GetUserResponse{
		Role:     role,
		Phone:    model.Phone,
		FullName: fullname,
	}, nil
}
