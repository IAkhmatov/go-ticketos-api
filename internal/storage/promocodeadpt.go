package storage

import (
	"database/sql"
	"time"

	"go-ticketos/internal/domain"

	"github.com/google/uuid"
)

type PromocodeAdapter interface {
	ToSchema(promocode domain.Promocode) promocodeSchema
	ToDomain(promocode promocodeSchema) domain.Promocode
}

type promocodeAdapter struct{}

var _ PromocodeAdapter = (*promocodeAdapter)(nil)

// nolint: revive
func NewPromocodeAdapter() *promocodeAdapter {
	return &promocodeAdapter{}
}

func (a *promocodeAdapter) ToSchema(promocode domain.Promocode) promocodeSchema {
	schema := promocodeSchema{
		ID:        promocode.ID,
		LimitUse:  promocode.LimitUse,
		CreatedAt: promocode.CreatedAt,
		UpdatedAt: promocode.UpdatedAt,
	}
	if promocode.DiscountPercent != nil {
		schema.DiscountPercent = sql.NullInt32{
			Int32: int32(*promocode.DiscountPercent),
			Valid: true,
		}
	}
	if promocode.DiscountValue != nil {
		schema.DiscountValue = sql.NullInt32{
			Int32: int32(*promocode.DiscountValue),
			Valid: true,
		}
	}
	if len(promocode.TicketCategories) > 0 {
		var ptcs []promocodeTicketCategorySchema
		for _, tc := range promocode.TicketCategories {
			now := time.Now().UTC()
			ptc := promocodeTicketCategorySchema{
				ID:               uuid.New(),
				TicketCategoryID: tc.ID,
				PromocodeID:      promocode.ID,
				CreatedAt:        now,
				UpdatedAt:        now,
			}
			ptcs = append(ptcs, ptc)
		}
		schema.TicketCategories = ptcs
	}
	return schema
}

func (a *promocodeAdapter) ToDomain(promocode promocodeSchema) domain.Promocode {
	dom := domain.Promocode{
		ID:        promocode.ID,
		LimitUse:  promocode.LimitUse,
		CreatedAt: promocode.CreatedAt,
		UpdatedAt: promocode.UpdatedAt,
	}
	if promocode.DiscountPercent.Valid {
		u := uint(promocode.DiscountPercent.Int32)
		dom.DiscountPercent = &u
	}
	if promocode.DiscountValue.Valid {
		u := uint(promocode.DiscountValue.Int32)
		dom.DiscountValue = &u
	}
	// add adpt ticket categories
	return dom
}
