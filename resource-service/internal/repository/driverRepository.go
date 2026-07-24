package repository

import (
	"context"
	"user-service/internal/repository/model"

	"gorm.io/gorm"
)

type DriverFilter struct {
	Status   *int
	Category *string
}

type DriverRepository struct {
	db *gorm.DB
}

func NewDriverRepository(db *gorm.DB) *DriverRepository {
	return &DriverRepository{db}
}

func (r *DriverRepository) NewDriver(ctx context.Context, driver *model.Driver) error {
	err := r.db.Model(&model.Driver{}).Create(driver).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *DriverRepository) NewDriverCategory(ctx *context.Context, driverCategory *model.DriverCategory) error {
	err := r.db.Model(&model.DriverCategory{}).Create(driverCategory).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *DriverRepository) GetDriverByID(ctx context.Context, ID int) (*model.Driver, error) {
	var driver model.Driver
	err := r.db.WithContext(ctx).Model(&model.Driver{}).Where("id = ?", ID).First(&driver).Error
	if err != nil {
		return nil, err
	}
	return &driver, err
}

func (r *DriverRepository) GetDriverCategoriesByDriverID(ctx context.Context, DriverID int) (*[]model.DriverCategory, error) {
	var driverCategories []model.DriverCategory
	err := r.db.WithContext(ctx).Model(&model.DriverCategory{}).Where("driver_id = ?", DriverID).Scan(&driverCategories).Error
	if err != nil {
		return nil, err
	}
	return &driverCategories, err
}

func (r *DriverRepository) GetDriversByFilter(ctx context.Context, filter *DriverFilter) (*[]model.Driver, error) {
	var drivers []model.Driver
	query := r.db.WithContext(ctx).Model(&model.Driver{})
	if filter.Status != nil {
		query.Where("status = ?", *filter.Status)
	}
	if filter.Category != nil {
		query.Joins("JOIN driver_category ON driver_category.driver_id = driver.id").Where("driver_category.category = ?", *filter.Category).Distinct()
	}
	err := query.Scan(&drivers).Error
	if err != nil {
		return nil, err
	}
	return &drivers, nil
}
func (r *DriverRepository) DriverStatusUpdate(ctx context.Context, ID int, status int) error {
	err := r.db.WithContext(ctx).Model(&model.Driver{}).Where("id = ?", ID).Update("status", status).Error
	if err != nil {
		return err
	}
	return nil
}
