package service

import (
	"auth-service/internal/features"
	"auth-service/internal/repository"
	"auth-service/internal/repository/model"
	pb "auth-service/proto/userpb"
	"context"
	"errors"
)

type AuthService interface {
	LoginByLogin(ctx context.Context, login, password string) (int, error)
	LoginByEmail(ctx context.Context, email, password string) (int, error)
	SignUp(ctx context.Context, login, email, password string) error
}
type authService struct {
	authRepo   *repository.AuthRepository
	userClient pb.UserServiceClient
}

func NewAuthService(authRepository *repository.AuthRepository, userClient pb.UserServiceClient) *authService {
	return &authService{authRepo: authRepository,
		userClient: userClient}
}

func (s *authService) LoginByEmail(ctx context.Context, email, password string) (int, error) {
	user, err := s.authRepo.GetByEmail(ctx, email)
	if err != nil {
		return -1, err
	}
	if !features.ComparePassword(password, user.Password) {
		return -1, errors.New("Wrong login or password")
	}

	return user.ID, nil
}

func (s *authService) LoginByLogin(ctx context.Context, login, password string) (int, error) {
	user, err := s.authRepo.GetByLogin(ctx, login)
	if err != nil {
		return -1, err
	}
	if !features.ComparePassword(password, user.Password) {
		return -1, errors.New("Wrong login or password")
	}

	return user.ID, nil
}

func (s *authService) SignUp(ctx context.Context, login, email, password string) error {
	user := &model.Auth{
		Login: login,
		Email: email,
	}

	if s.authRepo.ExistsByEmail(ctx, email) {
		return errors.New("Почта уже занята")
	}
	hashedpassword, err := features.HashPassword(password)
	if err != nil {
		return err
	}
	user.Password = hashedpassword
	err = s.authRepo.Create(ctx, user)
	if err != nil {
		return err
	}
	status, err := s.userClient.CreateUser(ctx, &pb.CreateUserRequest{Id: int64(user.ID)})
	if !status.Success {
		return errors.New("Ошибка передачи")
	}
	if err != nil {
		return err
	}
	return nil
}
