package service_test

import (
	"errors"
	"testing"

	"go-ticketos/internal/config"
	"go-ticketos/internal/domain"
	"go-ticketos/internal/service"
	"go-ticketos/pkg/log"
	"go-ticketos/pkg/paymasterclient"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type orderServiceTestSuite struct {
	suite.Suite
	a                  *assert.Assertions
	cfg                *config.Config
	orderRepo          *domain.MockOrderRepo
	ticketCategoryRepo *domain.MockTicketCategoryRepo
	promocodeRepo      *domain.MockPromocodeRepo
	payMasterClient    *paymasterclient.MockPayMasterClient
	svc                domain.OrderService
}

func TestOrderServiceTestSuite(t *testing.T) {
	suite.Run(t, &orderServiceTestSuite{})
}

func (s *orderServiceTestSuite) SetupTest() {
	s.a = assert.New(s.T())
	cfg, err := config.NewConfig()
	s.a.NoError(err)
	s.cfg = cfg
	s.orderRepo = domain.NewMockOrderRepo(s.T())
	s.ticketCategoryRepo = domain.NewMockTicketCategoryRepo(s.T())
	s.promocodeRepo = domain.NewMockPromocodeRepo(s.T())
	s.payMasterClient = paymasterclient.NewMockPayMasterClient(s.T())
	svc, err := service.NewOrderService(
		s.orderRepo,
		s.ticketCategoryRepo,
		s.promocodeRepo,
		s.payMasterClient,
		cfg,
		log.NewLogger(cfg),
	)
	s.a.NoError(err)
	s.svc = svc
}

func (s *orderServiceTestSuite) TeardownTest() {
	s.orderRepo.AssertExpectations(s.T())
	s.ticketCategoryRepo.AssertExpectations(s.T())
	s.promocodeRepo.AssertExpectations(s.T())
}

func (s *orderServiceTestSuite) TestCreate_GetPromocodeError() {
	props, _, p, _ := s.createEntities()
	s.promocodeRepo.On("GetByID", p.ID).Once().Return(nil, errors.New("test"))

	o, err := s.svc.Create(props)

	s.a.Error(err)
	s.a.Nil(o)
}

func (s *orderServiceTestSuite) TestCreate_GetTicketCategoryError() {
	props, tc, p, _ := s.createEntities()
	s.promocodeRepo.On("GetByID", p.ID).Once().Return(&p, nil)
	s.ticketCategoryRepo.On("GetByID", tc.ID).Once().Return(nil, errors.New("test"))

	o, err := s.svc.Create(props)

	s.a.Error(err)
	s.a.Nil(o)
}

func (s *orderServiceTestSuite) TestCreate_CreateOrderError() {
	props, tc, p, intOrder := s.createEntities()
	s.promocodeRepo.On("GetByID", p.ID).Once().Return(&p, nil)
	s.ticketCategoryRepo.On("GetByID", tc.ID).Once().Return(&tc, nil)
	s.orderRepo.On(
		"Create",
		mock.MatchedBy(
			func(o domain.Order) bool {
				return o.Name == intOrder.Name &&
					o.Email == intOrder.Email &&
					o.Status == intOrder.Status &&
					o.Tickets[0].TicketCategory.ID == intOrder.Tickets[0].TicketCategory.ID &&
					*o.Tickets[0].PromocodeID == *intOrder.Tickets[0].PromocodeID
			},
		),
	).Once().Return(errors.New("test"))

	o, err := s.svc.Create(props)

	s.a.Error(err)
	s.a.Nil(o)
}

func (s *orderServiceTestSuite) TestCreate_CreatePaymentError() {
	props, tc, p, intOrder := s.createEntities()
	s.promocodeRepo.On("GetByID", p.ID).Once().Return(&p, nil)
	s.ticketCategoryRepo.On("GetByID", tc.ID).Once().Return(&tc, nil)
	s.orderRepo.On(
		"Create",
		mock.MatchedBy(
			func(o domain.Order) bool {
				return o.Name == intOrder.Name &&
					o.Email == intOrder.Email &&
					o.Status == intOrder.Status &&
					o.Tickets[0].TicketCategory.ID == intOrder.Tickets[0].TicketCategory.ID &&
					*o.Tickets[0].PromocodeID == *intOrder.Tickets[0].PromocodeID
			},
		),
	).Once().Return(nil)
	s.payMasterClient.On(
		"CreateInvoice",
		mock.MatchedBy(
			func(dto paymasterclient.CreateInvoiceRequestDTO) bool {
				return dto.MerchantID == s.cfg.PayMasterMerchantID &&
					dto.PaymentMethod == paymasterclient.PaymentMethodBankCard &&
					dto.Amount.Value == float32(intOrder.BuyPrice())
			},
		),
	).Once().Return(nil, errors.New("test"))

	order, err := s.svc.Create(props)

	s.a.Error(err)
	s.a.Nil(order)
}

func (s *orderServiceTestSuite) TestCreate_UpdateOrderError() {
	props, tc, p, intOrder := s.createEntities()
	s.promocodeRepo.On("GetByID", p.ID).Once().Return(&p, nil)
	s.ticketCategoryRepo.On("GetByID", tc.ID).Once().Return(&tc, nil)
	s.orderRepo.On(
		"Create",
		mock.MatchedBy(
			func(o domain.Order) bool {
				return o.Name == intOrder.Name &&
					o.Email == intOrder.Email &&
					o.Status == intOrder.Status &&
					o.Tickets[0].TicketCategory.ID == intOrder.Tickets[0].TicketCategory.ID &&
					*o.Tickets[0].PromocodeID == *intOrder.Tickets[0].PromocodeID
			},
		),
	).Once().Return(nil)
	s.payMasterClient.On(
		"CreateInvoice",
		mock.MatchedBy(
			func(dto paymasterclient.CreateInvoiceRequestDTO) bool {
				return dto.MerchantID == s.cfg.PayMasterMerchantID &&
					dto.PaymentMethod == paymasterclient.PaymentMethodBankCard &&
					dto.Amount.Value == float32(intOrder.BuyPrice())
			},
		),
	).Once().Return(&paymasterclient.CreateInvoiceResponseDTO{
		PaymentID: "1",
		URL:       "https://test.com",
	},
		nil,
	)
	s.orderRepo.On(
		"GetByID",
		mock.AnythingOfType("uuid.UUID"),
	).Once().Return(nil, errors.New("test"))

	o, err := s.svc.Create(props)

	s.a.Error(err)
	s.a.Nil(o)
}

func (s *orderServiceTestSuite) TestCreate_GoodCase() {
	props, tc, p, intOrder := s.createEntities()
	s.promocodeRepo.On("GetByID", p.ID).Once().Return(&p, nil)
	s.ticketCategoryRepo.On("GetByID", tc.ID).Once().Return(&tc, nil)
	s.payMasterClient.On(
		"CreateInvoice",
		mock.MatchedBy(
			func(dto paymasterclient.CreateInvoiceRequestDTO) bool {
				return dto.MerchantID == s.cfg.PayMasterMerchantID &&
					dto.PaymentMethod == paymasterclient.PaymentMethodBankCard &&
					dto.Amount.Value == float32(intOrder.BuyPrice())
			},
		),
	).Once().Return(&paymasterclient.CreateInvoiceResponseDTO{
		PaymentID: "1",
		URL:       "https://test.com",
	},
		nil,
	)
	s.orderRepo.On(
		"Create",
		mock.MatchedBy(
			func(o domain.Order) bool {
				return o.Name == intOrder.Name &&
					o.Email == intOrder.Email &&
					o.Tickets[0].TicketCategory.ID == intOrder.Tickets[0].TicketCategory.ID &&
					*o.Tickets[0].PromocodeID == *intOrder.Tickets[0].PromocodeID
			},
		),
	).Once().Return(nil)
	s.orderRepo.On(
		"GetByID",
		mock.AnythingOfType("uuid.UUID"),
	).Once().Return(&intOrder, nil)
	s.orderRepo.On(
		"Update",
		mock.MatchedBy(
			func(o domain.Order) bool {
				return o.Name == intOrder.Name &&
					o.Email == intOrder.Email &&
					o.Status == domain.OrderStatusAwaitingPayment &&
					o.Tickets[0].TicketCategory.ID == intOrder.Tickets[0].TicketCategory.ID &&
					*o.Tickets[0].PromocodeID == *intOrder.Tickets[0].PromocodeID
			},
		),
	).Once().Return(nil)

	o, err := s.svc.Create(props)

	s.a.NoError(err)
	s.a.Equal(intOrder.Name, o.Name)
	s.a.Equal(intOrder.Email, o.Email)
	s.a.Equal(intOrder.Tickets[0].TicketCategory.ID, o.Tickets[0].TicketCategory.ID)
	s.a.Equal(intOrder.Tickets[0].PromocodeID, o.Tickets[0].PromocodeID)
}

func (s *orderServiceTestSuite) TestUpdate_ValidationError() {
	_, _, _, intOrder := s.createEntities()
	strPointer := func(s string) *string {
		return &s
	}
	props := domain.UpdateOrderProps{
		OrderID:    intOrder.ID,
		Status:     domain.OrderStatusCompleted,
		PaymentID:  strPointer("test"),
		PaymentURL: strPointer("https://test.com"),
	}

	dbOrder, err := s.svc.Update(props)

	s.a.Error(err)
	s.a.Nil(dbOrder)
}

func (s *orderServiceTestSuite) TestUpdate_GettingOrderError() {
	_, _, _, intOrder := s.createEntities()
	strPointer := func(s string) *string {
		return &s
	}
	props := domain.UpdateOrderProps{
		OrderID:    intOrder.ID,
		Status:     domain.OrderStatusAwaitingPayment,
		PaymentID:  strPointer("test"),
		PaymentURL: strPointer("https://test.com"),
	}
	s.orderRepo.On("GetByID", props.OrderID).Once().Return(nil, errors.New("test"))

	dbOrder, err := s.svc.Update(props)

	s.a.Error(err)
	s.a.Nil(dbOrder)
}

func (s *orderServiceTestSuite) TestUpdate_UpdatingStatusFromCompletedToAwaitingError() {
	_, _, _, intOrder := s.createEntities()
	strPointer := func(s string) *string {
		return &s
	}
	props := domain.UpdateOrderProps{
		OrderID:    intOrder.ID,
		Status:     domain.OrderStatusAwaitingPayment,
		PaymentID:  strPointer("test"),
		PaymentURL: strPointer("https://test.com"),
	}
	intOrder.Status = domain.OrderStatusCompleted
	s.orderRepo.On("GetByID", props.OrderID).Once().Return(&intOrder, nil)

	dbOrder, err := s.svc.Update(props)

	s.a.Error(err)
	s.a.Nil(dbOrder)
}

func (s *orderServiceTestSuite) TestUpdate_UpdatingStatusFromPreparedToCompletedError() {
	_, _, _, intOrder := s.createEntities()
	props := domain.UpdateOrderProps{
		OrderID: intOrder.ID,
		Status:  domain.OrderStatusCompleted,
	}
	s.orderRepo.On("GetByID", props.OrderID).Once().Return(&intOrder, nil)

	dbOrder, err := s.svc.Update(props)

	s.a.Error(err)
	s.a.Nil(dbOrder)
}

func (s *orderServiceTestSuite) TestUpdate_UpdatingInDBError() {
	_, _, _, intOrder := s.createEntities()
	strPointer := func(s string) *string {
		return &s
	}
	props := domain.UpdateOrderProps{
		OrderID:    intOrder.ID,
		Status:     domain.OrderStatusAwaitingPayment,
		PaymentID:  strPointer("test"),
		PaymentURL: strPointer("https://test.com"),
	}
	s.orderRepo.On("GetByID", props.OrderID).Once().Return(&intOrder, nil)
	s.orderRepo.On(
		"Update",
		mock.MatchedBy(
			func(o domain.Order) bool {
				return o.Name == intOrder.Name &&
					o.Email == intOrder.Email &&
					o.Status == domain.OrderStatusAwaitingPayment &&
					o.Tickets[0].TicketCategory.ID == intOrder.Tickets[0].TicketCategory.ID &&
					*o.Tickets[0].PromocodeID == *intOrder.Tickets[0].PromocodeID
			},
		)).Once().Return(errors.New("test"))
	dbOrder, err := s.svc.Update(props)

	s.a.Error(err)
	s.a.Nil(dbOrder)
}

func (s *orderServiceTestSuite) TestUpdate_ToAwaitingPaymentGoodCase() {
	_, _, _, intOrder := s.createEntities()
	strPointer := func(s string) *string {
		return &s
	}
	props := domain.UpdateOrderProps{
		OrderID:    intOrder.ID,
		Status:     domain.OrderStatusAwaitingPayment,
		PaymentID:  strPointer("test"),
		PaymentURL: strPointer("https://test.com"),
	}
	s.orderRepo.On("GetByID", props.OrderID).Once().Return(&intOrder, nil)
	s.orderRepo.On(
		"Update",
		mock.MatchedBy(
			func(o domain.Order) bool {
				return o.Name == intOrder.Name &&
					o.Email == intOrder.Email &&
					o.Status == domain.OrderStatusAwaitingPayment &&
					o.Tickets[0].TicketCategory.ID == intOrder.Tickets[0].TicketCategory.ID &&
					*o.Tickets[0].PromocodeID == *intOrder.Tickets[0].PromocodeID
			},
		)).Once().Return(nil)
	dbOrder, err := s.svc.Update(props)

	s.a.NoError(err)
	s.a.Equal(intOrder.ID, dbOrder.ID)
}

func (s *orderServiceTestSuite) TestUpdate_ToCompletedGoodCase() {
	_, _, _, intOrder := s.createEntities()
	props := domain.UpdateOrderProps{
		OrderID: intOrder.ID,
		Status:  domain.OrderStatusCompleted,
	}
	err := intOrder.UpdateStatus(domain.OrderStatusAwaitingPayment)
	s.a.NoError(err)
	payment := domain.Payment{
		ID:  "test",
		URL: "https://test.test",
	}
	intOrder.UpdatePayment(payment)
	s.orderRepo.On("GetByID", props.OrderID).Once().Return(&intOrder, nil)
	s.orderRepo.On(
		"Update",
		mock.MatchedBy(
			func(o domain.Order) bool {
				return o.Name == intOrder.Name &&
					o.Email == intOrder.Email &&
					o.Status == domain.OrderStatusCompleted &&
					o.Tickets[0].TicketCategory.ID == intOrder.Tickets[0].TicketCategory.ID &&
					*o.Tickets[0].PromocodeID == *intOrder.Tickets[0].PromocodeID
			},
		)).Once().Return(nil)

	dbOrder, err := s.svc.Update(props)

	s.a.NoError(err)
	s.a.Equal(intOrder.ID, dbOrder.ID)
}

func (s *orderServiceTestSuite) TestGetByID_GettingOrderError() {
	_, _, _, intOrder := s.createEntities()
	s.orderRepo.On("GetByID", intOrder.ID).Once().Return(nil, errors.New("test"))

	dbOrder, err := s.svc.GetByID(intOrder.ID)

	s.a.Error(err)
	s.a.Nil(dbOrder)
}

func (s *orderServiceTestSuite) TestGetByID_GoodCase() {
	_, _, _, intOrder := s.createEntities()
	s.orderRepo.On("GetByID", intOrder.ID).Once().Return(&intOrder, nil)

	dbOrder, err := s.svc.GetByID(intOrder.ID)

	s.a.NoError(err)
	s.a.Equal(intOrder.ID, dbOrder.ID)
}

func (s *orderServiceTestSuite) createEntities() (
	domain.CreateOrderProps,
	domain.TicketCategory,
	domain.Promocode,
	domain.Order,
) {
	s.T().Helper()
	tc := domain.NewTicketCategory(uuid.New(), 1000, "test1", nil)
	dv := uint(500)
	p, err := domain.NewPromocode(1, &dv, nil, []domain.TicketCategory{tc})
	s.a.NoError(err)
	props := domain.CreateOrderProps{
		Name:              "Test name",
		Email:             "test@email.com",
		Phone:             "779999999999",
		PromocodeID:       &p.ID,
		TicketCategoryIDs: []uuid.UUID{tc.ID},
	}
	ticket, err := domain.NewTicket(tc, p)
	s.a.NoError(err)
	o, err := domain.NewOrder(
		props.Name,
		props.Email,
		props.Phone,
		[]domain.Ticket{*ticket},
	)
	s.a.NoError(err)
	return props, tc, *p, *o
}
