package transport

import (
	"go-ticketos/internal/domain"
)

type TicketCategoryAdapter interface {
	ToResponseDTO(tc domain.TicketCategory) TicketCategoryResponseDTO
}

type ticketCategoryAdapter struct{}

// nolint: revive
func NewTicketCategoryAdapter() *ticketCategoryAdapter {
	return &ticketCategoryAdapter{}
}

var _ TicketCategoryAdapter = (*ticketCategoryAdapter)(nil)

func (o ticketCategoryAdapter) ToResponseDTO(tc domain.TicketCategory) TicketCategoryResponseDTO {
	dto := TicketCategoryResponseDTO{
		ID:          tc.ID,
		EventID:     tc.EventID,
		Price:       tc.Price,
		Name:        tc.Name,
		Description: tc.Description,
	}
	return dto
}
