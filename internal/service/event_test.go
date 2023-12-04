package service_test

import (
	"errors"
	"testing"
	"time"

	"go-ticketos/internal/domain"
	"go-ticketos/internal/service"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type eventServiceTestSuite struct {
	suite.Suite
	a         *assert.Assertions
	eventRepo *domain.MockEventRepo
	svc       domain.EventService
}

func TestEventServiceTestSuite(t *testing.T) {
	suite.Run(t, &eventServiceTestSuite{})
}

func (s *eventServiceTestSuite) SetupTest() {
	s.a = assert.New(s.T())
	s.eventRepo = domain.NewMockEventRepo(s.T())
	svc, err := service.NewEventService(s.eventRepo)
	s.a.NoError(err)
	s.svc = svc
}

func (s *eventServiceTestSuite) TeardownTest() {
	s.eventRepo.AssertExpectations(s.T())
}

func (s *eventServiceTestSuite) TestCreate_InsertError() {
	props := domain.CreateEventProps{
		Name:        "Test",
		Description: nil,
		Place:       "test place",
		AgeRating:   16,
		StartAt:     time.Now().UTC().Add(1 * time.Hour),
		EndAt:       time.Now().UTC().Add(2 * time.Hour),
	}
	s.eventRepo.On("Create", mock.Anything).Once().Return(errors.New("test"))

	e, err := s.svc.Create(props)

	s.a.Error(err)
	s.a.Nil(e)
}

func (s *eventServiceTestSuite) TestCreate_GoodCase() {
	props := domain.CreateEventProps{
		Name:        "Test",
		Description: nil,
		Place:       "test place",
		AgeRating:   16,
		StartAt:     time.Now().UTC().Add(1 * time.Hour),
		EndAt:       time.Now().UTC().Add(2 * time.Hour),
	}
	s.eventRepo.On("Create", mock.Anything).Once().Return(nil)

	e, err := s.svc.Create(props)

	s.a.NoError(err)
	s.a.Equal(props.Name, e.Name)
}
