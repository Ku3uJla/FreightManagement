package controller

import (
	"auth-service/internal/service"
	pb "auth-service/proto/authpb"
	"context"
)

type AuthGrpcController struct {
	srv service.AuthService
	pb.UnimplementedAuthServiceServer
}

func NewAuthGrpcController(authSrv service.AuthService) *AuthGrpcController {
	return &AuthGrpcController{srv: authSrv}
}

func (c *AuthGrpcController) GetLogin(ctx context.Context, req *pb.GetUserRequest) (*pb.GetLoginResponse, error) {
	login, err := c.srv.GetLogin(ctx, int(req.Id))
	if err != nil {
		return nil, err
	}
	return &pb.GetLoginResponse{Login: login}, nil
}

func (c *AuthGrpcController) GetEmail(ctx context.Context, req *pb.GetUserRequest) (*pb.GetEmailResponse, error) {
	email, err := c.srv.GetEmail(ctx, int(req.Id))
	if err != nil {
		return nil, err
	}
	return &pb.GetEmailResponse{Email: email}, nil
}
