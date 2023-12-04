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
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type orderControllerTestSuite struct {
	suite.Suite
	a                  *assert.Assertions
	db                 *sqlx.DB
	eventRepo          domain.EventRepo
	orderRepo          domain.OrderRepo
	ticketCategoryRepo domain.TicketCategoryRepo
	pmClient           *paymasterclient.MockPayMasterClient
	promocodeRepo      domain.PromocodeRepo
	svc                domain.OrderService
	controller         transport.OrderController[fiber.Ctx]
	app                *fiber.App
}

func TestOrderControllerTestSuite(t *testing.T) {
	suite.Run(t, &orderControllerTestSuite{})
}

func (s *orderControllerTestSuite) SetupSuite() {
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
	pmClient := paymasterclient.NewMockPayMasterClient(s.T())
	s.pmClient = pmClient
	logger := log.NewLogger(cfg)
	svc, err := service.NewOrderService(
		orderRepo,
		tcRepo,
		promocodeRepo,
		pmClient,
		cfg,
		logger,
	)
	s.a.NoError(err)
	s.svc = svc

	controller, err := transport.NewOrderController(s.svc, logger)
	s.a.NoError(err)
	s.controller = controller

	s.app = fiber.New()
	api := s.app.Group("/api")
	v1 := api.Group("/v1")
	v1.Post("/order", s.controller.Create)
}

func (s *orderControllerTestSuite) TearDownSuite() {
	err := storage.RecreateSchema(s.db)
	s.a.NoError(err)
}

func (s *orderControllerTestSuite) TearDownTest() {
	err := storage.CleanAllTables(s.db)
	s.a.NoError(err)
	s.pmClient.AssertExpectations(s.T())
}

