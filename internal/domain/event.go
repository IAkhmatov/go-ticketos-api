package domain

import (
	"time"

	"github.com/google/uuid"
)

type (
	Event struct {
		ID          uuid.UUID
		Name        string
		Description *string
		Place       string
		AgeRating   int
		StartAt     time.Time
		EndAt       time.Time
		CreatedAt   time.Time
		UpdatedAt   time.Time
	}

	CreateEventProps struct {
		Name        string
		Description *string
		Place       string
		AgeRating   int
		StartAt     time.Time
		EndAt       time.Time
	}

	EventRepo interface {
		Create(event Event) error
	}

	EventService interface {
		Create(props CreateEventProps) (*Event, error)
	}
)

func NewEvent(
	name string,
	description *string,
	place string,
	ageRating int,
	startAt time.Time,
	endAt time.Time,
) Event {
	now := time.Now().UTC()
	event := Event{
		ID:        uuid.New(),
		Name:      name,
		Place:     place,
		AgeRating: ageRating,
		StartAt:   startAt,
		EndAt:     endAt,
		CreatedAt: now,
		UpdatedAt: now,
	}
	if description != nil {
		desc := *description
		event.Description = &desc
	}
	return event
}
