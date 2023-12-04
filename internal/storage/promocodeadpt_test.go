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

type promocodeAdapterTestSuite struct {
	suite.Suite
	a       *assert.Assertions
	adapter PromocodeAdapter
}

func TestPromocodeAdapterTestSuite(t *testing.T) {
	suite.Run(t, &promocodeAdapterTestSuite{})
}

func (s *promocodeAdapterTestSuite) SetupSuite() {
	s.a = assert.New(s.T())
	s.adapter = NewPromocodeAdapter()
}

func (s *promocodeAdapterTestSuite) TestToSchema() {
	type args struct {
		promocode domain.Promocode
	}

	tc := domain.NewTicketCategory(uuid.New(), 1000, "name", nil)
	dv := uint(3)
	p1, err := domain.NewPromocode(10, &dv, nil, []domain.TicketCategory{tc})
	s.a.NoError(err)
	dp := uint(5)
	p2, err := domain.NewPromocode(10, nil, &dp, []domain.TicketCategory{tc})
	s.a.NoError(err)

	tests := []struct {
		name string
		args args
		want promocodeSchema
	}{
		{
			name: "with value",
			args: args{
				promocode: *p1,
			},
			want: promocodeSchema{
				ID:       p1.ID,
				LimitUse: p1.LimitUse,
				DiscountValue: sql.NullInt32{
					Int32: 3,
					Valid: true,
				},
				TicketCategories: []promocodeTicketCategorySchema{
					{
						ID:               uuid.New(),
						TicketCategoryID: p1.TicketCategories[0].ID,
						PromocodeID:      p1.ID,
						CreatedAt:        time.Time{},
						UpdatedAt:        time.Time{},
					},
				},
				CreatedAt: p1.CreatedAt,
				UpdatedAt: p1.UpdatedAt,
			},
		},
		{
			name: "with percent",
			args: args{
				promocode: *p2,
			},
			want: promocodeSchema{
				ID:       p2.ID,
				LimitUse: p2.LimitUse,
				DiscountPercent: sql.NullInt32{
					Int32: 5,
					Valid: true,
				},
				TicketCategories: []promocodeTicketCategorySchema{
					{
						ID:               uuid.New(),
						TicketCategoryID: p2.TicketCategories[0].ID,
						PromocodeID:      p2.ID,
						CreatedAt:        time.Time{},
						UpdatedAt:        time.Time{},
					},
				},
				CreatedAt: p2.CreatedAt,
				UpdatedAt: p2.UpdatedAt,
			},
		},
	}
	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			currentTest := tt
			actualDto := s.adapter.ToSchema(currentTest.args.promocode)
			s.a.Equal(currentTest.want.ID, actualDto.ID)
			s.a.Equal(currentTest.want.LimitUse, actualDto.LimitUse)
			s.a.Equal(currentTest.want.DiscountPercent.Valid, actualDto.DiscountPercent.Valid)
			s.a.Equal(currentTest.want.DiscountPercent.Int32, actualDto.DiscountPercent.Int32)
			s.a.Equal(currentTest.want.DiscountValue.Valid, actualDto.DiscountValue.Valid)
			s.a.Equal(currentTest.want.DiscountValue.Int32, actualDto.DiscountValue.Int32)
			s.a.WithinDuration(currentTest.want.CreatedAt, actualDto.CreatedAt, time.Millisecond)
			s.a.WithinDuration(currentTest.want.UpdatedAt, actualDto.UpdatedAt, time.Millisecond)
			s.a.Len(actualDto.TicketCategories, 1)
			s.a.Equal(currentTest.want.TicketCategories[0].TicketCategoryID, actualDto.TicketCategories[0].TicketCategoryID)
			s.a.Equal(currentTest.want.TicketCategories[0].PromocodeID, actualDto.TicketCategories[0].PromocodeID)
		})
	}
}

func (s *promocodeAdapterTestSuite) TestToDomain() {
	type args struct {
		promocode promocodeSchema
	}

	now := time.Now().UTC()
	p1 := promocodeSchema{
		ID:       uuid.New(),
		LimitUse: 11,
		DiscountValue: sql.NullInt32{
			Int32: 2000,
			Valid: true,
		},
		TicketCategories: nil,
		CreatedAt:        now,
		UpdatedAt:        now,
	}
	u1 := uint(2000)
	p2 := promocodeSchema{
		ID:       uuid.New(),
		LimitUse: 11,
		DiscountPercent: sql.NullInt32{
			Int32: 20,
			Valid: true,
		},
		TicketCategories: nil,
		CreatedAt:        now,
		UpdatedAt:        now,
	}
	u2 := uint(20)
	tests := []struct {
		name string
		args args
		want domain.Promocode
	}{
		{
			name: "with value",
			args: args{
				promocode: p1,
			},
			want: domain.Promocode{
				ID:               p1.ID,
				LimitUse:         p1.LimitUse,
				DiscountValue:    &u1,
				TicketCategories: nil,
				CreatedAt:        p1.CreatedAt,
				UpdatedAt:        p1.UpdatedAt,
			},
		},
		{
			name: "with percent",
			args: args{
				promocode: p2,
			},
			want: domain.Promocode{
				ID:               p2.ID,
				LimitUse:         p2.LimitUse,
				DiscountPercent:  &u2,
				TicketCategories: nil,
				CreatedAt:        p2.CreatedAt,
				UpdatedAt:        p2.UpdatedAt,
			},
		},
	}
	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			currentTest := tt
			actualDto := s.adapter.ToDomain(currentTest.args.promocode)
			s.a.Equal(currentTest.want.ID, actualDto.ID)
			s.a.Equal(currentTest.want.LimitUse, actualDto.LimitUse)
			if currentTest.want.DiscountValue != nil {
				s.a.Equal(currentTest.want.DiscountValue, actualDto.DiscountValue)
			}
			if currentTest.want.DiscountPercent != nil {
				s.a.Equal(currentTest.want.DiscountPercent, actualDto.DiscountPercent)
			}
			s.a.WithinDuration(currentTest.want.CreatedAt, actualDto.CreatedAt, time.Millisecond)
			s.a.WithinDuration(currentTest.want.UpdatedAt, actualDto.UpdatedAt, time.Millisecond)
		})
	}
}
