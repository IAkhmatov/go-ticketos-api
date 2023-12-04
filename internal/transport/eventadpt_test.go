package transport_test

import (
	"testing"
	"time"

	"go-ticketos/internal/domain"
	"go-ticketos/internal/transport"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type eventAdapterTestSuite struct {
	suite.Suite
	a       *assert.Assertions
	adapter transport.EventAdapter
}

func TestEventAdapterTestSuite(t *testing.T) {
	suite.Run(t, &eventAdapterTestSuite{})
}

func (s *eventAdapterTestSuite) SetupSuite() {
	s.a = assert.New(s.T())
	s.adapter = transport.NewEventAdapter()
}

func (s *eventAdapterTestSuite) TestToResponseDTO() {
	type args struct {
		event domain.Event
	}
	desc := "desc"
	evt1 := domain.Event{
		ID:          uuid.New(),
		Name:        "name",
		Description: &desc,
		Place:       "place",
		AgeRating:   18,
		StartAt:     time.Now().UTC().Add(1 * time.Hour),
		EndAt:       time.Now().UTC().Add(2 * time.Hour),
		CreatedAt:   time.Now().UTC().Add(-1 * 24 * time.Hour),
		UpdatedAt:   time.Now().UTC().Add(-1 * 20 * time.Hour),
	}
	evt2 := evt1
	evt2.Description = nil
	tests := []struct {
		name string
		args args
		want transport.EventResponseDTO
	}{
		{
			name: "with description",
			args: args{
				event: evt1,
			},
			want: transport.EventResponseDTO{
				ID:          evt1.ID,
				Name:        evt1.Name,
				Description: evt1.Description,
				Place:       evt1.Place,
				AgeRating:   evt1.AgeRating,
				StartAt:     evt1.StartAt,
				EndAt:       evt1.EndAt,
			},
		},
		{
			name: "without description",
			args: args{
				event: evt2,
			},
			want: transport.EventResponseDTO{
				ID:          evt2.ID,
				Name:        evt2.Name,
				Description: evt2.Description,
				Place:       evt2.Place,
				AgeRating:   evt2.AgeRating,
				StartAt:     evt2.StartAt,
				EndAt:       evt2.EndAt,
			},
		},
	}
	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			currentTest := tt
			actualDto := s.adapter.ToResponseDTO(currentTest.args.event)
			s.a.Equal(currentTest.want.ID, actualDto.ID)
			s.a.Equal(currentTest.want.Name, actualDto.Name)
			if currentTest.want.Description != nil {
				s.a.Equal(currentTest.want.Description, actualDto.Description)
			} else {
				s.a.Nil(actualDto.Description)
			}
			s.a.Equal(currentTest.want.Place, actualDto.Place)
			s.a.Equal(currentTest.want.AgeRating, actualDto.AgeRating)
			s.a.WithinDuration(currentTest.want.StartAt, actualDto.StartAt, time.Millisecond)
			s.a.WithinDuration(currentTest.want.EndAt, actualDto.EndAt, time.Millisecond)
		})
	}
}
