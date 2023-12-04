package service

import (
	"errors"
	"fmt"

	"go-ticketos/internal/domain"
)

type eventService struct {
	eventRepo domain.EventRepo
}

// nolint: revive
func NewEventService(
	eventRepo domain.EventRepo,
) (*eventService, error) {
	if eventRepo == nil {
		return nil, errors.New("NewEventService: eventRepo in nil")
	}
	return &eventService{
		eventRepo: eventRepo,
	}, nil
}

var _ domain.EventService = (*eventService)(nil)

func (e eventService) Create(props domain.CreateEventProps) (*domain.Event, error) {
	event := domain.NewEvent(
		props.Name,
		props.Description,
		props.Place,
		props.AgeRating,
		props.StartAt,
		props.EndAt,
	)
	if err := e.eventRepo.Create(event); err != nil {
		return nil, fmt.Errorf("Create: can not insert event in db: %w", err)
	}
	return &event, nil
}
