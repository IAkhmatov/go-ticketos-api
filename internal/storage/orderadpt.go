package storage

import (
	"database/sql"
	"fmt"

	"go-ticketos/internal/domain"

	"github.com/google/uuid"
)

type OrderAdapter interface {
	ToSchema(order domain.Order) orderSchema
	ToDomain(schema orderSchema) (*domain.Order, error)
}

type orderAdapter struct {
	ticketCategoryAdapter TicketCategoryAdapter
}

var _ OrderAdapter = (*orderAdapter)(nil)

// nolint: revive
func NewOrderAdapter() *orderAdapter {
	return &orderAdapter{
		ticketCategoryAdapter: NewTickerCategoryAdapter(),
	}
}

func (a orderAdapter) ToSchema(order domain.Order) orderSchema {
	schema := orderSchema{
		ID:        order.ID,
		Name:      order.Name,
		Email:     order.Email,
		Phone:     order.Phone,
		Status:    order.Status.String(),
		CreatedAt: order.CreatedAt,
		UpdatedAt: order.UpdatedAt,
	}
	if order.Payment != nil {
		schema.PaymentID = sql.NullString{
			String: order.Payment.ID,
			Valid:  true,
		}
		schema.PaymentURL = sql.NullString{
			String: order.Payment.URL,
			Valid:  true,
		}
	}
	for _, ticket := range order.Tickets {
		tSchema := ticketSchema{
			ID:               ticket.ID,
			OrderID:          order.ID,
			TicketCategoryID: ticket.TicketCategory.ID,
			TicketCategory:   a.ticketCategoryAdapter.ToSchema(ticket.TicketCategory),
			FullPrice:        ticket.FullPrice,
			BuyPrice:         ticket.BuyPrice,
			CreatedAt:        ticket.CreatedAt,
			UpdatedAt:        ticket.UpdatedAt,
		}
		if ticket.PromocodeID != nil {
			tSchema.PromocodeID = uuid.NullUUID{
				UUID:  *ticket.PromocodeID,
				Valid: true,
			}
		}
		schema.Tickets = append(schema.Tickets, tSchema)
	}
	return schema
}

func (a orderAdapter) ToDomain(schema orderSchema) (*domain.Order, error) {
	dom := domain.Order{
		ID:        schema.ID,
		Name:      schema.Name,
		Email:     schema.Email,
		Phone:     schema.Phone,
		Payment:   nil,
		Tickets:   nil,
		CreatedAt: schema.CreatedAt,
		UpdatedAt: schema.UpdatedAt,
	}
	status, err := domain.ParseOrderStatus(schema.Status)
	if err != nil {
		return nil, fmt.Errorf("ToDomain: can not parse status: %w", err)
	}
	dom.Status = status
	var payment domain.Payment
	if schema.PaymentID.Valid && schema.PaymentURL.Valid {
		payment.ID = schema.PaymentID.String
		payment.URL = schema.PaymentURL.String
		dom.Payment = &payment
	}
	for _, ticket := range schema.Tickets {
		ticketDom := domain.Ticket{
			ID:             ticket.ID,
			TicketCategory: a.ticketCategoryAdapter.ToDomain(ticket.TicketCategory),
			FullPrice:      ticket.FullPrice,
			BuyPrice:       ticket.BuyPrice,
			CreatedAt:      ticket.CreatedAt,
			UpdatedAt:      ticket.UpdatedAt,
		}
		if ticket.PromocodeID.Valid {
			promocodeID := ticket.PromocodeID.UUID
			ticketDom.PromocodeID = &promocodeID
		}
		dom.Tickets = append(dom.Tickets, ticketDom)
	}
	return &dom, nil
}
