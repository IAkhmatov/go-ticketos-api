package domain_test

import (
	"testing"
	"time"

	"go-ticketos/internal/domain"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestNewTicketCategory(t *testing.T) {
	a := assert.New(t)
	type args struct {
		eventID     uuid.UUID
		price       uint
		name        string
		description *string
	}
	testDesc := "test desc"
	tests := []struct {
		name string
		args args
	}{
		{
			name: "without desc",
			args: args{
				eventID:     uuid.New(),
				price:       1000,
				name:        "test",
				description: nil,
			},
		},
		{
			name: "with desc",
			args: args{
				eventID:     uuid.New(),
				price:       1000,
				name:        "test",
				description: &testDesc,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tc := domain.NewTicketCategory(tt.args.eventID, tt.args.price, tt.args.name, tt.args.description)
			a.NotEqual(uuid.Nil, tc.ID)
			a.Equal(tt.args.name, tc.Name)
			a.Equal(tt.args.price, tc.Price)
			a.Equal(tt.args.name, tc.Name)
			if tt.args.description != nil {
				a.Equal(*tt.args.description, *tc.Description)
			} else {
				a.Nil(tc.Description)
			}
			a.False(tc.CreatedAt.IsZero())
			a.False(tc.UpdatedAt.IsZero())
			a.WithinDuration(tc.CreatedAt, tc.UpdatedAt, time.Millisecond)
		})
	}
}

func TestTicketCategory_CanUsePromocode(t *testing.T) {
	a := assert.New(t)
	tc1 := domain.NewTicketCategory(uuid.New(), 1000, "test", nil)
	tc2 := domain.NewTicketCategory(uuid.New(), 1000, "test", nil)
	dp := uint(50)
	p1, err := domain.NewPromocode(1, nil, &dp, []domain.TicketCategory{tc1})
	a.NoError(err)
	p2, err := domain.NewPromocode(1, nil, &dp, []domain.TicketCategory{tc2})
	a.NoError(err)

	a.True(tc1.CanUsePromocode(*p1))
	a.False(tc1.CanUsePromocode(*p2))
}
