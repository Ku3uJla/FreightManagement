package service

import (
	"context"
	"fmt"
	"time"

	"strconv"
	"trips-service/dto"
	"trips-service/internal/repository"
	"trips-service/internal/repository/model"
	"trips-service/rabbitmq"
)

type TripService interface {
	CreateTrip(ctx context.Context, trip *model.Trip) error
	GetTripByID(ctx context.Context, id int) (*model.Trip, error)
	GetTrips(ctx context.Context, filter *dto.TripFilter, page, pageSize int) ([]model.Trip, int64, error)
	AssignDriver(ctx context.Context, id int, input *dto.UpdateTripInput) error
	AssignAuto(ctx context.Context, id int, input *dto.UpdateTripInput) error
	ChangeTripStatus(ctx context.Context, tripID int, newStatus string) error
}

type tripService struct {
	repo     repository.TripsRepository
	producer rabbitmq.TripEventProducer
}

func NewTripsService(repo repository.TripsRepository, producer rabbitmq.TripEventProducer) TripService {
	return &tripService{
		repo:     repo,
		producer: producer,
	}
}

// CreateTrip инициализирует новый рейс
func (s *tripService) CreateTrip(ctx context.Context, trip *model.Trip) error {
	if err := s.repo.CreateTrip(ctx, trip); err != nil {
		return fmt.Errorf("service.CreateTrip: %w", err)
	}

	return nil
}

func (s *tripService) GetTripByID(ctx context.Context, id int) (*model.Trip, error) {
	trip, err := s.repo.GetTripByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("service.GetTripByID: %w", err)
	}
	if trip == nil {
		return nil, fmt.Errorf("trip with id %v not found", id)
	}
	return trip, nil
}

func (s *tripService) GetTrips(ctx context.Context, filter *dto.TripFilter, page, pageSize int) ([]model.Trip, int64, error) {
	return s.repo.GetTrips(ctx, filter, page, pageSize)
}

func (s *tripService) AssignDriver(ctx context.Context, id int, input *dto.UpdateTripInput) error {
	_, err := s.GetTripByID(ctx, id)
	if err != nil {
		return err
	}
	return s.repo.UpdateTrip(ctx, id, input)
}

func (s *tripService) AssignAuto(ctx context.Context, id int, input *dto.UpdateTripInput) error {
	_, err := s.GetTripByID(ctx, id)
	if err != nil {
		return err
	}
	return s.repo.UpdateTrip(ctx, id, input)
}

func (s *tripService) ChangeTripStatus(ctx context.Context, tripID int, newStatus string) error {
	trip, err := s.GetTripByID(ctx, tripID)
	if err != nil {
		return err
	}

	if !isValidStatusTransition(trip.Status, newStatus) {
		return fmt.Errorf("invalid status transition from %s to %s", trip.Status, newStatus)
	}

	oldStatus := trip.Status

	if err := s.repo.UpdateTripStatus(ctx, tripID, newStatus); err != nil {
		return fmt.Errorf("failed to update status in DB: %w", err)
	}

	event := dto.TripStatusUpdatedPayload{
		TripID:    strconv.Itoa(trip.ID),
		OrderID:   strconv.Itoa(trip.OrderID),
		DriverID:  strconv.Itoa(*trip.DriverID),
		AutoID:    strconv.Itoa(*trip.AutoID),
		OldStatus: oldStatus,
		NewStatus: newStatus,
		UpdatedAt: time.Now(),
	}

	if err := s.producer.PublishStatusChanged(ctx, event); err != nil {

		return fmt.Errorf("status updated in DB, but failed to publish event to RabbitMQ: %w", err)
	}

	return nil
}

func isValidStatusTransition(currentStatus, newStatus string) bool {
	allowedTransitions := map[string][]string{
		"ASSIGNED":     {"AT_LOADING", "CANCELLED"},
		"AT_LOADING":   {"IN_TRANSIT", "CANCELLED", "FAILED"},
		"IN_TRANSIT":   {"AT_UNLOADING", "FAILED"},
		"AT_UNLOADING": {"COMPLETED", "FAILED"},
		"COMPLETED":    {}, // Финальный статус
		"CANCELLED":    {}, // Финальный статус
		"FAILED":       {}, // Финальный статус
	}

	validNextStatuses, exists := allowedTransitions[currentStatus]
	if !exists {
		return false
	}

	for _, status := range validNextStatuses {
		if status == newStatus {
			return true
		}
	}

	return false
}
