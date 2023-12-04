package storage

import (
	"database/sql"
	"errors"
	"fmt"

	"go-ticketos/internal/domain"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type orderPSQLRepo struct {
	db           *sqlx.DB
	orderAdapter OrderAdapter
}

var _ domain.OrderRepo = (*orderPSQLRepo)(nil)

// nolint: revive
func NewOrderRepo(db *sqlx.DB) (*orderPSQLRepo, error) {
	if db == nil {
		return nil, fmt.Errorf("NewEventRepo: db is null")
	}
	return &orderPSQLRepo{
		db:           db,
		orderAdapter: NewOrderAdapter(),
	}, nil
}

func (o orderPSQLRepo) GetByID(id uuid.UUID) (*domain.Order, error) {
	stmt := "SELECT %s FROM orders WHERE id = $1"
	rows := []string{
		"orders.id",
		"orders.name",
		"orders.email",
		"orders.phone",
		"orders.status",
		"orders.payment_id",
		"orders.payment_url",
		"orders.created_at",
		"orders.updated_at",
	}
	stmt = ToSelect(stmt, rows)
	var orderS orderSchema
	err := o.db.Get(&orderS, stmt, id)
	if err != nil {
		return nil, fmt.Errorf("GetByID: can not get orders from db: %w", err)
	}

	stmt = `SELECT %s 
			FROM tickets 
			JOIN ticket_categories ON tickets.ticket_category_id = ticket_categories.id
		    WHERE order_id = $1`
	rows = []string{
		"tickets.id",
		"tickets.promocode_id",
		"tickets.ticket_category_id",
		"tickets.full_price",
		"tickets.buy_price",
		"tickets.created_at",
		"tickets.updated_at",

		"ticket_categories.id",
		"ticket_categories.event_id",
		"ticket_categories.price",
		"ticket_categories.name",
		"ticket_categories.description",
		"ticket_categories.created_at",
		"ticket_categories.updated_at",
	}
	stmt = ToSelect(stmt, rows)
	var ticketsS []ticketSchema
	err = o.db.Select(&ticketsS, stmt, orderS.ID)
	if err != nil {
		return nil, fmt.Errorf("GetByID: can not get tickets from db: %w", err)
	}
	orderS.Tickets = ticketsS

	dom, err := o.orderAdapter.ToDomain(orderS)
	if err != nil {
		return nil, fmt.Errorf("GetByID: can not adapt order to domain model: %w", err)
	}
	return dom, nil
}

func (o orderPSQLRepo) Create(order domain.Order) error {
	tx, err := o.db.Beginx()
	if err != nil {
		return fmt.Errorf("Create: can not create tx: %w", err)
	}
	defer func() {
		err = tx.Rollback()
		if err != nil && !errors.Is(err, sql.ErrTxDone) {
			msg := fmt.Sprintf("Create: can not rollback tx: %s", err)
			// nolint: forbidigo
			fmt.Println(msg)
		}
	}()

	schema := o.orderAdapter.ToSchema(order)
	rows := []string{
		"orders.id",
		"orders.name",
		"orders.email",
		"orders.phone",
		"orders.status",
		"orders.payment_id",
		"orders.payment_url",
		"orders.created_at",
		"orders.updated_at",
	}
	stmt := "INSERT INTO orders %s VALUES %s"
	stmt = ToNamedInsert(stmt, rows)
	_, err = tx.NamedExec(stmt, schema)
	if err != nil {
		return fmt.Errorf("Create: can not insert row: %w", err)
	}

	err = o.createTickets(tx, schema.Tickets)
	if err != nil {
		return fmt.Errorf("Create: can not create tickets: %w", err)
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("Create: can not commit tx: %w", err)
	}
	return nil
}

func (o orderPSQLRepo) Update(order domain.Order) error {
	stmt := "UPDATE orders SET %s WHERE id=:orders.id"
	rows := []string{
		"orders.status",
		"orders.payment_id",
		"orders.payment_url",
		"orders.updated_at",
	}
	stmt = ToNamedUpdate(stmt, rows)
	schema := o.orderAdapter.ToSchema(order)
	res, err := o.db.NamedExec(stmt, schema)
	if err != nil {
		return fmt.Errorf("Update: can not update row: %w", err)
	}
	ra, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("Update: can not get rows affected: %w", err)
	}
	if ra == 0 {
		return fmt.Errorf("Update: rows affected equal 0")
	}
	return nil
}

func (o orderPSQLRepo) createTickets(tx *sqlx.Tx, schemas []ticketSchema) error {
	stmt := "INSERT INTO tickets %s VALUES %s"
	rows := []string{
		"tickets.id",
		"tickets.order_id",
		"tickets.ticket_category_id",
		"tickets.promocode_id",
		"tickets.full_price",
		"tickets.buy_price",
		"tickets.created_at",
		"tickets.updated_at",
	}
	stmt = ToNamedInsert(stmt, rows)
	_, err := tx.NamedExec(stmt, schemas)
	if err != nil {
		return fmt.Errorf("createTickets: can not insert row: %w", err)
	}
	return nil
}
