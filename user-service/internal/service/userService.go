package service

import (
	"context"
	"user-service/internal/repository"
	"user-service/internal/repository/model"
)

type UserService struct {
	UserRepository *repository.UserRepository
}

func NewUserService(UserRepository *repository.UserRepository) *UserService {
	return &UserService{UserRepository}
}

func (s *UserService) GetUser(ctx context.Context, id string) (*model.User, error) {
	return s.UserRepository.GetByID(ctx, id)
}

func (s *UserService) GetRoleByID(ctx context.Context, id string) (int, error) {
	user, err := s.UserRepository.GetByID(ctx, id)
	if err != nil {
		return -1, err
	}
	return user.Role, err
}
