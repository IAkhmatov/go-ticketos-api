package transport_test

import (
	"testing"

	"go-ticketos/internal/domain"
	"go-ticketos/internal/transport"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type orderAdapterTestSuite struct {
	suite.Suite
	a       *assert.Assertions
	adapter transport.OrderAdapter
}

func TestOrderAdapterTestSuite(t *testing.T) {
	suite.Run(t, &orderAdapterTestSuite{})
}

func (s *orderAdapterTestSuite) SetupSuite() {
	s.a = assert.New(s.T())
	s.adapter = transport.NewOrderAdapter()
}

func (s *orderAdapterTestSuite) TestToResponseDTO() {
	tc := domain.NewTicketCategory(uuid.New(), 10000, "test", nil)
	dp := uint(50)
	p, err := domain.NewPromocode(100, nil, &dp, []domain.TicketCategory{tc})
	s.a.NoError(err)
	t, err := domain.NewTicket(tc, p)
	s.a.NoError(err)
	o, err := domain.NewOrder(
		"test",
		"ggr@ggr.crr",
		"79999999999",
		[]domain.Ticket{*t},
	)
	s.a.NoError(err)

	orderDTO := s.adapter.ToResponseDTO(*o)

	s.a.Equal(o.ID, orderDTO.ID)
	s.a.Equal(o.Name, orderDTO.Name)
	s.a.Equal(o.Email, orderDTO.Email)
	s.a.Equal(o.Phone, orderDTO.Phone)
	s.a.Nil(orderDTO.Payment)
	s.a.Equal(o.Status.String(), orderDTO.Status)
	s.a.Equal(o.FullPrice(), orderDTO.FullPrice)
	s.a.Equal(o.BuyPrice(), orderDTO.BuyPrice)
	s.a.Equal(o.Tickets[0].ID, orderDTO.Tickets[0].ID)
	s.a.Equal(o.Tickets[0].TicketCategory.ID, orderDTO.Tickets[0].TicketCategory.ID)
	s.a.Equal(o.Tickets[0].PromocodeID, orderDTO.Tickets[0].PromocodeID)
	s.a.Equal(o.Tickets[0].FullPrice, orderDTO.Tickets[0].FullPrice)
	s.a.Equal(o.Tickets[0].BuyPrice, orderDTO.Tickets[0].BuyPrice)

	// test payment
	err = o.UpdateStatus(domain.OrderStatusAwaitingPayment)
	s.a.NoError(err)
	payment := domain.NewPayment("test", "https://test.test")
	o.UpdatePayment(payment)

	orderDTO = s.adapter.ToResponseDTO(*o)

	s.a.Equal(payment.ID, orderDTO.Payment.ID)
	s.a.Equal(payment.URL, orderDTO.Payment.URL)
}
