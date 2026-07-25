package service

import (
	"context"
	"errors"
	"resource-service/internal/filters"
	"resource-service/internal/repository"
	"resource-service/internal/repository/model"

	"gorm.io/gorm"
)

// DriverService — бизнес-слой для работы с водителями.
type DriverService struct {
	repo *repository.DriverRepository
}

// NewDriverService создаёт новый экземпляр сервиса.
func NewDriverService(repo *repository.DriverRepository) *DriverService {
	return &DriverService{repo: repo}
}

// CreateDriver создаёт нового водителя.
// Можно добавить валидацию, проверку дубликатов и т.д.
func (s *DriverService) CreateDriver(ctx context.Context, driver *model.Driver) error {
	if driver == nil {
		return errors.New("driver cannot be nil")
	}
	// Здесь можно добавить бизнес-правила, например, проверку, что номер телефона уникален.
	return s.repo.NewDriver(ctx, driver)
}

// CreateDriverCategory добавляет категорию для водителя.
func (s *DriverService) CreateDriverCategory(ctx context.Context, driverCategory *model.DriverCategory) error {
	if driverCategory == nil {
		return errors.New("driver category cannot be nil")
	}
	// Можно дополнительно проверить, что водитель с таким driver_id существует,
	// но это может делать репозиторий или база данных через внешний ключ.
	return s.repo.NewDriverCategory(ctx, driverCategory)
}

// GetDriverByID возвращает водителя по ID.
func (s *DriverService) GetDriverByID(ctx context.Context, id int) (*model.Driver, error) {
	if id <= 0 {
		return nil, errors.New("invalid driver id")
	}
	driver, err := s.repo.GetDriverByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("driver not found")
		}
		return nil, err
	}
	return driver, nil
}

// GetDriverCategories возвращает все категории водителя.
func (s *DriverService) GetDriverCategories(ctx context.Context, driverID int) ([]model.DriverCategory, error) {
	if driverID <= 0 {
		return nil, errors.New("invalid driver id")
	}
	categoriesPtr, err := s.repo.GetDriverCategoriesByDriverID(ctx, driverID)
	if err != nil {
		return nil, err
	}
	if categoriesPtr == nil {
		return []model.DriverCategory{}, nil // возвращаем пустой слайс, а не nil
	}
	return *categoriesPtr, nil
}

// GetDriversByFilter возвращает список водителей по фильтрам.
func (s *DriverService) GetDriversByFilter(ctx context.Context, filter filters.DriverFilter) ([]model.Driver, error) {
	driversPtr, err := s.repo.GetDriversByFilter(ctx, filter)
	if err != nil {
		return nil, err
	}
	if driversPtr == nil {
		return []model.Driver{}, nil
	}
	return *driversPtr, nil
}

// UpdateDriverStatus обновляет статус водителя.
func (s *DriverService) UpdateDriverStatus(ctx context.Context, id int, status int) error {
	if id <= 0 {
		return errors.New("invalid driver id")
	}
	// Здесь можно добавить валидацию допустимых статусов, например, только 1, 2, 3.
	validStatuses := map[int]bool{1: true, 2: true, 3: true}
	if !validStatuses[status] {
		return errors.New("invalid status value")
	}
	return s.repo.DriverStatusUpdate(ctx, id, status)
}
