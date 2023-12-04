package transport

import (
	"go-ticketos/internal/domain"
)

type EventAdapter interface {
	ToResponseDTO(event domain.Event) EventResponseDTO
}

type eventAdapter struct{}

var _ EventAdapter = (*eventAdapter)(nil)

// nolint: revive
func NewEventAdapter() *eventAdapter {
	return &eventAdapter{}
}
func (a *eventAdapter) ToResponseDTO(event domain.Event) EventResponseDTO {
	return EventResponseDTO{
		ID:          event.ID,
		Name:        event.Name,
		Description: event.Description,
		Place:       event.Place,
		AgeRating:   event.AgeRating,
		StartAt:     event.StartAt,
		EndAt:       event.EndAt,
	}
}
