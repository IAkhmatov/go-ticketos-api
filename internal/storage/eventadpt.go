package storage

import (
	"database/sql"

	"go-ticketos/internal/domain"
)

type EventAdapter interface {
	ToSchema(event domain.Event) eventSchema
}

type eventAdapter struct{}

var _ EventAdapter = (*eventAdapter)(nil)

// nolint: revive
func NewEventAdapter() *eventAdapter {
	return &eventAdapter{}
}

func (a eventAdapter) ToSchema(event domain.Event) eventSchema {
	schema := eventSchema{
		ID:        event.ID,
		Name:      event.Name,
		Place:     event.Place,
		AgeRating: event.AgeRating,
		StartAt:   event.StartAt,
		EndAt:     event.EndAt,
		CreatedAt: event.CreatedAt,
		UpdatedAt: event.UpdatedAt,
	}
	if event.Description != nil {
		schema.Description = sql.NullString{
			String: *event.Description,
			Valid:  true,
		}
	}
	return schema
}
