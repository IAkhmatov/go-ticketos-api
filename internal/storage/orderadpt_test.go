// nolint: testpackage
package storage

import (
	"database/sql"
	"testing"

	"go-ticketos/internal/domain"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type orderAdapterTestSuite struct {
	suite.Suite
	a       *assert.Assertions
	adapter *orderAdapter
}

func TestOrderAdapterTestSuite(t *testing.T) {
	suite.Run(t, &orderAdapterTestSuite{})
}

func (s *orderAdapterTestSuite) SetupSuite() {
	s.a = assert.New(s.T())
	s.adapter = NewOrderAdapter()
}

func (s *orderAdapterTestSuite) TestToSchema() {
	type args struct {
		order domain.Order
	}
	tc := domain.NewTicketCategory(uuid.New(), 1000, "name", nil)
	dp := uint(50)
	p, err := domain.NewPromocode(1, nil, &dp, []domain.TicketCategory{tc})
	s.a.NoError(err)
	ticket, err := domain.NewTicket(tc, p)
	s.a.NoError(err)
	ord1, err := domain.NewOrder(
		"name",
		"test@tmail.co",
		"79999999999",
		[]domain.Ticket{*ticket},
	)
	s.a.NoError(err)
	ord2 := *ord1
	payment := domain.Payment{
		ID:  "test",
		URL: "https://test.test",
	}
	ord2.Payment = &payment
	err = ord2.UpdateStatus(domain.OrderStatusAwaitingPayment)
	s.a.NoError(err)

	tests := []struct {
		name string
		args args
		want orderSchema
	}{
		{
			name: "default",
			args: args{
				order: *ord1,
			},
			want: orderSchema{
				ID:   ord1.ID,
				Name: ord1.Name,
				Tickets: []ticketSchema{
					{
						ID:      ord1.Tickets[0].ID,
						OrderID: ord1.ID,
						TicketCategory: ticketCategorySchema{
							ID:        ord1.Tickets[0].TicketCategory.ID,
							EventID:   ord1.Tickets[0].TicketCategory.EventID,
							Price:     ord1.Tickets[0].TicketCategory.Price,
							Name:      ord1.Tickets[0].TicketCategory.Name,
							CreatedAt: ord1.Tickets[0].TicketCategory.CreatedAt,
							UpdatedAt: ord1.Tickets[0].TicketCategory.UpdatedAt,
						},
						PromocodeID: uuid.NullUUID{
							UUID:  *ord1.Tickets[0].PromocodeID,
							Valid: true,
						},
						FullPrice: ord1.Tickets[0].FullPrice,
						BuyPrice:  ord1.Tickets[0].BuyPrice,
						CreatedAt: ord1.Tickets[0].CreatedAt,
						UpdatedAt: ord1.Tickets[0].UpdatedAt,
					},
				},
				Email:     ord1.Email,
				Phone:     ord1.Phone,
				Status:    ord1.Status.String(),
				CreatedAt: ord1.CreatedAt,
				UpdatedAt: ord1.UpdatedAt,
			},
		},
		{
			name: "with payment",
			args: args{
				order: ord2,
			},
			want: orderSchema{
				ID:   ord2.ID,
				Name: ord2.Name,
				Tickets: []ticketSchema{
					{
						ID:      ord2.Tickets[0].ID,
						OrderID: ord2.ID,
						TicketCategory: ticketCategorySchema{
							ID:        ord1.Tickets[0].TicketCategory.ID,
							EventID:   ord1.Tickets[0].TicketCategory.EventID,
							Price:     ord1.Tickets[0].TicketCategory.Price,
							Name:      ord1.Tickets[0].TicketCategory.Name,
							CreatedAt: ord1.Tickets[0].TicketCategory.CreatedAt,
							UpdatedAt: ord1.Tickets[0].TicketCategory.UpdatedAt,
						},
						PromocodeID: uuid.NullUUID{
							UUID:  *ord2.Tickets[0].PromocodeID,
							Valid: true,
						},
						FullPrice: ord2.Tickets[0].FullPrice,
						BuyPrice:  ord2.Tickets[0].BuyPrice,
						CreatedAt: ord2.Tickets[0].CreatedAt,
						UpdatedAt: ord2.Tickets[0].UpdatedAt,
					},
				},
				Email:  ord2.Email,
				Phone:  ord2.Phone,
				Status: ord2.Status.String(),
				PaymentID: sql.NullString{
					String: ord2.Payment.ID,
					Valid:  true,
				},
				PaymentURL: sql.NullString{
					String: ord2.Payment.URL,
					Valid:  true,
				},
				CreatedAt: ord2.CreatedAt,
				UpdatedAt: ord2.UpdatedAt,
			},
		},
	}
	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			currentTest := tt
			actualDto := s.adapter.ToSchema(currentTest.args.order)
			s.a.Equal(currentTest.want.ID, actualDto.ID)
			s.a.Equal(currentTest.want.Name, actualDto.Name)
			s.a.Equal(currentTest.want.Email, actualDto.Email)
			s.a.Equal(currentTest.want.Phone, actualDto.Phone)
			s.a.Equal(currentTest.want.PaymentURL, actualDto.PaymentURL)
			s.a.Equal(currentTest.want.PaymentID, actualDto.PaymentID)
			s.a.Equal(currentTest.want.Status, actualDto.Status)
			s.a.Equal(currentTest.want.CreatedAt, actualDto.CreatedAt)
			s.a.Equal(currentTest.want.UpdatedAt, actualDto.UpdatedAt)
			s.a.Len(actualDto.Tickets, 1)
			s.a.Equal(currentTest.want.Tickets[0].ID, actualDto.Tickets[0].ID)
			s.a.Equal(currentTest.want.Tickets[0].OrderID, actualDto.Tickets[0].OrderID)
			s.a.Equal(currentTest.want.Tickets[0].TicketCategory.ID, actualDto.Tickets[0].TicketCategory.ID)
			s.a.Equal(currentTest.want.Tickets[0].FullPrice, actualDto.Tickets[0].FullPrice)
			s.a.Equal(currentTest.want.Tickets[0].BuyPrice, actualDto.Tickets[0].BuyPrice)
			s.a.Equal(currentTest.want.Tickets[0].CreatedAt, actualDto.Tickets[0].CreatedAt)
			s.a.Equal(currentTest.want.Tickets[0].UpdatedAt, actualDto.Tickets[0].UpdatedAt)
		})
	}
}

