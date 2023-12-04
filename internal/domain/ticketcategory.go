package domain

import (
	"time"

	"github.com/google/uuid"
)

type (
	TicketCategory struct {
		ID          uuid.UUID
		EventID     uuid.UUID
		Price       uint
		Name        string
		Description *string
		CreatedAt   time.Time
		UpdatedAt   time.Time
	}

	TicketCategoryRepo interface {
		Create(tc TicketCategory) error
		GetByID(id uuid.UUID) (*TicketCategory, error)
	}
)

func NewTicketCategory(
	eventID uuid.UUID,
	price uint,
	name string,
	description *string,
) TicketCategory {
	now := time.Now().UTC()
	tc := TicketCategory{
		ID:        uuid.New(),
		EventID:   eventID,
		Price:     price,
		Name:      name,
		CreatedAt: now,
		UpdatedAt: now,
	}
	if description != nil {
		desc := *description
		tc.Description = &desc
	}
	return tc
}

func (tc TicketCategory) CanUsePromocode(p Promocode) bool {
	for _, ptc := range p.TicketCategories {
		if tc.ID == ptc.ID {
			return true
		}
	}
	return false
}
