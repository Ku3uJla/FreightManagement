package service

import (
	"context"
	"errors"
	"resource-service/internal/filters"
	"resource-service/internal/repository"
	"resource-service/internal/repository/model"

	"gorm.io/gorm"
)

type DriverService interface {
	CreateDriver(ctx context.Context, driver *model.Driver) error
	CreateDriverCategory(ctx context.Context, driverCategory *model.DriverCategory) error
	GetDriverByID(ctx context.Context, id int) (*model.Driver, error)
	GetDriverCategories(ctx context.Context, driverID int) ([]model.DriverCategory, error)
	GetDriversByFilter(ctx context.Context, filter filters.DriverFilter) ([]model.Driver, error)
	UpdateDriverStatus(ctx context.Context, id int, status int) error
}
type driverService struct {
	repo *repository.DriverRepository
}

func NewDriverService(repo *repository.DriverRepository) *driverService {
	return &driverService{repo: repo}
}

func (s *driverService) CreateDriver(ctx context.Context, driver *model.Driver) error {
	if driver == nil {
		return errors.New("driver cannot be nil")
	}
	return s.repo.NewDriver(ctx, driver)
}

func (s *driverService) CreateDriverCategory(ctx context.Context, driverCategory *model.DriverCategory) error {
	if driverCategory == nil {
		return errors.New("driver category cannot be nil")
	}

	return s.repo.NewDriverCategory(ctx, driverCategory)
}

func (s *driverService) GetDriverByID(ctx context.Context, id int) (*model.Driver, error) {
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

func (s *driverService) GetDriverCategories(ctx context.Context, driverID int) ([]model.DriverCategory, error) {
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
func (s *driverService) GetDriversByFilter(ctx context.Context, filter filters.DriverFilter) ([]model.Driver, error) {
	driversPtr, err := s.repo.GetDriversByFilter(ctx, filter)
	if err != nil {
		return nil, err
	}
	if driversPtr == nil {
		return []model.Driver{}, nil
	}
	return *driversPtr, nil
}

func (s *driverService) UpdateDriverStatus(ctx context.Context, id int, status int) error {
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
