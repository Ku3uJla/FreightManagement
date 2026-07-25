package service

import (
	"context"
	"errors"
	"order-service/internal/dto"
	"order-service/internal/repository"
	"order-service/internal/repository/model"
)

type OrderService struct {
	repo *repository.OrderRepository
}

func NewOrderService(repo *repository.OrderRepository) *OrderService {
	return &OrderService{repo: repo}
}

func (s *OrderService) CreateOrder(ctx context.Context, req *dto.CreateOrderRequest) (*model.Order, error) {
	if req == nil {
		return nil, errors.New("request cannot be nil")
	}

	// Валидация обязательных полей
	if req.UserID <= 0 {
		return nil, errors.New("user_id is required and must be positive")
	}
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
		UserID:          req.UserID,
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

	// Сохраняем
	err := s.repo.CreateOrder(ctx, order)
	if err != nil {
		return nil, err
	}
	return &order, nil
}

func (s *OrderService) GetOrderByID(ctx context.Context, id int) (*model.Order, error) {
	if id <= 0 {
		return nil, errors.New("id must be positive")
	}
	return s.repo.GetOrderByID(ctx, id)
}

func (s *OrderService) UpdateOrder(ctx context.Context, id int, req *dto.UpdateOrderRequest) error {
	if id <= 0 {
		return errors.New("id must be positive")
	}
	if req == nil {
		return errors.New("update request cannot be nil")
	}

	// Валидация цены, если передана
	if req.Price != nil && *req.Price < 0 {
		return errors.New("price cannot be negative")
	}

	return s.repo.UpdateOrder(ctx, id, req)
}

func (s *OrderService) UpdateStatus(ctx context.Context, id int, status int) error {
	if id <= 0 {
		return errors.New("id must be positive")
	}
	if status < 1 || status > 5 {
		return errors.New("status must be between 1 and 5")
	}
	return s.repo.UpdateStatus(ctx, id, status)
}

func (s *OrderService) AssignManager(ctx context.Context, orderID int, managerID int) error {
	if orderID <= 0 {
		return errors.New("order_id must be positive")
	}
	return s.repo.AssignManager(ctx, orderID, managerID)
}

func (s *OrderService) GetOrders(ctx context.Context, filter *dto.OrderFilter, page, pageSize int) ([]model.Order, int64, error) {
	if page < 1 {
		page = 1
	}
	if filter == nil {
		filter = &dto.OrderFilter{}
	}
	return s.repo.GetList(ctx, filter, page, pageSize)
}

func (s *OrderService) GetOrdersByUser(ctx context.Context, userID int, page, pageSize int) ([]model.Order, int64, error) {
	if userID <= 0 {
		return nil, 0, errors.New("user_id must be positive")
	}
	if page < 1 {
		page = 1
	}
	return s.repo.GetOrdersByUser(ctx, userID, page, pageSize)
}

func (s *OrderService) GetOrdersByDriver(ctx context.Context, driverID int, page, pageSize int) ([]model.Order, error) {
	if driverID <= 0 {
		return nil, errors.New("driver_id must be positive")
	}
	if page < 1 {
		page = 1
	}

	return s.repo.GetOrdersByDriver(ctx, driverID, page, pageSize)
}
