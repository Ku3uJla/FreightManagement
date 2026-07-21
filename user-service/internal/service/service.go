package service

import (
	"context"
	"errors"
	"strconv"
	"user-service/internal/features"
	"user-service/internal/repository"
	"user-service/internal/repository/model"
)

type UserService struct {
	userRepository *repository.UserRepository
}

func NewUserService(UserRepository *repository.UserRepository) *UserService {
	return &UserService{UserRepository}
}

func (s *UserService) Login(ctx context.Context, email, password string) (string, error) {
	if !s.userRepository.ExistsByEmail(ctx, email) {
		return "", errors.New("Почта указана не верно")
	}
	user, err := s.userRepository.GetByEmail(ctx, email)
	if err != nil {
		return "", err
	}
	if features.ComparePassword(password, user.Password) {
		return "", errors.New("Неверный пароль")
	}

	ID := strconv.Itoa(user.User_ID)
	token, err := features.CreateToken(ID)
	if err != nil {
		return "", err
	}
	return token, nil
}

func (s *UserService) SignUp(ctx context.Context, email, password string) error {
	if s.userRepository.ExistsByEmail(ctx, email) {
		return errors.New("Почта уже занята")
	}
	hashedpassword, err := features.HashPassword(password)
	if err != nil {
		return err
	}
	err = s.userRepository.Create(ctx, &model.User{Email: email, Password: hashedpassword})
	if err != nil {
		return err
	}
	return nil
}

func (s *UserService) GetUser(ctx context.Context, id string) (*model.User, error) {
	return s.userRepository.GetByID(ctx, id)
}
