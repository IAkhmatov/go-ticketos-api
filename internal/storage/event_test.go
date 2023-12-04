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

type eventRepoTestSuite struct {
	suite.Suite
	a    *assert.Assertions
	db   *sqlx.DB
	repo domain.EventRepo
}

func TestEventRepoTestSuite(t *testing.T) {
	suite.Run(t, &eventRepoTestSuite{})
}

func (s *eventRepoTestSuite) SetupSuite() {
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

	repo, err := storage.NewEventRepo(db)
	s.a.NoError(err)
	s.repo = repo
}

func (s *eventRepoTestSuite) TearDownSuite() {
	err := storage.RecreateSchema(s.db)
	s.a.NoError(err)
}

func (s *eventRepoTestSuite) TearDownTest() {
	err := storage.CleanAllTables(s.db)
	s.a.NoError(err)
}

func (s *eventRepoTestSuite) TestCreate() {
	event := domain.NewEvent(
		"name",
		nil,
		"place",
		18,
		time.Now().UTC().Add(1*time.Hour),
		time.Now().UTC().Add(2*time.Hour),
	)

	err := s.repo.Create(event)

	s.a.NoError(err)
}
