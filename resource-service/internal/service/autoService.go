package service

import (
	"context"
	filters "user-service/internal/filters"
	"user-service/internal/repository"
	"user-service/internal/repository/model"
)

type AutoService struct {
	autoRepository *repository.AutoRepository
}

func NewAutoService(AutoRepository *repository.AutoRepository) *AutoService {
	return &AutoService{autoRepository: AutoRepository}
}

func (s *AutoService) GetAuto(ctx context.Context, ID int) (*model.Auto, error) {
	return s.autoRepository.GetAutoByID(ctx, ID)
}
func (s *AutoService) GetAutosWithFilter(ctx context.Context, filter filters.AutoFilter) (*[]model.Auto, error) {
	return s.autoRepository.GetAutosByFilter(ctx, filter)
}
func (s *AutoService) UpdateStatusAuto(ctx context.Context, status, autoID int) error {
	return s.autoRepository.UpdateStatusAuto(ctx, autoID, status)
}
func (s *AutoService) CreateAuto(ctx context.Context, auto *model.Auto) error {
	return s.autoRepository.CreateAuto(ctx, *auto)
}