func (s *orderAdapterTestSuite) TestToDomain() {
	type args struct {
		order orderSchema
	}
	tc := domain.NewTicketCategory(uuid.New(), 1000, "name", nil)
	dp := uint(50)
	p, err := domain.NewPromocode(1, nil, &dp, []domain.TicketCategory{tc})
	s.a.NoError(err)
	ticket, err := domain.NewTicket(tc, p)
	s.a.NoError(err)
	ord1, err := domain.NewOrder(
		"name",
		"test@tmail.co",
		"79999999999",
		[]domain.Ticket{*ticket},
	)
	s.a.NoError(err)
	ord2 := *ord1
	payment := domain.Payment{
		ID:  "test",
		URL: "https://test.test",
	}
	ord2.Payment = &payment
	err = ord2.UpdateStatus(domain.OrderStatusAwaitingPayment)
	s.a.NoError(err)

	tests := []struct {
		name string
		args args
		want domain.Order
	}{
		{
			name: "default",
			args: args{
				order: orderSchema{
					ID:   ord1.ID,
					Name: ord1.Name,
					Tickets: []ticketSchema{
						{
							ID:      ord1.Tickets[0].ID,
							OrderID: ord1.ID,
							TicketCategory: ticketCategorySchema{
								ID:        ord1.Tickets[0].TicketCategory.ID,
								EventID:   ord1.Tickets[0].TicketCategory.EventID,
								Price:     ord1.Tickets[0].TicketCategory.Price,
								Name:      ord1.Tickets[0].TicketCategory.Name,
								CreatedAt: ord1.Tickets[0].TicketCategory.CreatedAt,
								UpdatedAt: ord1.Tickets[0].TicketCategory.UpdatedAt,
							},
							PromocodeID: uuid.NullUUID{
								UUID:  *ord1.Tickets[0].PromocodeID,
								Valid: true,
							},
							FullPrice: ord1.Tickets[0].FullPrice,
							BuyPrice:  ord1.Tickets[0].BuyPrice,
							CreatedAt: ord1.Tickets[0].CreatedAt,
							UpdatedAt: ord1.Tickets[0].UpdatedAt,
						},
					},
					Email:     ord1.Email,
					Phone:     ord1.Phone,
					Status:    ord1.Status.String(),
					CreatedAt: ord1.CreatedAt,
					UpdatedAt: ord1.UpdatedAt,
				},
			},
			want: *ord1,
		},
		{
			name: "with payment",
			args: args{
				order: orderSchema{
					ID:   ord2.ID,
					Name: ord2.Name,
					Tickets: []ticketSchema{
						{
							ID:      ord2.Tickets[0].ID,
							OrderID: ord2.ID,
							TicketCategory: ticketCategorySchema{
								ID:        ord2.Tickets[0].TicketCategory.ID,
								EventID:   ord2.Tickets[0].TicketCategory.EventID,
								Price:     ord2.Tickets[0].TicketCategory.Price,
								Name:      ord2.Tickets[0].TicketCategory.Name,
								CreatedAt: ord2.Tickets[0].TicketCategory.CreatedAt,
								UpdatedAt: ord2.Tickets[0].TicketCategory.UpdatedAt,
							},
							PromocodeID: uuid.NullUUID{
								UUID:  *ord2.Tickets[0].PromocodeID,
								Valid: true,
							},
							FullPrice: ord2.Tickets[0].FullPrice,
							BuyPrice:  ord2.Tickets[0].BuyPrice,
							CreatedAt: ord2.Tickets[0].CreatedAt,
							UpdatedAt: ord2.Tickets[0].UpdatedAt,
						},
					},
					Email:  ord2.Email,
					Phone:  ord2.Phone,
					Status: ord2.Status.String(),
					PaymentID: sql.NullString{
						String: ord2.Payment.ID,
						Valid:  true,
					},
					PaymentURL: sql.NullString{
						String: ord2.Payment.URL,
						Valid:  true,
					},
					CreatedAt: ord2.CreatedAt,
					UpdatedAt: ord2.UpdatedAt,
				},
			},
			want: ord2,
		},
	}
	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			currentTest := tt
			actualDto, err2 := s.adapter.ToDomain(currentTest.args.order)
			s.a.NoError(err2)
			s.a.Equal(currentTest.want.ID, actualDto.ID)
			s.a.Equal(currentTest.want.Name, actualDto.Name)
			s.a.Equal(currentTest.want.Email, actualDto.Email)
			s.a.Equal(currentTest.want.Phone, actualDto.Phone)
			if currentTest.want.Payment != nil {
				s.a.Equal(currentTest.want.Payment.ID, actualDto.Payment.ID)
				s.a.Equal(currentTest.want.Payment.URL, actualDto.Payment.URL)
			} else {
				s.a.Nil(actualDto.Payment)
			}
			s.a.Equal(currentTest.want.Status, actualDto.Status)
			s.a.Equal(currentTest.want.CreatedAt, actualDto.CreatedAt)
			s.a.Equal(currentTest.want.UpdatedAt, actualDto.UpdatedAt)
			s.a.Len(actualDto.Tickets, 1)
			s.a.Equal(currentTest.want.Tickets[0].ID, actualDto.Tickets[0].ID)
			s.a.Equal(currentTest.want.Tickets[0].TicketCategory.ID, actualDto.Tickets[0].TicketCategory.ID)
			s.a.Equal(currentTest.want.Tickets[0].FullPrice, actualDto.Tickets[0].FullPrice)
			s.a.Equal(currentTest.want.Tickets[0].BuyPrice, actualDto.Tickets[0].BuyPrice)
			s.a.Equal(currentTest.want.Tickets[0].CreatedAt, actualDto.Tickets[0].CreatedAt)
			s.a.Equal(currentTest.want.Tickets[0].UpdatedAt, actualDto.Tickets[0].UpdatedAt)
		})
	}
}
