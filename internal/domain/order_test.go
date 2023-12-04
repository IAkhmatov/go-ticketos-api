package domain_test

import (
	"testing"
	"time"

	"go-ticketos/internal/domain"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestNewOrder(t *testing.T) {
	a := assert.New(t)
	type args struct {
		name    string
		email   string
		phone   string
		tickets []domain.Ticket
	}
	tc := domain.NewTicketCategory(uuid.New(), 1000, "name", nil)
	dp := uint(50)
	p, err := domain.NewPromocode(1, nil, &dp, []domain.TicketCategory{tc})
	a.NoError(err)
	t1, err := domain.NewTicket(tc, p)
	a.NoError(err)
	t2, err := domain.NewTicket(tc, p)
	a.NoError(err)
	tests := []struct {
		name        string
		args        args
		isErrResult bool
	}{
		{
			name: "without tickets",
			args: args{
				name:    "name",
				email:   "email",
				phone:   "79999999999",
				tickets: nil,
			},
			isErrResult: true,
		},
		{
			name: "good case",
			args: args{
				name:  "name",
				email: "email",
				phone: "79999999999",
				tickets: []domain.Ticket{
					*t1,
					*t2,
				},
			},
			isErrResult: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			o, err1 := domain.NewOrder(
				tt.args.name,
				tt.args.email,
				tt.args.phone,
				tt.args.tickets,
			)
			if tt.isErrResult {
				a.Error(err1)
			} else {
				a.NoError(err1)
				a.NotEqual(uuid.Nil, o.ID)
				a.Equal(tt.args.name, o.Name)
				a.Equal(tt.args.email, o.Email)
				a.Equal(domain.OrderStatusPrepared, o.Status)
				a.False(o.CreatedAt.IsZero())
				a.False(o.UpdatedAt.IsZero())
				a.WithinDuration(o.CreatedAt, o.UpdatedAt, time.Millisecond)
			}
		})
	}
}

func TestOrder_BuyPrice(t *testing.T) {
	a := assert.New(t)
	o := createOrder(t)

	a.Equal(uint(1000), o.BuyPrice())
}

func TestOrder_FullPrice(t *testing.T) {
	a := assert.New(t)
	o := createOrder(t)

	a.Equal(uint(2000), o.FullPrice())
}

func TestNewTicket(t *testing.T) {
	a := assert.New(t)
	type args struct {
		ticketCategory domain.TicketCategory
		promocode      *domain.Promocode
		isErrResult    bool
	}
	tc := domain.NewTicketCategory(uuid.New(), 1000, "name", nil)
	dp := uint(50)
	p, err := domain.NewPromocode(1, nil, &dp, []domain.TicketCategory{tc})
	a.NoError(err)
	tc2 := domain.NewTicketCategory(uuid.New(), 3000, "name", nil)
	p2, err := domain.NewPromocode(1, nil, &dp, []domain.TicketCategory{tc2})
	a.NoError(err)
	var tests = []struct {
		name string
		args args
	}{
		{
			name: "without promocode",
			args: args{
				ticketCategory: tc,
				promocode:      nil,
				isErrResult:    false,
			},
		},
		{
			name: "with promocode",
			args: args{
				ticketCategory: tc,
				promocode:      p,
				isErrResult:    false,
			},
		},
		{
			name: "with incorrect promocode",
			args: args{
				ticketCategory: tc,
				promocode:      p2,
				isErrResult:    true,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ticket, err2 := domain.NewTicket(
				tt.args.ticketCategory,
				tt.args.promocode,
			)
			// nolint: nestif
			if tt.args.isErrResult {
				a.Error(err2)
				a.Nil(ticket)
			} else {
				a.NoError(err2)
				a.NotEqual(uuid.Nil, ticket.ID)
				a.Equal(tt.args.ticketCategory.ID, ticket.TicketCategory.ID)
				a.Equal(ticket.TicketCategory.Price, ticket.FullPrice)
				if tt.args.promocode != nil {
					a.Equal(tt.args.promocode.ID, *ticket.PromocodeID)
				} else {
					a.Nil(ticket.PromocodeID)
				}
				if tt.args.promocode != nil {
					a.Equal(p.CalcBuyPrice(tc.Price), ticket.BuyPrice)
				} else {
					a.Equal(ticket.FullPrice, ticket.BuyPrice)
				}
				a.False(ticket.CreatedAt.IsZero())
				a.False(ticket.UpdatedAt.IsZero())
				a.WithinDuration(ticket.CreatedAt, ticket.UpdatedAt, time.Millisecond)
			}
		})
	}
}

func TestUpdateStatus(t *testing.T) {
	a := assert.New(t)
	o := createOrder(t)

	// test OrderStatusPrepared
	err := o.UpdateStatus(domain.OrderStatusCompleted)
	a.Error(err)
	a.Equal(domain.OrderStatusPrepared, o.Status)

	err = o.UpdateStatus(domain.OrderStatusAwaitingPayment)
	a.NoError(err)
	a.Equal(domain.OrderStatusAwaitingPayment, o.Status)

	// test OrderStatusAwaitingPayment
	err = o.UpdateStatus(domain.OrderStatusPrepared)
	a.Error(err)
	a.Equal(domain.OrderStatusAwaitingPayment, o.Status)

	err = o.UpdateStatus(domain.OrderStatusCompleted)
	a.NoError(err)
	a.Equal(domain.OrderStatusCompleted, o.Status)

	// test OrderStatusCompleted
	err = o.UpdateStatus(domain.OrderStatusPrepared)
	a.Error(err)
	a.Equal(domain.OrderStatusCompleted, o.Status)

	err = o.UpdateStatus(domain.OrderStatusAwaitingPayment)
	a.Error(err)
	a.Equal(domain.OrderStatusCompleted, o.Status)
}

func createOrder(t *testing.T) domain.Order {
	t.Helper()
	a := assert.New(t)
	tc := domain.NewTicketCategory(uuid.New(), 1000, "name", nil)
	dp := uint(50)
	p, err := domain.NewPromocode(1, nil, &dp, []domain.TicketCategory{tc})
	a.NoError(err)
	t1, err := domain.NewTicket(tc, p)
	a.NoError(err)
	t2, err := domain.NewTicket(tc, p)
	a.NoError(err)
	o, err := domain.NewOrder(
		"name",
		"email",
		"79999999999",
		[]domain.Ticket{
			*t1,
			*t2,
		},
	)
	a.NoError(err)
	return *o
}
