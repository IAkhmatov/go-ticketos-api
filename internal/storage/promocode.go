package storage

import (
	"database/sql"
	"errors"
	"fmt"

	"go-ticketos/internal/domain"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type promocodePSQLRepo struct {
	db                 *sqlx.DB
	adapter            PromocodeAdapter
	ticketCategoryRepo domain.TicketCategoryRepo
}

// nolint: revive
func NewPromocodeRepo(
	db *sqlx.DB,
	ticketCategoryRepo domain.TicketCategoryRepo,
) (*promocodePSQLRepo, error) {
	if db == nil {
		return nil, errors.New("NewPromocodeRepo: db is nil")
	}
	return &promocodePSQLRepo{
		db:                 db,
		adapter:            NewPromocodeAdapter(),
		ticketCategoryRepo: ticketCategoryRepo,
	}, nil
}

var _ domain.PromocodeRepo = (*promocodePSQLRepo)(nil)

func (p promocodePSQLRepo) Create(promocode domain.Promocode) error {
	tx, err := p.db.Beginx()
	if err != nil {
		return fmt.Errorf("Create: can not create tx: %w", err)
	}
	defer func() {
		if err = tx.Rollback(); err != nil && !errors.Is(err, sql.ErrTxDone) {
			msg := fmt.Sprintf("Create: can not rollback tx: %s", err)
			// nolint: forbidigo
			fmt.Println(msg)
		}
	}()

	schema := p.adapter.ToSchema(promocode)
	rows := []string{
		"promocodes.id",
		"promocodes.limit_use",
		"promocodes.discount_value",
		"promocodes.discount_percent",
		"promocodes.created_at",
		"promocodes.updated_at",
	}
	stmt := "INSERT INTO promocodes %s VALUES %s"
	stmt = ToNamedInsert(stmt, rows)
	if _, err = tx.NamedExec(stmt, schema); err != nil {
		return fmt.Errorf("Create: can not insert promocodes: %w", err)
	}

	rows = []string{
		"promocodes_ticket_categories.id",
		"promocodes_ticket_categories.ticket_category_id",
		"promocodes_ticket_categories.promocode_id",
		"promocodes_ticket_categories.created_at",
		"promocodes_ticket_categories.updated_at",
	}
	stmt = "INSERT INTO promocodes_ticket_categories %s VALUES %s"
	stmt = ToNamedInsert(stmt, rows)
	if _, err = tx.NamedExec(stmt, schema.TicketCategories); err != nil {
		return fmt.Errorf("Create: can not insert promocodes ticket categories: %w", err)
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("Create: can not commit tx: %w", err)
	}
	return nil
}

func (p promocodePSQLRepo) GetByID(id uuid.UUID) (*domain.Promocode, error) {
	stmt := "SELECT %s FROM promocodes WHERE id = $1"
	rows := []string{
		"promocodes.id",
		"promocodes.limit_use",
		"promocodes.discount_value",
		"promocodes.discount_percent",
		"promocodes.created_at",
		"promocodes.updated_at",
	}
	stmt = ToSelect(stmt, rows)
	var s promocodeSchema
	if err := p.db.Get(&s, stmt, id); err != nil {
		return nil, fmt.Errorf("GetByID: can not get promocodes from db: %w", err)
	}
	dom := p.adapter.ToDomain(s)

	stmt = "SELECT %s FROM promocodes_ticket_categories WHERE promocode_id = $1"
	rows = []string{
		"promocodes_ticket_categories.id",
		"promocodes_ticket_categories.ticket_category_id",
		"promocodes_ticket_categories.promocode_id",
		"promocodes_ticket_categories.created_at",
		"promocodes_ticket_categories.updated_at",
	}
	stmt = ToSelect(stmt, rows)
	if err := p.db.Select(&s.TicketCategories, stmt, id); err != nil {
		return nil, fmt.Errorf("GetByID: can not get promocodes ticket categories from db: %w", err)
	}

	var tcDom []domain.TicketCategory
	for _, ptc := range s.TicketCategories {
		tc, err := p.ticketCategoryRepo.GetByID(ptc.TicketCategoryID)
		if err != nil {
			return nil, fmt.Errorf("GetByID: can not get ticket category from db: %w", err)
		}
		tcDom = append(tcDom, *tc)
	}
	dom.TicketCategories = tcDom
	return &dom, nil
}
