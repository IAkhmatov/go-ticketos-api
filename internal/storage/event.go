package storage

import (
	"errors"
	"fmt"

	"go-ticketos/internal/domain"

	"github.com/jmoiron/sqlx"
)

type eventPSQLRepo struct {
	db      *sqlx.DB
	adapter EventAdapter
}

var _ domain.EventRepo = (*eventPSQLRepo)(nil)

// nolint: revive
func NewEventRepo(db *sqlx.DB) (*eventPSQLRepo, error) {
	if db == nil {
		return nil, errors.New("NewEventRepo: db is null")
	}
	return &eventPSQLRepo{
		db:      db,
		adapter: NewEventAdapter(),
	}, nil
}

func (r eventPSQLRepo) Create(event domain.Event) error {
	schema := r.adapter.ToSchema(event)
	rows := []string{
		"events.id",
		"events.name",
		"events.description",
		"events.place",
		"events.age_rating",
		"events.start_at",
		"events.end_at",
		"events.created_at",
		"events.updated_at",
	}
	stmt := "INSERT INTO events %s VALUES %s"
	stmt = ToNamedInsert(stmt, rows)
	if _, err := r.db.NamedExec(stmt, schema); err != nil {
		return fmt.Errorf("Create: can not insert row: %w", err)
	}
	return nil
}
