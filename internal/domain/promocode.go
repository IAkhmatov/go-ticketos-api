package domain

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

type (
	Promocode struct {
		ID               uuid.UUID
		LimitUse         uint
		DiscountValue    *uint
		DiscountPercent  *uint
		TicketCategories []TicketCategory
		CreatedAt        time.Time
		UpdatedAt        time.Time
	}

	PromocodeRepo interface {
		Create(promocode Promocode) error
		GetByID(id uuid.UUID) (*Promocode, error)
	}
)

func NewPromocode(
	limitUse uint,
	discountValue *uint,
	discountPercent *uint,
	ticketCategories []TicketCategory,
) (*Promocode, error) {
	now := time.Now().UTC()
	if discountPercent != nil && discountValue != nil {
		return nil, errors.New("NewPromocode: you need to set discount value or discount percent")
	}
	if discountPercent == nil && discountValue == nil {
		return nil, errors.New("NewPromocode: you need to set discount value or discount percent")
	}
	if len(ticketCategories) == 0 {
		return nil, errors.New("NewPromocode: promocode must have at least one ticket category")
	}
	if limitUse == 0 {
		return nil, errors.New("NewPromocode: limitUse must be more than 0")
	}
	promocode := &Promocode{
		ID:        uuid.New(),
		LimitUse:  limitUse,
		CreatedAt: now,
		UpdatedAt: now,
	}
	ticketCategoriesCopy := make([]TicketCategory, len(ticketCategories))
	copy(ticketCategoriesCopy, ticketCategories)
	promocode.TicketCategories = ticketCategoriesCopy
	if discountValue != nil {
		dvCopy := *discountValue
		promocode.DiscountValue = &dvCopy
	}
	if discountPercent != nil {
		dpCopy := *discountPercent
		promocode.DiscountPercent = &dpCopy
	}
	return promocode, nil
}

func (p Promocode) CalcBuyPrice(fullPrice uint) uint {
	if p.DiscountValue != nil {
		return fullPrice - *p.DiscountValue
	}
	if p.DiscountPercent != nil {
		buyPrice := float64(fullPrice) * (1 - float64(*p.DiscountPercent)/100)
		return uint(buyPrice)
	}
	return fullPrice
}
