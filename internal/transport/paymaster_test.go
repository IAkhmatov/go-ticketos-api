package transport_test

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"testing"
	"time"

	"go-ticketos/internal/config"
	"go-ticketos/internal/domain"
	"go-ticketos/internal/service"
	"go-ticketos/internal/storage"
	"go-ticketos/internal/transport"
	"go-ticketos/pkg/log"
	"go-ticketos/pkg/paymasterclient"

	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type payMasterControllerTestSuite struct {
	suite.Suite
	a                  *assert.Assertions
	db                 *sqlx.DB
	eventRepo          domain.EventRepo
	orderRepo          domain.OrderRepo
	ticketCategoryRepo domain.TicketCategoryRepo
	promocodeRepo      domain.PromocodeRepo
	orderService       domain.OrderService
	useCase            domain.PayMasterWebHookUseCase
	controller         transport.PayMasterController[fiber.Ctx]
	app                *fiber.App
}

func TestPayMasterControllerTestSuite(t *testing.T) {
	suite.Run(t, &payMasterControllerTestSuite{})
}

func (s *payMasterControllerTestSuite) TearDownSuite() {
	err := storage.RecreateSchema(s.db)
	s.a.NoError(err)
}

func (s *payMasterControllerTestSuite) TearDownTest() {
	err := storage.CleanAllTables(s.db)
	s.a.NoError(err)
}

func (s *payMasterControllerTestSuite) SetupSuite() {
	s.a = assert.New(s.T())

	cfg, err := config.NewConfig()
	s.a.NoError(err)
	db, err := storage.NewSqlxDB(cfg.TestDBConnectString)
	s.a.NoError(err)
	err = storage.RecreateSchema(db)
	s.a.NoError(err)
	err = storage.CreateTables(db)
	s.a.NoError(err)
	s.db = db

	eventRepo, err := storage.NewEventRepo(db)
	s.a.NoError(err)
	s.eventRepo = eventRepo
	orderRepo, err := storage.NewOrderRepo(db)
	s.a.NoError(err)
	s.orderRepo = orderRepo
	tcRepo, err := storage.NewTicketCategoryRepo(db)
	s.a.NoError(err)
	s.ticketCategoryRepo = tcRepo
	promocodeRepo, err := storage.NewPromocodeRepo(db, tcRepo)
	s.a.NoError(err)
	s.promocodeRepo = promocodeRepo
	logger := log.NewLogger(cfg)
	pmClient, err := paymasterclient.NewPayMasterClient(cfg, logger)
	s.a.NoError(err)
	svc, err := service.NewOrderService(
		orderRepo,
		tcRepo,
		promocodeRepo,
		pmClient,
		cfg,
		logger,
	)
	s.a.NoError(err)
	s.orderService = svc
	useCase, err := service.NewPayMasterWebHookUseCase(s.orderService, logger)
	s.a.NoError(err)
	s.useCase = useCase

	controller, err := transport.NewPayMasterController(s.useCase, logger)
	s.a.NoError(err)
	s.controller = controller

	s.app = fiber.New()
	api := s.app.Group("/api")
	v1 := api.Group("/v1")
	v1.Post("/webhook/paymaster", s.controller.WebHook)
}

func (s *payMasterControllerTestSuite) TestWebHook() {
	order := s.prepareOrder()
	err := s.orderRepo.Create(order)
	s.a.NoError(err)
	tests := []struct {
		name         string
		method       string
		url          string
		body         transport.PayMasterWebHookRequestDTO
		expectedCode int
	}{
		{
			name:   "incorrect input",
			method: "POST",
			url:    "/api/v1/webhook/paymaster",
			body: transport.PayMasterWebHookRequestDTO{
				ID:         "123124",
				Created:    time.Now().UTC(),
				TestMode:   true,
				Status:     "Settled",
				MerchantID: "f100e9ef-f838-400e-89a8-0f906362d479",
				Amount: transport.Amount{
					Value:    float32(-1),
					Currency: "RUB",
				},
				Invoice: transport.Invoice{
					Description: "Test",
					OrderNo:     order.ID.String(),
				},
				PaymentData: transport.PaymentData{
					PaymentMethod:          "BankCard",
					PaymentInstrumentTitle: "410000XXXXXX0001",
				},
			},
			expectedCode: http.StatusBadRequest,
		},
		{
			name:   "good case",
			method: "POST",
			url:    "/api/v1/webhook/paymaster",
			body: transport.PayMasterWebHookRequestDTO{
				ID:         order.Payment.ID,
				Created:    time.Now().UTC(),
				TestMode:   true,
				Status:     "Settled",
				MerchantID: "f100e9ef-f838-400e-89a8-0f906362d479",
				Amount: transport.Amount{
					Value:    float32(order.BuyPrice() * 100),
					Currency: "RUB",
				},
				Invoice: transport.Invoice{
					Description: "Test",
					OrderNo:     order.ID.String(),
				},
				PaymentData: transport.PaymentData{
					PaymentMethod:          "BankCard",
					PaymentInstrumentTitle: "410000XXXXXX0001",
				},
			},
			expectedCode: http.StatusOK,
		},
	}
	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			err = storage.CleanOrders(s.db)
			s.a.NoError(err)
			err = s.orderRepo.Create(order)
			s.a.NoError(err)
			b := s.generateBodyString(tt.body)
			req, err2 := http.NewRequest(tt.method, tt.url, b)
			s.a.NoError(err2)
			req.Header.Set("Content-Type", "application/json")
			resp, err3 := s.app.Test(req, 1000)
			s.a.NoError(err3)
			s.Equalf(tt.expectedCode, resp.StatusCode, tt.name)
		})
	}
}

func (s *payMasterControllerTestSuite) generateBodyString(
	payMasterWebHookRequestDTO transport.PayMasterWebHookRequestDTO,
) io.Reader {
	s.T().Helper()
	b, err := json.Marshal(payMasterWebHookRequestDTO)
	c := string(b)
	_ = c
	s.a.NoError(err)
	return bytes.NewReader(b)
}

func (s *payMasterControllerTestSuite) prepareOrder() domain.Order {
	s.T().Helper()
	event := domain.NewEvent(
		"test name",
		nil,
		"test place",
		15,
		time.Now().UTC().Add(1000*time.Hour),
		time.Now().UTC().Add(1100*time.Hour),
	)
	err := s.eventRepo.Create(event)
	s.a.NoError(err)

	tc := domain.NewTicketCategory(event.ID, 1000, "Common ticket", nil)
	err = s.ticketCategoryRepo.Create(tc)
	s.a.NoError(err)

	ticket, err := domain.NewTicket(tc, nil)
	s.a.NoError(err)
	order, err := domain.NewOrder("Test name", "test@test.com", "79997776644", []domain.Ticket{*ticket})
	s.a.NoError(err)
	order.UpdatePayment(domain.Payment{ID: "test", URL: "https://test.test"})
	err = order.UpdateStatus(domain.OrderStatusAwaitingPayment)
	s.a.NoError(err)
	return *order
}
