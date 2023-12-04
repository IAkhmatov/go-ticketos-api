package storage

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type (
	eventSchema struct {
		ID          uuid.UUID      `db:"events.id"`
		Name        string         `db:"events.name"`
		Description sql.NullString `db:"events.description"`
		Place       string         `db:"events.place"`
		AgeRating   int            `db:"events.age_rating"`
		StartAt     time.Time      `db:"events.start_at"`
		EndAt       time.Time      `db:"events.end_at"`
		CreatedAt   time.Time      `db:"events.created_at"`
		UpdatedAt   time.Time      `db:"events.updated_at"`
	}

	orderSchema struct {
		ID         uuid.UUID      `db:"orders.id"`
		Name       string         `db:"orders.name"`
		Tickets    []ticketSchema `db:"-"`
		Email      string         `db:"orders.email"`
		Phone      string         `db:"orders.phone"`
		Status     string         `db:"orders.status"`
		PaymentID  sql.NullString `db:"orders.payment_id"`
		PaymentURL sql.NullString `db:"orders.payment_url"`
		CreatedAt  time.Time      `db:"orders.created_at"`
		UpdatedAt  time.Time      `db:"orders.updated_at"`
	}

	ticketSchema struct {
		ID               uuid.UUID            `db:"tickets.id"`
		OrderID          uuid.UUID            `db:"tickets.order_id"`
		TicketCategory   ticketCategorySchema `db:""`
		TicketCategoryID uuid.UUID            `db:"tickets.ticket_category_id"`
		PromocodeID      uuid.NullUUID        `db:"tickets.promocode_id"`
		FullPrice        uint                 `db:"tickets.full_price"`
		BuyPrice         uint                 `db:"tickets.buy_price"`
		CreatedAt        time.Time            `db:"tickets.created_at"`
		UpdatedAt        time.Time            `db:"tickets.updated_at"`
	}

	ticketCategorySchema struct {
		ID          uuid.UUID      `db:"ticket_categories.id"`
		EventID     uuid.UUID      `db:"ticket_categories.event_id"`
		Price       uint           `db:"ticket_categories.price"`
		Name        string         `db:"ticket_categories.name"`
		Description sql.NullString `db:"ticket_categories.description"`
		CreatedAt   time.Time      `db:"ticket_categories.created_at"`
		UpdatedAt   time.Time      `db:"ticket_categories.updated_at"`
	}

	promocodeSchema struct {
		ID               uuid.UUID                       `db:"promocodes.id"`
		LimitUse         uint                            `db:"promocodes.limit_use"`
		DiscountValue    sql.NullInt32                   `db:"promocodes.discount_value"`
		DiscountPercent  sql.NullInt32                   `db:"promocodes.discount_percent"`
		TicketCategories []promocodeTicketCategorySchema `db:"-"`
		CreatedAt        time.Time                       `db:"promocodes.created_at"`
		UpdatedAt        time.Time                       `db:"promocodes.updated_at"`
	}

	promocodeTicketCategorySchema struct {
		ID               uuid.UUID `db:"promocodes_ticket_categories.id"`
		TicketCategoryID uuid.UUID `db:"promocodes_ticket_categories.ticket_category_id"`
		PromocodeID      uuid.UUID `db:"promocodes_ticket_categories.promocode_id"`
		CreatedAt        time.Time `db:"promocodes_ticket_categories.created_at"`
		UpdatedAt        time.Time `db:"promocodes_ticket_categories.updated_at"`
	}
)