func (s *orderControllerTestSuite) TestCreate() {
	strPointer := func(s string) *string {
		return &s
	}
	tc11, tc12, tc2, p1, p2 := s.createEntities()

	tests := []struct {
		name         string
		method       string
		url          string
		body         transport.CreateOrderRequestDTO
		expectedCode int
	}{
		{
			name:   "incorrect name",
			method: "POST",
			url:    "/api/v1/order",
			body: transport.CreateOrderRequestDTO{
				Name:              "te",
				Email:             "test@te.te",
				Phone:             "79999999999",
				PromocodeID:       strPointer(p1.ID.String()),
				TicketCategoryIDs: []string{tc11.ID.String()},
			},
			expectedCode: http.StatusBadRequest,
		},
		{
			name:   "incorrect email",
			method: "POST",
			url:    "/api/v1/order",
			body: transport.CreateOrderRequestDTO{
				Name:              "testtest",
				Email:             "test",
				Phone:             "79999999999",
				PromocodeID:       strPointer(p1.ID.String()),
				TicketCategoryIDs: []string{tc11.ID.String()},
			},
			expectedCode: http.StatusBadRequest,
		},
		{
			name:   "incorrect phone",
			method: "POST",
			url:    "/api/v1/order",
			body: transport.CreateOrderRequestDTO{
				Name:              "testtest",
				Email:             "test@te.te",
				Phone:             "ddd",
				PromocodeID:       strPointer(p1.ID.String()),
				TicketCategoryIDs: []string{tc11.ID.String()},
			},
			expectedCode: http.StatusBadRequest,
		},
		{
			name:   "incorrect promocode id",
			method: "POST",
			url:    "/api/v1/order",
			body: transport.CreateOrderRequestDTO{
				Name:              "testtest",
				Email:             "test@te.te",
				Phone:             "79999999999",
				PromocodeID:       strPointer("test"),
				TicketCategoryIDs: []string{tc11.ID.String()},
			},
			expectedCode: http.StatusBadRequest,
		},
		{
			name:   "incorrect ticket category id",
			method: "POST",
			url:    "/api/v1/order",
			body: transport.CreateOrderRequestDTO{
				Name:              "testtest",
				Email:             "test@te.te",
				Phone:             "79999999999",
				PromocodeID:       strPointer(p1.ID.String()),
				TicketCategoryIDs: []string{},
			},
			expectedCode: http.StatusBadRequest,
		},
		{
			name:   "good case with unapplied promocode. Trying to use promocode for another ticket category",
			method: "POST",
			url:    "/api/v1/order",
			body: transport.CreateOrderRequestDTO{
				Name:              "testtest",
				Email:             "test@te.te",
				Phone:             "79999999999",
				PromocodeID:       strPointer(p2.ID.String()),
				TicketCategoryIDs: []string{tc11.ID.String()},
			},
			expectedCode: http.StatusCreated,
		},
		{
			name:   "good case",
			method: "POST",
			url:    "/api/v1/order",
			body: transport.CreateOrderRequestDTO{
				Name:              "testtest",
				Email:             "test@te.te",
				Phone:             "79999999999",
				PromocodeID:       strPointer(p1.ID.String()),
				TicketCategoryIDs: []string{tc11.ID.String()},
			},
			expectedCode: http.StatusCreated,
		},
		{
			name:   "good case without promo",
			method: "POST",
			url:    "/api/v1/order",
			body: transport.CreateOrderRequestDTO{
				Name:              "testtest",
				Email:             "test@te.te",
				Phone:             "79999999999",
				PromocodeID:       nil,
				TicketCategoryIDs: []string{tc11.ID.String()},
			},
			expectedCode: http.StatusCreated,
		},
		{
			name:   "good case with 2 ticket categories. Promocode can apply only for 1",
			method: "POST",
			url:    "/api/v1/order",
			body: transport.CreateOrderRequestDTO{
				Name:              "testtest",
				Email:             "test@te.te",
				Phone:             "79999999999",
				PromocodeID:       strPointer(p1.ID.String()),
				TicketCategoryIDs: []string{tc11.ID.String(), tc2.ID.String()},
			},
			expectedCode: http.StatusCreated,
		},
		{
			name:   "good case with 2 ticket categories. Promocode can apply only for 2",
			method: "POST",
			url:    "/api/v1/order",
			body: transport.CreateOrderRequestDTO{
				Name:              "testtest",
				Email:             "test@te.te",
				Phone:             "79999999999",
				PromocodeID:       strPointer(p1.ID.String()),
				TicketCategoryIDs: []string{tc11.ID.String(), tc12.ID.String()},
			},
			expectedCode: http.StatusCreated,
		},
	}
	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			s.pmClient.On("CreateInvoice", mock.Anything).Return(
				&paymasterclient.CreateInvoiceResponseDTO{PaymentID: "test", URL: "https://test.test"},
				nil,
			)
			b := s.generateBodyString(tt.body)
			req, err := http.NewRequest(tt.method, tt.url, b)
			s.a.NoError(err)
			req.Header.Set("Content-Type", "application/json")
			resp, err := s.app.Test(req, 1000000000)
			s.a.NoError(err)
			defer resp.Body.Close()
			var dto transport.OrderResponseDTO
			err = json.NewDecoder(resp.Body).Decode(&dto)
			s.a.NoError(err)
			s.Equalf(tt.expectedCode, resp.StatusCode, tt.name)
			err = storage.CleanOrders(s.db)
			s.a.NoError(err)
		})
	}
}

func (s *orderControllerTestSuite) generateBodyString(createEventRequestDTO transport.CreateOrderRequestDTO) io.Reader {
	s.T().Helper()
	b, err := json.Marshal(createEventRequestDTO)
	c := string(b)
	_ = c
	s.a.NoError(err)
	return bytes.NewReader(b)
}

func (s *orderControllerTestSuite) createEntities() (
	domain.TicketCategory,
	domain.TicketCategory,
	domain.TicketCategory,
	domain.Promocode,
	domain.Promocode,
) {
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

	tc11 := domain.NewTicketCategory(event.ID, 1000, "Common ticket", nil)
	err = s.ticketCategoryRepo.Create(tc11)
	s.a.NoError(err)
	tc12 := domain.NewTicketCategory(event.ID, 1500, "Common ticket 2", nil)
	err = s.ticketCategoryRepo.Create(tc12)
	s.a.NoError(err)
	tc2 := domain.NewTicketCategory(event.ID, 2000, "Vip ticket", nil)
	err = s.ticketCategoryRepo.Create(tc2)
	s.a.NoError(err)

	dv := uint(500)
	p1, err := domain.NewPromocode(1, &dv, nil, []domain.TicketCategory{tc11, tc12})
	s.a.NoError(err)
	err = s.promocodeRepo.Create(*p1)
	s.a.NoError(err)
	p2, err := domain.NewPromocode(1, &dv, nil, []domain.TicketCategory{tc2})
	s.a.NoError(err)
	err = s.promocodeRepo.Create(*p2)
	s.a.NoError(err)
	return tc11, tc12, tc2, *p1, *p2
}
