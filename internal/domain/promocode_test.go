package domain_test

import (
	"testing"
	"time"

	"go-ticketos/internal/domain"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestNewPromocode(t *testing.T) {
	a := assert.New(t)
	type args struct {
		limitUse         uint
		discountValue    *uint
		discountPercent  *uint
		ticketCategories []domain.TicketCategory
	}
	dv := uint(100)
	dp := uint(30)
	tc := domain.NewTicketCategory(uuid.New(), 1000, "test", nil)
	tests := []struct {
		name        string
		args        args
		isErrResult bool
	}{
		{
			name: "incorrect discount 1",
			args: args{
				limitUse:         1,
				discountValue:    &dv,
				discountPercent:  &dp,
				ticketCategories: []domain.TicketCategory{tc},
			},
			isErrResult: true,
		},
		{
			name: "incorrect discount 2",
			args: args{
				limitUse:         1,
				discountValue:    nil,
				discountPercent:  nil,
				ticketCategories: []domain.TicketCategory{tc},
			},
			isErrResult: true,
		},
		{
			name: "incorrect limit Use",
			args: args{
				limitUse:         0,
				discountValue:    &dv,
				discountPercent:  nil,
				ticketCategories: []domain.TicketCategory{tc},
			},
			isErrResult: true,
		},
		{
			name: "incorrect ticket categories",
			args: args{
				limitUse:         1,
				discountValue:    &dv,
				discountPercent:  nil,
				ticketCategories: nil,
			},
			isErrResult: true,
		},
		{
			name: "good case with value",
			args: args{
				limitUse:         1,
				discountValue:    &dv,
				discountPercent:  nil,
				ticketCategories: []domain.TicketCategory{tc},
			},
			isErrResult: false,
		},
		{
			name: "good case with percent",
			args: args{
				limitUse:         1,
				discountValue:    nil,
				discountPercent:  &dp,
				ticketCategories: []domain.TicketCategory{tc},
			},
			isErrResult: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := domain.NewPromocode(
				tt.args.limitUse,
				tt.args.discountValue,
				tt.args.discountPercent,
				tt.args.ticketCategories,
			)
			if tt.isErrResult {
				a.Error(err)
			} else {
				a.NotEqual(uuid.Nil, got.ID)
				if tt.args.discountValue != nil {
					a.Equal(*tt.args.discountValue, *got.DiscountValue)
				}
				if tt.args.discountPercent != nil {
					a.Equal(*tt.args.discountPercent, *got.DiscountPercent)
				}
				a.Len(got.TicketCategories, len(tt.args.ticketCategories))
				a.False(got.CreatedAt.IsZero())
				a.False(got.UpdatedAt.IsZero())
				a.WithinDuration(got.CreatedAt, got.UpdatedAt, time.Millisecond)
			}
		})
	}
}
