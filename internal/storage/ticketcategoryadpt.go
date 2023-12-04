package storage

import (
	"database/sql"

	"go-ticketos/internal/domain"
)

type TicketCategoryAdapter interface {
	ToSchema(tc domain.TicketCategory) ticketCategorySchema
	ToDomain(schema ticketCategorySchema) domain.TicketCategory
}

type ticketCategoryAdapter struct{}

var _ TicketCategoryAdapter = (*ticketCategoryAdapter)(nil)

// nolint: revive
func NewTickerCategoryAdapter() *ticketCategoryAdapter {
	return &ticketCategoryAdapter{}
}

func (t ticketCategoryAdapter) ToSchema(tc domain.TicketCategory) ticketCategorySchema {
	schema := ticketCategorySchema{
		ID:        tc.ID,
		EventID:   tc.EventID,
		Price:     tc.Price,
		Name:      tc.Name,
		CreatedAt: tc.CreatedAt,
		UpdatedAt: tc.UpdatedAt,
	}
	if tc.Description != nil {
		schema.Description = sql.NullString{
			String: *tc.Description,
			Valid:  true,
		}
	}
	return schema
}

func (t ticketCategoryAdapter) ToDomain(schema ticketCategorySchema) domain.TicketCategory {
	dom := domain.TicketCategory{
		ID:        schema.ID,
		EventID:   schema.EventID,
		Price:     schema.Price,
		Name:      schema.Name,
		CreatedAt: schema.CreatedAt,
		UpdatedAt: schema.UpdatedAt,
	}
	if schema.Description.Valid {
		dom.Description = &schema.Description.String
	}
	return dom
}
