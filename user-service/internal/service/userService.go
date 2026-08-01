package service

import (
	"context"
	"log"
	"strconv"
	"time"
	"user-service/dto"
	"user-service/events"
	"user-service/internal/repository"
	"user-service/internal/repository/model"
)

type UserEventPublisher interface {
	PublishUserCreated(ctx context.Context, payload events.UserCreatedPayload) error
}

type UserService interface {
	GetUser(ctx context.Context, id int) (*model.User, error)
	GetRoleByID(ctx context.Context, id int) (string, error)
	UpdateRole(ctx context.Context, id, role int) error
	CreateUser(ctx context.Context, user dto.CreateUserGrpcRequest) error
}

type userService struct {
	userRepo      repository.UserRepository
	userPublisher UserEventPublisher
}

func NewUserService(UserRepository repository.UserRepository, pub UserEventPublisher) *userService {
	return &userService{UserRepository, pub}
}

func (s *userService) GetUser(ctx context.Context, id int) (*model.User, error) {
	return s.userRepo.GetByID(ctx, id)
}

func (s *userService) GetRoleByID(ctx context.Context, id int) (string, error) {
	user, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		return "", err
	}
	return user.Role, err
}

func (s *userService) UpdateRole(ctx context.Context, id, role int) error {
	return s.userRepo.UpdateRole(ctx, id, role)
}

func (s *userService) CreateUser(ctx context.Context, user dto.CreateUserGrpcRequest) error {
	log.Print("CreateUser был вызван", user)
	// 1. Создаем модель для БД
	newUser := &model.User{
		ID:       user.ID,
		Role:     "User",
		Phone:    user.Phone,
		Email:    user.Email,
		FullName: &user.FullName,
	}

	// 2. Сохраняем пользователя в Базу Данных
	if err := s.userRepo.Create(ctx, newUser); err != nil {
		return err
	}

	// 3. Формируем DTO для брокера сообщений
	payload := events.UserCreatedPayload{
		UserID:    strconv.Itoa(user.ID),
		Email:     user.Email,
		Name:      user.FullName,
		Role:      "User",
		CreatedAt: time.Now(),
	}

	if err := s.userPublisher.PublishUserCreated(ctx, payload); err != nil {
		log.Printf("[WARN] Пользователь %s создан в БД, но событие user.created не отправлено: %v", payload.UserID, err)
	}

	return nil
}
