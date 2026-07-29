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

func (s *UserService) GetUser(ctx context.Context, id int) (*model.User, error) {
	return s.UserRepository.GetByID(ctx, id)
}

func (s *UserService) GetRoleByID(ctx context.Context, id int) (int, error) {
	user, err := s.UserRepository.GetByID(ctx, id)
	if err != nil {
		return -1, err
	}
	return *user.Role, err
}

func (s *UserService) UpdateRole(ctx context.Context, id, role int) error {
	return s.UserRepository.UpdateRole(ctx, id, role)
}

func (s *UserService) CreateUser(ctx context.Context, id int) error {
	return s.UserRepository.Create(ctx, &model.User{ID: id})
}
