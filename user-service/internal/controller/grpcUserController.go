package controller

import (
	"context"
	"user-service/dto"
	"user-service/internal/service"

	pb "user-service/proto/userpb"
)

type UserGrpcController struct {
	UserService service.UserService
	pb.UnimplementedUserServiceServer
}

func NewUserGrpcController(srv service.UserService) *UserGrpcController {
	return &UserGrpcController{UserService: srv}
}

func (h *UserGrpcController) CreateUser(
	ctx context.Context,
	req *pb.CreateUserRequest,
) (*pb.CreateUserResponse, error) {

	dtoReq := dto.CreateUserGrpcRequest{
		ID:       int(req.Id), // int64 → int
		Email:    req.Email,
		FullName: req.FullName,
		Phone:    req.Phone,
	}

	err := h.UserService.CreateUser(ctx, dtoReq)
	if err != nil {
		return &pb.CreateUserResponse{Success: false}, err
	}

	return &pb.CreateUserResponse{Success: true}, nil
}

func (h *UserGrpcController) GetUser(
	ctx context.Context,
	req *pb.GetUserRequest,
) (*pb.GetUserResponse, error) {
	var fullname string
	userID := int(req.Id)
	model, err := h.UserService.GetUser(ctx, userID)

	if model.FullName != nil {
		fullname = *model.FullName
	}
	if err != nil {
		return &pb.GetUserResponse{}, err
	}
	return &pb.GetUserResponse{
		Role:     model.Role,
		Phone:    model.Phone,
		FullName: fullname,
		Email:    model.Email,
	}, nil
}
