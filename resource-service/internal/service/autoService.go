package service

import (
	"context"
	"resource-service/internal/filters"
	"resource-service/internal/repository"
	"resource-service/internal/repository/model"
)

type AutoService interface {
	GetAuto(ctx context.Context, ID int) (*model.Auto, error)
	GetAutosWithFilter(ctx context.Context, filter filters.AutoFilter) (*[]model.Auto, error)
	UpdateStatusAuto(ctx context.Context, status, autoID int) error
	CreateAuto(ctx context.Context, auto *model.Auto) error
}
type autoService struct {
	autoRepository repository.AutoRepository
}

func NewAutoService(AutoRepository repository.AutoRepository) *autoService {
	return &autoService{autoRepository: AutoRepository}
}

func (s *autoService) GetAuto(ctx context.Context, ID int) (*model.Auto, error) {
	return s.autoRepository.GetAutoByID(ctx, ID)
}
func (s *autoService) GetAutosWithFilter(ctx context.Context, filter filters.AutoFilter) (*[]model.Auto, error) {
	return s.autoRepository.GetAutosByFilter(ctx, filter)
}
func (s *autoService) UpdateStatusAuto(ctx context.Context, status, autoID int) error {
	return s.autoRepository.UpdateStatusAuto(ctx, autoID, status)
}
func (s *autoService) CreateAuto(ctx context.Context, auto *model.Auto) error {
	return s.autoRepository.CreateAuto(ctx, auto)
}
