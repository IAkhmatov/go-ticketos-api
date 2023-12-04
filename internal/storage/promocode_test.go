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

type promocodeRepoTestSuite struct {
	suite.Suite
	a                  *assert.Assertions
	db                 *sqlx.DB
	repo               domain.PromocodeRepo
	ticketCategoryRepo domain.TicketCategoryRepo
	eventRepo          domain.EventRepo
}

func TestPromocodeRepoTestSuite(t *testing.T) {
	suite.Run(t, &promocodeRepoTestSuite{})
}

func (s *promocodeRepoTestSuite) SetupSuite() {
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

	tcRepo, err := storage.NewTicketCategoryRepo(db)
	s.a.NoError(err)
	s.ticketCategoryRepo = tcRepo

	repo, err := storage.NewPromocodeRepo(db, tcRepo)
	s.a.NoError(err)
	s.repo = repo

	eRepo, err := storage.NewEventRepo(db)
	s.a.NoError(err)
	s.eventRepo = eRepo
}

func (s *promocodeRepoTestSuite) TearDownSuite() {
	err := storage.RecreateSchema(s.db)
	s.a.NoError(err)
}

func (s *promocodeRepoTestSuite) TearDownTest() {
	err := storage.CleanAllTables(s.db)
	s.a.NoError(err)
}

func (s *promocodeRepoTestSuite) TestCreate() {
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
	p, err := domain.NewPromocode(
		10,
		nil,
		&dp,
		[]domain.TicketCategory{tc},
	)
	s.a.NoError(err)

	err = s.repo.Create(*p)

	s.a.NoError(err)
}

func (s *promocodeRepoTestSuite) TestGetByID() {
	pdb := s.createPromo()

	p, err := s.repo.GetByID(pdb.ID)

	s.a.NoError(err)
	s.a.Equal(pdb.ID, p.ID)
	s.a.Equal(pdb.LimitUse, p.LimitUse)
	s.a.Equal(pdb.DiscountValue, p.DiscountValue)
	s.a.Equal(pdb.DiscountPercent, p.DiscountPercent)
	s.a.Equal(pdb.DiscountPercent, p.DiscountPercent)
	s.a.Len(p.TicketCategories, 2)
	s.a.WithinDuration(pdb.CreatedAt, p.CreatedAt, time.Millisecond)
	s.a.WithinDuration(pdb.UpdatedAt, p.UpdatedAt, time.Millisecond)
}

func (s *promocodeRepoTestSuite) createPromo() domain.Promocode {
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
	tc2 := domain.NewTicketCategory(
		e.ID,
		4000,
		"name2",
		nil)
	err = s.ticketCategoryRepo.Create(tc2)
	s.a.NoError(err)
	dp := uint(50)
	p, err := domain.NewPromocode(
		10,
		nil,
		&dp,
		[]domain.TicketCategory{tc, tc2},
	)
	s.a.NoError(err)

	err = s.repo.Create(*p)

	s.a.NoError(err)

	return *p
}
