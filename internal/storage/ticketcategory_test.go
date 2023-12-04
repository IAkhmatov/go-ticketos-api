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

type ticketCategoryRepoTestSuite struct {
	suite.Suite
	a         *assert.Assertions
	db        *sqlx.DB
	eventRepo domain.EventRepo
	repo      domain.TicketCategoryRepo
}

func TestTicketCategoryRepoTestSuite(t *testing.T) {
	suite.Run(t, &ticketCategoryRepoTestSuite{})
}

func (s *ticketCategoryRepoTestSuite) SetupSuite() {
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

	repo, err := storage.NewTicketCategoryRepo(db)
	s.a.NoError(err)
	s.repo = repo
}

func (s *ticketCategoryRepoTestSuite) TearDownSuite() {
	err := storage.RecreateSchema(s.db)
	s.a.NoError(err)
}

func (s *ticketCategoryRepoTestSuite) TearDownTest() {
	err := storage.CleanAllTables(s.db)
	s.a.NoError(err)
}

func (s *ticketCategoryRepoTestSuite) TestCreate() {
	d := "desc"
	event := domain.NewEvent(
		"name",
		nil,
		"place",
		18,
		time.Now().UTC().Add(1*time.Hour),
		time.Now().UTC().Add(2*time.Hour),
	)
	err := s.eventRepo.Create(event)
	s.a.NoError(err)
	tc := domain.NewTicketCategory(
		event.ID,
		1000,
		"test",
		&d,
	)

	err = s.repo.Create(tc)

	s.a.NoError(err)
}

func (s *ticketCategoryRepoTestSuite) TestGetByID() {
	dbtc := s.createTicketCategory()

	tc, err := s.repo.GetByID(dbtc.ID)

	s.a.NoError(err)
	s.a.Equal(dbtc.ID, tc.ID)
	s.a.Equal(dbtc.EventID, tc.EventID)
	s.a.Equal(dbtc.Price, tc.Price)
	s.a.Equal(dbtc.Name, tc.Name)
	s.a.Equal(dbtc.Description, tc.Description)
	s.a.WithinDuration(dbtc.CreatedAt, tc.CreatedAt, time.Millisecond)
	s.a.WithinDuration(dbtc.UpdatedAt, tc.UpdatedAt, time.Millisecond)
}

func (s *ticketCategoryRepoTestSuite) createTicketCategory() domain.TicketCategory {
	s.T().Helper()
	d := "desc"
	event := domain.NewEvent(
		"name",
		nil,
		"place",
		18,
		time.Now().UTC().Add(1*time.Hour),
		time.Now().UTC().Add(2*time.Hour),
	)
	err := s.eventRepo.Create(event)
	s.a.NoError(err)
	tc := domain.NewTicketCategory(
		event.ID,
		1000,
		"test",
		&d,
	)
	err = s.repo.Create(tc)
	s.a.NoError(err)
	return tc
}
