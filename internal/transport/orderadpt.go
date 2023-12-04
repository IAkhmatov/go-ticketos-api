package transport

import (
	"go-ticketos/internal/domain"
)

type OrderAdapter interface {
	ToResponseDTO(order domain.Order) OrderResponseDTO
}

type orderAdapter struct {
	ticketCategoryAdapter TicketCategoryAdapter
}

// nolint: revive
func NewOrderAdapter() *orderAdapter {
	return &orderAdapter{
		ticketCategoryAdapter: NewTicketCategoryAdapter(),
	}
}

var _ OrderAdapter = (*orderAdapter)(nil)

func (o orderAdapter) ToResponseDTO(order domain.Order) OrderResponseDTO {
	dto := OrderResponseDTO{
		ID:        order.ID,
		Name:      order.Name,
		Email:     order.Email,
		Phone:     order.Phone,
		Tickets:   o.toTicketResponseDTOMany(order.Tickets),
		Status:    order.Status.String(),
		FullPrice: order.FullPrice(),
		BuyPrice:  order.BuyPrice(),
	}
	if order.Payment != nil {
		paymentResponseDTO := o.toPaymentResponseDTO(*order.Payment)
		dto.Payment = &paymentResponseDTO
	}
	return dto
}

func (o orderAdapter) toPaymentResponseDTO(payment domain.Payment) PaymentResponseDTO {
	return PaymentResponseDTO{
		ID:  payment.ID,
		URL: payment.URL,
	}
}

func (o orderAdapter) toTicketResponseDTO(ticket domain.Ticket) TicketResponseDTO {
	return TicketResponseDTO{
		ID:             ticket.ID,
		TicketCategory: o.ticketCategoryAdapter.ToResponseDTO(ticket.TicketCategory),
		PromocodeID:    ticket.PromocodeID,
		FullPrice:      ticket.FullPrice,
		BuyPrice:       ticket.BuyPrice,
	}
}

func (o orderAdapter) toTicketResponseDTOMany(tickets []domain.Ticket) []TicketResponseDTO {
	var dtos []TicketResponseDTO
	for _, t := range tickets {
		dto := o.toTicketResponseDTO(t)
		dtos = append(dtos, dto)
	}
	return dtos
}
