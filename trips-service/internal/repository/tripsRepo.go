package repository

import (
	"context"
	"trips-service/dto"
	"trips-service/internal/repository/model"

	"errors"
	"fmt"

	"gorm.io/gorm"
)

// TripsRepository описывает методы работы с БД для рейсов
type TripsRepository interface {
	CreateTrip(ctx context.Context, trip *model.Trip) error
	GetTripByID(ctx context.Context, id int) (*model.Trip, error)
	GetTrips(ctx context.Context, filter *dto.TripFilter, page, pageSize int) ([]model.Trip, int64, error)
	UpdateTrip(ctx context.Context, id int, input *dto.UpdateTripInput) error
	UpdateTripStatus(ctx context.Context, id int, newStatus string) error
	DeleteTrip(ctx context.Context, id int) error
}

type tripsRepository struct {
	db *gorm.DB
}

// NewTripsRepository возвращает реализацию репозитория
func NewTripsRepository(db *gorm.DB) TripsRepository {
	return &tripsRepository{db: db}
}

// CreateTrip создает новую запись о рейсе
func (r *tripsRepository) CreateTrip(ctx context.Context, trip *model.Trip) error {
	if err := r.db.WithContext(ctx).Create(trip).Error; err != nil {
		return fmt.Errorf("failed to create trip: %w", err)
	}
	return nil
}

// GetTripByID получает рейс по его уникальному ID
func (r *tripsRepository) GetTripByID(ctx context.Context, id int) (*model.Trip, error) {
	var trip model.Trip
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&trip).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // Рейс не найден
		}
		return nil, fmt.Errorf("failed to get trip by id: %w", err)
	}
	return &trip, nil
}

func (r *tripsRepository) GetTrips(ctx context.Context, filter *dto.TripFilter, page, pageSize int) ([]model.Trip, int64, error) {
	var trips []model.Trip
	var totalRecords int64

	query := r.db.WithContext(ctx).Model(&model.Trip{})

	// Явно накладываем только те фильтры, которые переданы (не nil)
	if filter != nil {
		if filter.DriverID != nil {
			query = query.Where("driver_id = ?", *filter.DriverID)
		}
		if filter.AutoID != nil {
			query = query.Where("auto_id = ?", *filter.AutoID)
		}
		if len(filter.Statuses) > 0 {
			query = query.Where("status IN ?", filter.Statuses)
		}
		if filter.DateFrom != nil {
			query = query.Where("created_at >= ?", *filter.DateFrom)
		}
		if filter.DateTo != nil {
			query = query.Where("created_at <= ?", *filter.DateTo)
		}
	}

	// Считаем общее количество
	if err := query.Count(&totalRecords).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count trips: %w", err)
	}

	offset := (page - 1) * pageSize
	err := query.
		Order("created_at DESC").
		Limit(pageSize).
		Offset(offset).
		Find(&trips).Error

	return trips, totalRecords, err
}

// UpdateTrip обновляет все заполненные поля рейса (использует Updates)
func (r *tripsRepository) UpdateTrip(ctx context.Context, id int, input *dto.UpdateTripInput) error {
	updates := make(map[string]interface{})

	if input.DriverID != nil {
		updates["driver_id"] = *input.DriverID
	}
	if input.AutoID != nil {
		updates["auto_id"] = *input.AutoID
	}
	if input.Comment != nil {
		updates["comment"] = *input.Comment
	}

	if len(updates) == 0 {
		return nil // Нечего обновлять
	}

	updates["updated_at"] = gorm.Expr("NOW()")

	result := r.db.WithContext(ctx).
		Model(&model.Trip{}).
		Where("id = ?", id).
		Updates(updates)

	if result.Error != nil {
		return fmt.Errorf("failed to update trip: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return errors.New("trip not found")
	}

	return nil
}

// UpdateTripStatus атомарно обновляет статус рейса и сопутствующую временную метку
func (r *tripsRepository) UpdateTripStatus(ctx context.Context, id int, newStatus string) error {
	updates := map[string]interface{}{
		"status":     newStatus,
		"updated_at": gorm.Expr("NOW()"),
	}

	// Автоматически фиксируем время прохождения этапов
	switch newStatus {
	case "AT_LOADING":
		updates["arrived_loading_at"] = gorm.Expr("NOW()")
	case "IN_TRANSIT":
		updates["departed_loading_at"] = gorm.Expr("NOW()")
	case "AT_UNLOADING":
		updates["arrived_unloading_at"] = gorm.Expr("NOW()")
	case "COMPLETED":
		updates["finished_at"] = gorm.Expr("NOW()")
	}

	result := r.db.WithContext(ctx).
		Model(&model.Trip{}).
		Where("id = ?", id).
		Updates(updates)

	if result.Error != nil {
		return fmt.Errorf("failed to update trip status: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return errors.New("trip not found")
	}

	return nil
}

// DeleteTrip выполняет мягкое или физическое удаление рейса
func (r *tripsRepository) DeleteTrip(ctx context.Context, id int) error {
	result := r.db.WithContext(ctx).Where("id = ?", id).Delete(&model.Trip{})
	if result.Error != nil {
		return fmt.Errorf("failed to delete trip: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return errors.New("trip not found")
	}
	return nil
}
