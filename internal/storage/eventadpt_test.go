// nolint: testpackage
package storage

import (
	"database/sql"
	"testing"
	"time"

	"go-ticketos/internal/domain"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type eventAdapterTestSuite struct {
	suite.Suite
	a       *assert.Assertions
	adapter *eventAdapter
}

func TestEventAdapterTestSuite(t *testing.T) {
	suite.Run(t, &eventAdapterTestSuite{})
}

func (s *eventAdapterTestSuite) SetupSuite() {
	s.a = assert.New(s.T())
	s.adapter = NewEventAdapter()
}

func (s *eventAdapterTestSuite) TestToSchema() {
	type args struct {
		event domain.Event
	}
	desc := "desc"
	evt1 := domain.NewEvent(
		"name",
		&desc,
		"place",
		18,
		time.Now().UTC().Add(1*time.Hour),
		time.Now().UTC().Add(2*time.Hour),
	)
	evt2 := evt1
	evt2.Description = nil

	tests := []struct {
		name string
		args args
		want eventSchema
	}{
		{
			name: "with description",
			args: args{
				event: evt1,
			},
			want: eventSchema{
				ID:   evt1.ID,
				Name: evt1.Name,
				Description: sql.NullString{
					String: *evt1.Description,
					Valid:  true,
				},
				Place:     evt1.Place,
				AgeRating: evt1.AgeRating,
				StartAt:   evt1.StartAt,
				EndAt:     evt1.EndAt,
			},
		},
		{
			name: "without description",
			args: args{
				event: evt2,
			},
			want: eventSchema{
				ID:   evt2.ID,
				Name: evt2.Name,
				Description: sql.NullString{
					Valid: false,
				},
				Place:     evt2.Place,
				AgeRating: evt2.AgeRating,
				StartAt:   evt2.StartAt,
				EndAt:     evt2.EndAt,
			},
		},
	}
	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			currentTest := tt
			actualDto := s.adapter.ToSchema(currentTest.args.event)
			s.a.Equal(currentTest.want.ID, actualDto.ID)
			s.a.Equal(currentTest.want.Name, actualDto.Name)
			s.a.Equal(currentTest.want.Description.Valid, actualDto.Description.Valid)
			s.a.Equal(currentTest.want.Description.String, actualDto.Description.String)
			s.a.Equal(currentTest.want.Place, actualDto.Place)
			s.a.Equal(currentTest.want.AgeRating, actualDto.AgeRating)
			s.a.WithinDuration(currentTest.want.StartAt, actualDto.StartAt, time.Millisecond)
			s.a.WithinDuration(currentTest.want.EndAt, actualDto.EndAt, time.Millisecond)
		})
	}
}
