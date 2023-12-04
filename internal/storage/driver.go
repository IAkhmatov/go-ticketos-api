package storage

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"go-ticketos/internal/config"

	// import db driver.
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
)

const (
	DriverName = "pgx"
)

func ToSelect(stmt string, rows []string) string {
	qargs := make([]string, len(rows))
	for i, k := range rows {
		ss := strings.Split(k, ".")
		var p string
		if len(ss) == 1 {
			p = ss[0]
		} else {
			p = ss[len(ss)-1]
		}
		_ = p
		qargs[i] = fmt.Sprintf(`%s "%s"`, k, k)
	}
	return fmt.Sprintf(stmt, strings.Join(qargs, ", "))
}

func ToNamedInsert(stmt string, rows []string) string {
	args := make([]string, len(rows))
	for i, k := range rows {
		args[i] = strings.Split(k, ".")[1]
	}
	qArgs := make([]string, len(rows))
	for i, k := range rows {
		qArgs[i] = ":" + k
	}
	return fmt.Sprintf(
		stmt,
		fmt.Sprintf("(%s)", strings.Join(args, ", ")),
		fmt.Sprintf("(%s)", strings.Join(qArgs, ", ")),
	)
}

func ToNamedUpdate(stmt string, rows []string) string {
	qargs := make([]string, len(rows))
	for i, k := range rows {
		qargs[i] = strings.Split(k, ".")[1] + "=:" + k
	}
	return fmt.Sprintf(stmt, strings.Join(qargs, ", "))
}

func NewSqlxDB(connectString string) (*sqlx.DB, error) {
	db, err := sqlx.Connect(DriverName, connectString)
	if err != nil {
		return nil, fmt.Errorf("NewTestDB: can not connect to db: %w", err)
	}
	return db, nil
}

func NewSqlxDBFromConfig(cfg *config.Config) (*sqlx.DB, error) {
	db, err := NewSqlxDB(cfg.DBConnectString)
	if err != nil {
		return nil, fmt.Errorf("NewSqlxDBFromConfig: can not create db connect: %w", err)
	}
	return db, err
}

func RecreateSchema(db *sqlx.DB) error {
	stmt := "DROP SCHEMA public CASCADE; CREATE SCHEMA public;"
	if _, err := db.Exec(stmt); err != nil {
		return fmt.Errorf("RecreateSchema: can not recreate schema conn: %w", err)
	}
	return nil
}

func CreateTables(db *sqlx.DB) error {
	path := filepath.Join("../..", "2-create-tables.sql")
	c, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("CreateSchemas: can not read create tables script: %w", err)
	}
	if _, err = db.Exec(string(c)); err != nil {
		return fmt.Errorf("CreateSchemas: can not exec query: %w", err)
	}
	return nil
}

func CleanAllTables(db *sqlx.DB) error {
	stmt := `TRUNCATE TABLE events CASCADE;
			 TRUNCATE TABLE promocodes CASCADE;
			 TRUNCATE TABLE ticket_categories CASCADE;
			 TRUNCATE TABLE promocodes_ticket_categories CASCADE;
			 TRUNCATE TABLE orders CASCADE;
			 TRUNCATE TABLE tickets CASCADE;`
	if _, err := db.Exec(stmt); err != nil {
		return fmt.Errorf("CleanAllTables: can not run sql query: %w", err)
	}
	return nil
}

func CleanOrders(db *sqlx.DB) error {
	stmt := `TRUNCATE TABLE orders CASCADE;
			 TRUNCATE TABLE tickets CASCADE;`
	if _, err := db.Exec(stmt); err != nil {
		return fmt.Errorf("CleanAllTables: can not run sql query: %w", err)
	}
	return nil
}
