package db

import (
	"context"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
)

// Handler is a function that executes in a transaction.
type Handler func(ctx context.Context) error

// TxManager is a transaction manager.
type TxManager interface {
	ReadCommited(ctx context.Context, f Handler) error
}

// Transactor is a transaction beginner.
type Transactor interface {
	BeginTx(ctx context.Context, txOptions pgx.TxOptions) (pgx.Tx, error)
}

// Query is a query with name.
type Query struct {
	Name     string
	RawQuery string
}

// Client is a database client.
type Client interface {
	DB() DB
	Close() error
}

// SQLExecer is a query executor.
type SQLExecer interface {
	QueryExecer
	NamedExecer
}

// QueryExecer is a query executor.
type QueryExecer interface {
	ExecContext(ctx context.Context, query Query, args ...interface{}) (pgconn.CommandTag, error)
	QueryRowContext(ctx context.Context, Query Query, args ...interface{}) pgx.Row
}

// NamedExecer is a query executor.
type NamedExecer interface {
	ScanOneContext(ctx context.Context, dest interface{}, q Query, args ...interface{}) error
	ScanAllContext(ctx context.Context, dest interface{}, q Query, args ...interface{}) error
}

// Pinger is a pinger.
type Pinger interface {
	Ping(ctx context.Context) error
}

// DB is a database.
type DB interface {
	SQLExecer
	NamedExecer
	Transactor
	Pinger
	Close()
}
