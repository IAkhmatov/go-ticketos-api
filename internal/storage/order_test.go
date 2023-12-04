package storage_test

import (
	"testing"
	"time"

	"go-ticketos/internal/config"
	"go-ticketos/internal/domain"
	"go-ticketos/internal/storage"

	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type orderRepoTestSuite struct {
	suite.Suite
	a                  *assert.Assertions
	db                 *sqlx.DB
	repo               domain.OrderRepo
	ticketCategoryRepo domain.TicketCategoryRepo
	eventRepo          domain.EventRepo
	promocodeRepo      domain.PromocodeRepo
}

func TestOrderRepoTestSuite(t *testing.T) {
	suite.Run(t, &orderRepoTestSuite{})
}

func (s *orderRepoTestSuite) SetupSuite() {
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

	repo, err := storage.NewOrderRepo(db)
	s.a.NoError(err)
	s.repo = repo

	tcRepo, err := storage.NewTicketCategoryRepo(db)
	s.a.NoError(err)
	s.ticketCategoryRepo = tcRepo

	pRepo, err := storage.NewPromocodeRepo(db, tcRepo)
	s.a.NoError(err)
	s.promocodeRepo = pRepo

	eRepo, err := storage.NewEventRepo(db)
	s.a.NoError(err)
	s.eventRepo = eRepo
}

func (s *orderRepoTestSuite) TearDownSuite() {
	err := storage.RecreateSchema(s.db)
	s.a.NoError(err)
}

func (s *orderRepoTestSuite) TearDownTest() {
	err := storage.CleanAllTables(s.db)
	s.a.NoError(err)
}

func (s *orderRepoTestSuite) TestCreate() {
	order := s.prepareOrder()

	err := s.repo.Create(order)

	s.a.NoError(err)
}

func (s *orderRepoTestSuite) TestUpdate() {
	order := s.prepareOrder()
	err := s.repo.Create(order)
	s.a.NoError(err)
	err = order.UpdateStatus(domain.OrderStatusAwaitingPayment)
	s.a.NoError(err)
	order.UpdatePayment(domain.Payment{ID: "test", URL: "https://test.test"})

	err = s.repo.Update(order)

	s.a.NoError(err)
}

func (s *orderRepoTestSuite) TestGetByID() {
	order := s.prepareOrder()
	err := s.repo.Create(order)
	s.a.NoError(err)
	err = order.UpdateStatus(domain.OrderStatusAwaitingPayment)
	s.a.NoError(err)
	order.UpdatePayment(domain.Payment{ID: "test", URL: "https://test.test"})
	err = s.repo.Update(order)
	s.a.NoError(err)

	dbOrder, err := s.repo.GetByID(order.ID)

	s.a.NoError(err)
	s.a.Equal(order.ID, dbOrder.ID)
	s.a.Equal(order.Name, dbOrder.Name)
	s.a.Equal(order.Email, dbOrder.Email)
	s.a.Equal(order.Phone, dbOrder.Phone)
	s.a.Equal(order.Status, dbOrder.Status)
	s.a.Equal(order.Payment.ID, order.Payment.ID)
	s.a.Equal(order.Payment.URL, order.Payment.URL)
	s.a.Equal(len(order.Tickets), len(dbOrder.Tickets))
	s.a.Equal(order.Tickets[0].ID, dbOrder.Tickets[0].ID)
	s.a.WithinDuration(order.CreatedAt, dbOrder.CreatedAt, time.Millisecond)
	s.a.WithinDuration(order.UpdatedAt, dbOrder.UpdatedAt, time.Millisecond)
}

func (s *orderRepoTestSuite) prepareOrder() domain.Order {
	s.T().Helper()
	e := domain.NewEvent(
		"name",
		nil,
		"place",
		15,
		time.Now().UTC().Add(10*24*time.Hour),
		time.Now().UTC().Add(20*24*time.Hour),
	)
	err := s.eventRepo.Create(e)
	s.a.NoError(err)
	tc := domain.NewTicketCategory(
		e.ID,
		5000,
		"name",
		nil)
	err = s.ticketCategoryRepo.Create(tc)
	s.a.NoError(err)
	dp := uint(50)
	p, err := domain.NewPromocode(1, nil, &dp, []domain.TicketCategory{tc})
	s.a.NoError(err)
	err = s.promocodeRepo.Create(*p)
	s.a.NoError(err)
	ticket, err := domain.NewTicket(tc, p)
	s.a.NoError(err)
	order, err := domain.NewOrder(
		"test",
		"email",
		"79999999999",
		[]domain.Ticket{*ticket},
	)
	s.a.NoError(err)
	return *order
}
