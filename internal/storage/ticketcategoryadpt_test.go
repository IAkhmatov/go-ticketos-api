// nolint: testpackage
package storage

import (
	"database/sql"
	"testing"
	"time"

	"go-ticketos/internal/domain"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type ticketCategoryAdapterTestSuite struct {
	suite.Suite
	a       *assert.Assertions
	adapter *ticketCategoryAdapter
}

func TestTicketCategoryAdapterTestSuite(t *testing.T) {
	suite.Run(t, &ticketCategoryAdapterTestSuite{})
}

func (s *ticketCategoryAdapterTestSuite) SetupSuite() {
	s.a = assert.New(s.T())
	s.adapter = NewTickerCategoryAdapter()
}

func (s *ticketCategoryAdapterTestSuite) TestToSchema() {
	type args struct {
		tc domain.TicketCategory
	}
	desc := "desc"
	tc1 := domain.NewTicketCategory(
		uuid.New(),
		1000,
		"test",
		nil,
	)

	tc2 := tc1
	tc2.Description = &desc

	tests := []struct {
		name string
		args args
	}{
		{
			name: "without description",
			args: args{
				tc: tc1,
			},
		},
		{
			name: "with description",
			args: args{
				tc: tc2,
			},
		},
	}
	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			actualDto := s.adapter.ToSchema(tt.args.tc)
			s.a.Equal(tt.args.tc.ID, actualDto.ID)
			s.a.Equal(tt.args.tc.EventID, actualDto.EventID)
			s.a.Equal(tt.args.tc.Price, actualDto.Price)
			s.a.Equal(tt.args.tc.Name, actualDto.Name)
			if tt.args.tc.Description != nil {
				s.a.Equal(*tt.args.tc.Description, actualDto.Description.String)
			} else {
				s.a.False(actualDto.Description.Valid)
			}
			s.a.False(actualDto.CreatedAt.IsZero())
			s.a.False(actualDto.UpdatedAt.IsZero())
			s.a.WithinDuration(actualDto.CreatedAt, actualDto.UpdatedAt, time.Millisecond)
		})
	}
}

func (s *ticketCategoryAdapterTestSuite) TestToDomain() {
	type args struct {
		tc ticketCategorySchema
	}

	now := time.Now().UTC()
	tests := []struct {
		name string
		args args
	}{
		{
			name: "without description",
			args: args{
				tc: ticketCategorySchema{
					ID:        uuid.New(),
					EventID:   uuid.New(),
					Price:     5000,
					Name:      "Test",
					CreatedAt: now,
					UpdatedAt: now,
				},
			},
		},
		{
			name: "with description",
			args: args{
				tc: ticketCategorySchema{
					ID:      uuid.New(),
					EventID: uuid.New(),
					Price:   1000,
					Name:    "Test",
					Description: sql.NullString{
						String: "Test desc",
						Valid:  true,
					},
					CreatedAt: now,
					UpdatedAt: now,
				},
			},
		},
	}
	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			actualDto := s.adapter.ToDomain(tt.args.tc)
			s.a.Equal(tt.args.tc.ID, actualDto.ID)
			s.a.Equal(tt.args.tc.EventID, actualDto.EventID)
			s.a.Equal(tt.args.tc.Price, actualDto.Price)
			s.a.Equal(tt.args.tc.Name, actualDto.Name)
			if tt.args.tc.Description.Valid {
				s.a.Equal(tt.args.tc.Description.String, *actualDto.Description)
			} else {
				s.a.Nil(actualDto.Description)
			}
			s.a.False(actualDto.CreatedAt.IsZero())
			s.a.False(actualDto.UpdatedAt.IsZero())
			s.a.WithinDuration(actualDto.CreatedAt, actualDto.UpdatedAt, time.Millisecond)
		})
	}
}
