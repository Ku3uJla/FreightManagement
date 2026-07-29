package service

import (
	"context"
	"errors"
	"strconv"
	"user-service/internal/features"
	"user-service/internal/repository"
	"user-service/internal/repository/model"
)

type AuthService struct {
	userRepository *repository.UserRepository
}

func NewAuthService(UserRepository *repository.UserRepository) *AuthService {
	return &AuthService{UserRepository}
}

func (s *AuthService) Login(ctx context.Context, email, password string) (string, int, error) {
	if !s.userRepository.ExistsByEmail(ctx, email) {
		return "", -1, errors.New("Почта указана не верно")
	}
	user, err := s.userRepository.GetByEmail(ctx, email)
	if err != nil {
		return "", -1, err
	}
	if features.ComparePassword(password, user.Password) {
		return "", -1, errors.New("Неверный пароль")
	}

	ID := strconv.Itoa(user.ID)
	return ID, user.Role, nil
}

func (s *AuthService) SignUp(ctx context.Context, login, email, password string) error {
	if s.userRepository.ExistsByEmail(ctx, email) {
		return errors.New("Почта уже занята")
	}
	hashedpassword, err := features.HashPassword(password)
	if err != nil {
		return err
	}
	err = s.userRepository.Create(ctx, &model.User{Login: login, Email: email, Password: hashedpassword})
	if err != nil {
		return err
	}
	return nil
}
