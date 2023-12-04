package storage

import (
	"errors"
	"fmt"

	"go-ticketos/internal/domain"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type ticketCategoryPSQLRepo struct {
	db      *sqlx.DB
	adapter TicketCategoryAdapter
}

var _ domain.TicketCategoryRepo = (*ticketCategoryPSQLRepo)(nil)

// nolint: revive
func NewTicketCategoryRepo(db *sqlx.DB) (*ticketCategoryPSQLRepo, error) {
	if db == nil {
		return nil, errors.New("NewEventRepo: db is null")
	}
	return &ticketCategoryPSQLRepo{
		db:      db,
		adapter: NewTickerCategoryAdapter(),
	}, nil
}

func (o ticketCategoryPSQLRepo) Create(tc domain.TicketCategory) error {
	schema := o.adapter.ToSchema(tc)
	rows := []string{
		"ticket_categories.id",
		"ticket_categories.event_id",
		"ticket_categories.price",
		"ticket_categories.name",
		"ticket_categories.description",
		"ticket_categories.created_at",
		"ticket_categories.updated_at",
	}
	stmt := "INSERT INTO ticket_categories %s VALUES %s"
	stmt = ToNamedInsert(stmt, rows)
	if _, err := o.db.NamedExec(stmt, schema); err != nil {
		return fmt.Errorf("Create: can not insert row: %w", err)
	}
	return nil
}

func (o ticketCategoryPSQLRepo) GetByID(id uuid.UUID) (*domain.TicketCategory, error) {
	stmt := "SELECT %s FROM ticket_categories WHERE id = $1"
	rows := []string{
		"ticket_categories.id",
		"ticket_categories.event_id",
		"ticket_categories.price",
		"ticket_categories.name",
		"ticket_categories.description",
		"ticket_categories.created_at",
		"ticket_categories.updated_at",
	}
	stmt = ToSelect(stmt, rows)
	var s ticketCategorySchema
	if err := o.db.Get(&s, stmt, id); err != nil {
		return nil, fmt.Errorf("GetByID: can not get data from db: %w", err)
	}
	dom := o.adapter.ToDomain(s)
	return &dom, nil
}
