package repository

import (
	"context"
	"errors"
	"order-service/dto"
	"order-service/internal/repository/model"

	"gorm.io/gorm"
)

type OrderRepository interface {
	CreateOrder(ctx context.Context, order model.Order) error
	UpdateOrder(ctx context.Context, id int, req *dto.UpdateOrderRequest) error
	GetOrderByID(ctx context.Context, id int) (*model.Order, error)
	GetList(ctx context.Context, filter *dto.OrderFilter, page, pageSize int) ([]model.Order, int64, error)
	UpdateStatus(ctx context.Context, id int, status int) error
	AssignManager(ctx context.Context, orderID int, managerID int) error
	GetOrdersByUser(ctx context.Context, userID int, page, pageSize int) ([]model.Order, int64, error)
	GetOrdersByDriver(ctx context.Context, driverID int, page, pageSize int) ([]model.Order, error)
}

type orderRepository struct {
	db *gorm.DB
}

func NewOrderRepository(db *gorm.DB) *orderRepository {
	return &orderRepository{db}
}

func (r *orderRepository) CreateOrder(ctx context.Context, order model.Order) error {
	return r.db.Model(&model.Order{}).Create(order).Error
}

func (r *orderRepository) UpdateOrder(ctx context.Context, id int, req *dto.UpdateOrderRequest) error {
	if req == nil {
		return errors.New("update request is empty")
	}
	return r.db.WithContext(ctx).Model(&model.Order{}).Where("id = ?", id).Updates(req).Error
}

func (r *orderRepository) GetOrderByID(ctx context.Context, id int) (*model.Order, error) {
	var order model.Order
	err := r.db.WithContext(ctx).Model(&model.Order{}).Where("id = ?", id).First(&order).Error
	if err != nil {
		return &model.Order{}, err
	}
	return &order, nil
}

// GetList с фильтрацией и пагинацией
func (r *orderRepository) GetList(ctx context.Context, filter *dto.OrderFilter, page, pageSize int) ([]model.Order, int64, error) {
	var orders []model.Order
	var total int64

	query := r.db.WithContext(ctx).Model(&model.Order{})

	// 1. Основные поля
	if filter.UserID != nil {
		query = query.Where("user_id = ?", *filter.UserID)
	}
	if filter.Status != nil {
		query = query.Where("status = ?", *filter.Status)
	}
	if filter.ManagerID != nil {
		query = query.Where("manager_id = ?", *filter.ManagerID)
	}
	if filter.Type != nil {
		query = query.Where("type = ?", *filter.Type)
	}

	// 2. Диапазоны вместимости
	if filter.CapacityMin != nil {
		query = query.Where("capacity >= ?", *filter.CapacityMin)
	}
	if filter.CapacityMax != nil {
		query = query.Where("capacity <= ?", *filter.CapacityMax)
	}
	if filter.LiftingCapacityMin != nil {
		query = query.Where("lifting_capacity >= ?", *filter.LiftingCapacityMin)
	}
	if filter.LiftingCapacityMax != nil {
		query = query.Where("lifting_capacity <= ?", *filter.LiftingCapacityMax)
	}

	// 3. Диапазоны цены
	if filter.PriceMin != nil {
		query = query.Where("price >= ?", *filter.PriceMin)
	}
	if filter.PriceMax != nil {
		query = query.Where("price <= ?", *filter.PriceMax)
	}

	// 4. Диапазоны дат
	if filter.DateCreateFrom != nil {
		query = query.Where("date_create >= ?", *filter.DateCreateFrom)
	}
	if filter.DateCreateTo != nil {
		query = query.Where("date_create <= ?", *filter.DateCreateTo)
	}
	if filter.DateStartFrom != nil {
		query = query.Where("date_start >= ?", *filter.DateStartFrom)
	}
	if filter.DateStartTo != nil {
		query = query.Where("date_start <= ?", *filter.DateStartTo)
	}
	if filter.DateEndFrom != nil {
		query = query.Where("date_end >= ?", *filter.DateEndFrom)
	}
	if filter.DateEndTo != nil {
		query = query.Where("date_end <= ?", *filter.DateEndTo)
	}
	if filter.PeriodFromFrom != nil {
		query = query.Where("period_from >= ?", *filter.PeriodFromFrom)
	}
	if filter.PeriodFromTo != nil {
		query = query.Where("period_from <= ?", *filter.PeriodFromTo)
	}
	if filter.PeriodToFrom != nil {
		query = query.Where("period_to >= ?", *filter.PeriodToFrom)
	}
	if filter.PeriodToTo != nil {
		query = query.Where("period_to <= ?", *filter.PeriodToTo)
	}

	// 5. Адреса (частичное совпадение, регистронезависимо)
	if filter.PickupAddress != nil && *filter.PickupAddress != "" {
		query = query.Where("pickup_address ILIKE ?", "%"+*filter.PickupAddress+"%")
	}
	if filter.DeliveryAddress != nil && *filter.DeliveryAddress != "" {
		query = query.Where("delivery_address ILIKE ?", "%"+*filter.DeliveryAddress+"%")
	}

	// 6. Связи через driver_auto (JOIN)
	if filter.DriverID != nil || filter.AutoID != nil {
		query = query.Joins("JOIN driver_auto ON driver_auto.order_id = orders.id")
		if filter.DriverID != nil {
			query = query.Where("driver_auto.driver_id = ?", *filter.DriverID)
		}
		if filter.AutoID != nil {
			query = query.Where("driver_auto.auto_id = ?", *filter.AutoID)
		}
		query = query.Distinct()
	}

	offset := (page - 1) * pageSize
	err := query.Limit(pageSize).Offset(offset).Order("id desc").Find(&orders).Error
	return orders, total, err
}

func (r *orderRepository) UpdateStatus(ctx context.Context, id int, status int) error {
	return r.db.WithContext(ctx).
		Model(&model.Order{}).
		Where("id = ?", id).
		Update("status", status).Error
}

func (r *orderRepository) AssignManager(ctx context.Context, orderID int, managerID int) error {
	return r.db.WithContext(ctx).
		Model(&model.Order{}).
		Where("id = ?", orderID).
		Update("manager_id", managerID).Error
}

func (r *orderRepository) GetOrdersByUser(ctx context.Context, userID int, page, pageSize int) ([]model.Order, int64, error) {
	filter := &dto.OrderFilter{
		UserID: &userID,
	}
	return r.GetList(ctx, filter, page, pageSize)
}

func (r *orderRepository) GetOrdersByDriver(ctx context.Context, driverID int, page, pageSize int) ([]model.Order, error) {
	var orders []model.Order
	err := r.db.WithContext(ctx).
		Joins("JOIN driver_auto ON driver_auto.order_id = orders.id").
		Where("driver_auto.driver_id = ?", driverID).
		Order("orders.id desc").
		Limit(pageSize).Offset((page - 1) * pageSize).
		Find(&orders).Error
	return orders, err
}
