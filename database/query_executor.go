package database

import (
	"context"
	"database/sql"
)

type queryExecutor interface {
	Exec(query string, args ...any) (sql.Result, error)
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
	Prepare(query string) (*sql.Stmt, error)
	PrepareContext(ctx context.Context, query string) (*sql.Stmt, error)
	Query(query string, args ...any) (*sql.Rows, error)
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
	QueryRow(query string, args ...any) *sql.Row
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
}

// NewExecutor creates a new executor.
func NewExecutor(db *sql.DB) executor {
	return executor{
		db: db,
	}
}

type executor struct {
	db queryExecutor
}

func (e *executor) DB() queryExecutor {
	return e.db
}

type Option interface {
	Apply(q *executor) error
}

func WithTransaction(tx *sql.Tx) Option {
	return &withTransaction{tx: tx}
}

type withTransaction struct {
	tx *sql.Tx
}

func (o *withTransaction) Apply(q *executor) error {
	q.db = o.tx
	return nil
}
