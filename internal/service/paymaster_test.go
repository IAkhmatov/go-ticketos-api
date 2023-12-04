package service_test

import (
	"errors"
	"testing"
	"time"

	"go-ticketos/internal/config"
	"go-ticketos/internal/domain"
	"go-ticketos/internal/service"
	"go-ticketos/pkg/log"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type payMasterWebHookUseCaseTestSuite struct {
	suite.Suite
	a            *assert.Assertions
	orderService *domain.MockOrderService
	uc           domain.PayMasterWebHookUseCase
}

func TestPayMasterWebHookUseCaseTestSuite(t *testing.T) {
	suite.Run(t, &payMasterWebHookUseCaseTestSuite{})
}

func (s *payMasterWebHookUseCaseTestSuite) SetupTest() {
	s.a = assert.New(s.T())
	s.orderService = domain.NewMockOrderService(s.T())
	cfg, err := config.NewConfig()
	s.a.NoError(err)
	uc, err := service.NewPayMasterWebHookUseCase(s.orderService, log.NewLogger(cfg))
	s.a.NoError(err)
	s.uc = uc
}

func (s *payMasterWebHookUseCaseTestSuite) TeardownTest() {
	s.orderService.AssertExpectations(s.T())
}

func (s *payMasterWebHookUseCaseTestSuite) TestExecute_ParsingOrderIDError() {
	props := domain.PayMasterWebHookUseCaseProps{
		ID:         "123213",
		Created:    time.Now().UTC(),
		TestMode:   true,
		Status:     "Settled",
		MerchantID: "54a1a08f-b344-44b1-9f06-fe5868170e4c",
		Amount: domain.Amount{
			Value:    1000,
			Currency: "RUB",
		},
		Invoice: domain.Invoice{
			Description: "Test",
			OrderNo:     "test",
		},
		PaymentData: domain.PaymentData{
			PaymentMethod:          "BankCard",
			PaymentInstrumentTitle: "410000XXXXXX0001",
		},
	}

	err := s.uc.Execute(props)

	s.a.Error(err)
}

func (s *payMasterWebHookUseCaseTestSuite) TestExecute_GettingOrderError() {
	props := domain.PayMasterWebHookUseCaseProps{
		ID:         "123213",
		Created:    time.Now().UTC(),
		TestMode:   true,
		Status:     "Settled",
		MerchantID: "54a1a08f-b344-44b1-9f06-fe5868170e4c",
		Amount: domain.Amount{
			Value:    1000,
			Currency: "RUB",
		},
		Invoice: domain.Invoice{
			Description: "Test",
			OrderNo:     "de22bf4d-882d-4d52-bc52-7b691cf4b1c6",
		},
		PaymentData: domain.PaymentData{
			PaymentMethod:          "BankCard",
			PaymentInstrumentTitle: "410000XXXXXX0001",
		},
	}
	s.orderService.On("GetByID", uuid.MustParse(props.Invoice.OrderNo)).Once().Return(nil, errors.New("test"))
	err := s.uc.Execute(props)

	s.a.Error(err)
}

func (s *payMasterWebHookUseCaseTestSuite) TestExecute_IncorrectAmountError() {
	order := s.createOrder()
	err := order.UpdateStatus(domain.OrderStatusAwaitingPayment)
	s.a.NoError(err)
	order.UpdatePayment(domain.Payment{
		ID:  "test_id",
		URL: "https://test.test",
	})
	props := domain.PayMasterWebHookUseCaseProps{
		ID:         order.Payment.ID,
		Created:    time.Now().UTC(),
		TestMode:   true,
		Status:     "Settled",
		MerchantID: "54a1a08f-b344-44b1-9f22-fe5262170e4c",
		Amount: domain.Amount{
			Value:    float32(order.BuyPrice()*100 + 5),
			Currency: "RUB",
		},
		Invoice: domain.Invoice{
			Description: "Test",
			OrderNo:     order.ID.String(),
		},
		PaymentData: domain.PaymentData{
			PaymentMethod:          "BankCard",
			PaymentInstrumentTitle: "410000XXXXXX0001",
		},
	}
	s.orderService.On("GetByID", uuid.MustParse(props.Invoice.OrderNo)).Once().Return(&order, nil)

	err = s.uc.Execute(props)

	s.a.Error(err)
}

func (s *payMasterWebHookUseCaseTestSuite) TestExecute_IncorrectPaymentIDError() {
	order := s.createOrder()
	err := order.UpdateStatus(domain.OrderStatusAwaitingPayment)
	s.a.NoError(err)
	order.UpdatePayment(domain.Payment{
		ID:  "test_id",
		URL: "https://test.test",
	})
	props := domain.PayMasterWebHookUseCaseProps{
		ID:         "newId",
		Created:    time.Now().UTC(),
		TestMode:   true,
		Status:     "Settled",
		MerchantID: "54a1a08f-b344-44b1-9f06-fe5868170e4c",
		Amount: domain.Amount{
			Value:    float32(order.BuyPrice() * 100),
			Currency: "RUB",
		},
		Invoice: domain.Invoice{
			Description: "Test",
			OrderNo:     order.ID.String(),
		},
		PaymentData: domain.PaymentData{
			PaymentMethod:          "BankCard",
			PaymentInstrumentTitle: "410000XXXXXX0001",
		},
	}
	s.orderService.On("GetByID", uuid.MustParse(props.Invoice.OrderNo)).Once().Return(&order, nil)

	err = s.uc.Execute(props)

	s.a.Error(err)
}

func (s *payMasterWebHookUseCaseTestSuite) TestExecute_IncorrectPaymentStatusError() {
	order := s.createOrder()
	err := order.UpdateStatus(domain.OrderStatusAwaitingPayment)
	s.a.NoError(err)
	order.UpdatePayment(domain.Payment{
		ID:  "test_id",
		URL: "https://test.test",
	})
	props := domain.PayMasterWebHookUseCaseProps{
		ID:         order.Payment.ID,
		Created:    time.Now().UTC(),
		TestMode:   true,
		Status:     "Some strange status",
		MerchantID: "54a1a08f-b344-44b1-9f06-fe5868170e4c",
		Amount: domain.Amount{
			Value:    float32(order.BuyPrice() * 100),
			Currency: "RUB",
		},
		Invoice: domain.Invoice{
			Description: "Test",
			OrderNo:     order.ID.String(),
		},
		PaymentData: domain.PaymentData{
			PaymentMethod:          "BankCard",
			PaymentInstrumentTitle: "410000XXXXXX0001",
		},
	}
	s.orderService.On("GetByID", uuid.MustParse(props.Invoice.OrderNo)).Once().Return(&order, nil)

	err = s.uc.Execute(props)

	s.a.Error(err)
}

func (s *payMasterWebHookUseCaseTestSuite) TestExecute_UpdatingOrderError() {
	order := s.createOrder()
	err := order.UpdateStatus(domain.OrderStatusAwaitingPayment)
	s.a.NoError(err)
	order.UpdatePayment(domain.Payment{
		ID:  "test_id",
		URL: "https://test.test",
	})
	props := domain.PayMasterWebHookUseCaseProps{
		ID:         order.Payment.ID,
		Created:    time.Now().UTC(),
		TestMode:   true,
		Status:     "Settled",
		MerchantID: "54a1a08f-b344-44b1-9f06-fe5868170e4c",
		Amount: domain.Amount{
			Value:    float32(order.BuyPrice() * 100),
			Currency: "RUB",
		},
		Invoice: domain.Invoice{
			Description: "Test",
			OrderNo:     order.ID.String(),
		},
		PaymentData: domain.PaymentData{
			PaymentMethod:          "BankCard",
			PaymentInstrumentTitle: "410000XXXXXX0001",
		},
	}
	s.orderService.On("GetByID", uuid.MustParse(props.Invoice.OrderNo)).Once().Return(&order, nil)
	updateProps := domain.UpdateOrderProps{
		OrderID: order.ID,
		Status:  domain.OrderStatusCompleted,
	}
	s.orderService.On("Update", updateProps).Once().Return(nil, errors.New("test"))

	err = s.uc.Execute(props)

	s.a.Error(err)
}

func (s *payMasterWebHookUseCaseTestSuite) TestExecute_GoodCase() {
	order := s.createOrder()
	err := order.UpdateStatus(domain.OrderStatusAwaitingPayment)
	s.a.NoError(err)
	order.UpdatePayment(domain.Payment{
		ID:  "test_id",
		URL: "https://test.test",
	})
	props := domain.PayMasterWebHookUseCaseProps{
		ID:         order.Payment.ID,
		Created:    time.Now().UTC(),
		TestMode:   true,
		Status:     "Settled",
		MerchantID: "54a1a08f-b344-44b1-9f06-fe5868170e4c",
		Amount: domain.Amount{
			Value:    float32(order.BuyPrice() * 100),
			Currency: "RUB",
		},
		Invoice: domain.Invoice{
			Description: "Test",
			OrderNo:     order.ID.String(),
		},
		PaymentData: domain.PaymentData{
			PaymentMethod:          "BankCard",
			PaymentInstrumentTitle: "410000XXXXXX0001",
		},
	}
	s.orderService.On("GetByID", uuid.MustParse(props.Invoice.OrderNo)).Once().Return(&order, nil)
	updateProps := domain.UpdateOrderProps{
		OrderID: order.ID,
		Status:  domain.OrderStatusCompleted,
	}
	s.orderService.On("Update", updateProps).Once().Return(&order, nil)

	err = s.uc.Execute(props)

	s.a.NoError(err)
}

func (s *payMasterWebHookUseCaseTestSuite) createOrder() domain.Order {
	s.T().Helper()
	tc := domain.NewTicketCategory(uuid.New(), 1000, "test1", nil)
	ticket, err := domain.NewTicket(tc, nil)
	s.a.NoError(err)
	o, err := domain.NewOrder(
		"test name",
		"test@test.te",
		"79997776644",
		[]domain.Ticket{*ticket},
	)
	s.a.NoError(err)
	return *o
}
