//go:generate go-enum --marshal --mustparse

package domain

import (
	"errors"
	"slices"
	"time"

	"github.com/google/uuid"
)

// ENUM(prepared, awaiting_payment, completed).
type OrderStatus string

type (
	Ticket struct {
		ID             uuid.UUID
		TicketCategory TicketCategory
		PromocodeID    *uuid.UUID
		FullPrice      uint
		BuyPrice       uint
		CreatedAt      time.Time
		UpdatedAt      time.Time
	}

	Payment struct {
		ID  string
		URL string
	}

	Order struct {
		ID        uuid.UUID
		Name      string
		Email     string
		Phone     string
		Payment   *Payment
		Tickets   []Ticket
		Status    OrderStatus
		CreatedAt time.Time
		UpdatedAt time.Time
	}

	CreateOrderProps struct {
		Name              string
		Email             string
		Phone             string
		PromocodeID       *uuid.UUID
		TicketCategoryIDs []uuid.UUID
	}

	// nolint: lll
	UpdateOrderProps struct {
		OrderID    uuid.UUID   `validate:"required"`
		Status     OrderStatus `validate:"required,oneof=awaiting_payment completed"`
		PaymentID  *string     `validate:"omitempty,required_if=Status awaiting_payment,excluded_unless=Status awaiting_payment"`
		PaymentURL *string     `validate:"omitempty,required_if=Status awaiting_payment,excluded_unless=Status awaiting_payment,url"`
	}

	OrderRepo interface {
		Create(order Order) error
		Update(order Order) error
		GetByID(id uuid.UUID) (*Order, error)
	}

	OrderService interface {
		Create(props CreateOrderProps) (*Order, error)
		Update(props UpdateOrderProps) (*Order, error)
		GetByID(id uuid.UUID) (*Order, error)
	}
)

func NewPayment(id, url string) Payment {
	return Payment{
		ID:  id,
		URL: url,
	}
}

// NewTicket create new ticket.
// If promocode can not append to this ticket return an error.
func NewTicket(tc TicketCategory, promocode *Promocode) (*Ticket, error) {
	now := time.Now().UTC()
	ticket := Ticket{
		ID:             uuid.New(),
		TicketCategory: tc,
		FullPrice:      tc.Price,
		BuyPrice:       tc.Price,
		CreatedAt:      now,
		UpdatedAt:      now,
	}
	if promocode != nil {
		if !tc.CanUsePromocode(*promocode) {
			return nil, errors.New("NewTicket: can not use this promocode with this ticket category")
		}
		promocodeCopy := *promocode
		ticket.PromocodeID = &promocodeCopy.ID
		ticket.BuyPrice = promocode.CalcBuyPrice(ticket.FullPrice)
	}
	return &ticket, nil
}

func NewOrder(
	name string,
	email string,
	phone string,
	tickets []Ticket,
) (*Order, error) {
	if len(tickets) == 0 {
		return nil, errors.New("NewOrder: len(tickets) = 0")
	}
	ticketsCopy := make([]Ticket, len(tickets))
	copy(ticketsCopy, tickets)
	orderID := uuid.New()
	now := time.Now().UTC()
	order := Order{
		ID:        orderID,
		Name:      name,
		Email:     email,
		Phone:     phone,
		Tickets:   tickets,
		Status:    OrderStatusPrepared,
		CreatedAt: now,
		UpdatedAt: now,
	}
	return &order, nil
}

func (o Order) BuyPrice() uint {
	p := uint(0)
	for _, ticket := range o.Tickets {
		p += ticket.BuyPrice
	}
	return p
}

func (o Order) FullPrice() uint {
	p := uint(0)
	for _, ticket := range o.Tickets {
		p += ticket.FullPrice
	}
	return p
}

func (o *Order) UpdatePayment(payment Payment) {
	o.Payment = &payment
	o.UpdatedAt = time.Now().UTC()
}

func (o *Order) UpdateStatus(newStatus OrderStatus) error {
	statusMap := map[OrderStatus][]OrderStatus{
		OrderStatusPrepared:        {OrderStatusAwaitingPayment},
		OrderStatusAwaitingPayment: {OrderStatusCompleted},
		OrderStatusCompleted:       nil,
	}
	if val, ok := statusMap[o.Status]; ok {
		if slices.Contains(val, newStatus) {
			o.Status = newStatus
			o.UpdatedAt = time.Now().UTC()
			return nil
		}
	}
	return errors.New("can not change status")
}
