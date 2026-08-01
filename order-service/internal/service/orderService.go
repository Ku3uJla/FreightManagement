package service

import (
	"context"
	"errors"
	"log"
	"order-service/dto"
	"order-service/events"
	"order-service/internal/repository"
	"order-service/internal/repository/model"
	"strconv"
	"time"
)

type OrderPublisher interface {
	PublishOrderCreated(ctx context.Context, payload events.OrderPayload) error
	PublishOrderUpdated(ctx context.Context, payload events.OrderPayload) error
	PublishOrderCanceled(ctx context.Context, payload events.OrderPayload) error
}

type OrderService interface {
	CreateOrder(ctx context.Context, userID int, req *dto.CreateOrderRequest) (*model.Order, error)
	GetOrderByID(ctx context.Context, id int) (*model.Order, error)
	UpdateOrder(ctx context.Context, id int, req *dto.UpdateOrderRequest) error
	UpdateStatus(ctx context.Context, id int, status int) error
	AssignManager(ctx context.Context, orderID int, managerID int) error
	GetOrders(ctx context.Context, filter *dto.OrderFilter, page, pageSize int) ([]model.Order, int64, error)
	GetOrdersByUser(ctx context.Context, userID int, page, pageSize int) ([]model.Order, int64, error)
	GetOrdersByDriver(ctx context.Context, driverID int, page, pageSize int) ([]model.Order, error)
}
type orderService struct {
	repo      repository.OrderRepository
	publisher OrderPublisher
}

func NewOrderService(repo repository.OrderRepository, pub OrderPublisher) *orderService {
	return &orderService{
		repo:      repo,
		publisher: pub,
	}
}

func (s *orderService) CreateOrder(ctx context.Context, userID int, req *dto.CreateOrderRequest) (*model.Order, error) {
	if req == nil {
		return nil, errors.New("request cannot be nil")
	}

	// Валидация обязательных полей
	if req.Capacity <= 0 {
		return nil, errors.New("capacity must be greater than 0")
	}
	if req.LiftingCapacity <= 0 {
		return nil, errors.New("lifting_capacity must be greater than 0")
	}
	if req.Price != nil && *req.Price < 0 {
		return nil, errors.New("price cannot be negative")
	}

	// Преобразуем DTO в модель
	order := model.Order{
		UserID:          userID,
		Capacity:        req.Capacity,
		LiftingCapacity: req.LiftingCapacity,
		Status:          1, // default статус
	}
	if req.Type != nil {
		order.Type = req.Type
	}
	if req.Price != nil {
		order.Price = req.Price
	}
	if req.PickupAddress != nil {
		order.PickupAddress = req.PickupAddress
	}
	if req.DeliveryAddress != nil {
		order.DeliveryAddress = req.DeliveryAddress
	}
	if req.DateStart != nil {
		order.DateStart = req.DateStart
	}
	if req.DateEnd != nil {
		order.DateEnd = req.DateEnd
	}
	if req.PeriodFrom != nil {
		order.PeriodFrom = req.PeriodFrom
	}
	if req.PeriodTo != nil {
		order.PeriodTo = req.PeriodTo
	}

	// 1. Сохраняем заказ в БД
	err := s.repo.CreateOrder(ctx, order)
	if err != nil {
		return nil, err
	}

	// 2. Отправляем событие order.created в RabbitMQ 👈
	if s.publisher != nil {
		var price float64
		if order.Price != nil {
			price = *order.Price
		}

		payload := events.OrderPayload{
			OrderID:   string(rune(order.ID)), // Или strconv.Itoa(order.ID)
			UserID:    string(rune(userID)),
			Price:     price,
			Status:    "created",
			Timestamp: time.Now(),
		}

		if pubErr := s.publisher.PublishOrderCreated(ctx, payload); pubErr != nil {
			log.Printf("[WARN] Заказ %d создан, но событие order.created не отправлено: %v", order.ID, pubErr)
		}
	}

	return &order, nil
}

func (s *orderService) GetOrderByID(ctx context.Context, id int) (*model.Order, error) {
	if id <= 0 {
		return nil, errors.New("id must be positive")
	}
	return s.repo.GetOrderByID(ctx, id)
}

func (s *orderService) UpdateOrder(ctx context.Context, id int, req *dto.UpdateOrderRequest) error {
	if id <= 0 {
		return errors.New("id must be positive")
	}
	if req == nil {
		return errors.New("update request cannot be nil")
	}

	if req.Price != nil && *req.Price < 0 {
		return errors.New("price cannot be negative")
	}

	return s.repo.UpdateOrder(ctx, id, req)
}
func (s *orderService) UpdateStatus(ctx context.Context, id int, status int) error {
	if id <= 0 {
		return errors.New("id must be positive")
	}
	// Разрешаем статусы: -1 (отмена) и 1..5
	if status != -1 && (status < 1 || status > 5) {
		return errors.New("status must be -1 or between 1 and 5")
	}

	// 1. Обновляем статус в БД
	err := s.repo.UpdateStatus(ctx, id, status)
	if err != nil {
		return err
	}

	// 2. Отправляем событие в RabbitMQ
	if s.publisher != nil {
		// Формируем payload (исправляем преобразование)
		payload := events.OrderPayload{
			OrderID:   strconv.Itoa(id), // правильное преобразование int → string
			Status:    strconv.Itoa(status),
			Timestamp: time.Now(),
		}

		// Выбираем тип события
		var pubErr error
		if status == -1 {
			pubErr = s.publisher.PublishOrderCanceled(ctx, payload)
		} else {
			pubErr = s.publisher.PublishOrderUpdated(ctx, payload)
		}

		if pubErr != nil {
			log.Printf("[WARN] Событие для заказа %d не отправлено: %v", id, pubErr)
		}
	}

	return nil
}

func (s *orderService) AssignManager(ctx context.Context, orderID int, managerID int) error {
	if orderID <= 0 {
		return errors.New("order_id must be positive")
	}
	return s.repo.AssignManager(ctx, orderID, managerID)
}

func (s *orderService) GetOrders(ctx context.Context, filter *dto.OrderFilter, page, pageSize int) ([]model.Order, int64, error) {
	if page < 1 {
		page = 1
	}
	if filter == nil {
		filter = &dto.OrderFilter{}
	}
	return s.repo.GetList(ctx, filter, page, pageSize)
}

func (s *orderService) GetOrdersByUser(ctx context.Context, userID int, page, pageSize int) ([]model.Order, int64, error) {
	if userID <= 0 {
		return nil, 0, errors.New("user_id must be positive")
	}
	if page < 1 {
		page = 1
	}
	return s.repo.GetOrdersByUser(ctx, userID, page, pageSize)
}

func (s *orderService) GetOrdersByDriver(ctx context.Context, driverID int, page, pageSize int) ([]model.Order, error) {
	if driverID <= 0 {
		return nil, errors.New("driver_id must be positive")
	}
	if page < 1 {
		page = 1
	}

	return s.repo.GetOrdersByDriver(ctx, driverID, page, pageSize)
}
