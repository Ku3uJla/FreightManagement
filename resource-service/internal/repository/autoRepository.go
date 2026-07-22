package repository

import (
	"context"
	"user-service/internal/repository/model"

	"gorm.io/gorm"
)

type AutoFilter struct {
	Capacity        *int
	LiftingCapacity *int
	Status          *int
}

type AutoRepository struct {
	db *gorm.DB
}

func NewAutoRepository(db *gorm.DB) *AutoRepository {
	return &AutoRepository{db: db}
}

func (r *AutoRepository) GetAutosID(ctx context.Context) ([]int, error) {
	var autoIDs []int
	result := r.db.WithContext(ctx).Raw("SELECT id FROM auto").Scan(&autoIDs)
	if result.Error != nil {
		return nil, result.Error
	}
	return autoIDs, nil
}

func (r *AutoRepository) GetAutoByID(ctx context.Context, id int) (*model.Auto, error) {
	var auto model.Auto
	err := r.db.WithContext(ctx).Model(&model.Auto{}).Where(id).First(&auto).Error
	if err != nil {
		return nil, err
	}
	return &auto, err
}

func (r *AutoRepository) GetAutosByFilter(ctx context.Context, filter AutoFilter) ([]model.Auto, error) {
	var autos []model.Auto
	query := r.db.WithContext(ctx).Model(&model.Auto{})
	if filter.Capacity != nil {
		query.Where("capacity > ?", *filter.Capacity)
	}
	if filter.LiftingCapacity != nil {
		query.Where("lifting_capacity > ?", *filter.LiftingCapacity)
	}
	if *filter.Status == 3 {
		query.Where("status = ?", *filter.Status)
	}
	err := query.Find(&autos).Error
	if err != nil {
		return nil, err
	}
	return autos, nil
}

func (r *AutoRepository) UpdateStatusAuto(ctx context.Context, autoID int, status int) error {
	err := r.db.Model(&model.Auto{}).Where("id = ?", autoID).Update("status", status).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *AutoRepository) CreateAuto(ctx context.Context, auto model.Auto) error {
	err := r.db.Create(auto).Error
	if err != nil {
		return err
	}
	return nil
}
