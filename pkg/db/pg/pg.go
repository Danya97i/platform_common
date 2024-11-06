package pg

import (
	"context"

	"github.com/georgysavva/scany/pgxscan"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"

	"github.com/Danya97i/platform_common/pkg/db"
)

type key string

// TxKey is a context key for the transaction
const TxKey key = "tx"

type pg struct {
	dbc *pgxpool.Pool
}

// NewDB returns a new DB instance
func NewDB(dbc *pgxpool.Pool) db.DB {
	return &pg{
		dbc: dbc,
	}
}

// ScanOneContext scans a single row from a query
func (p *pg) ScanOneContext(ctx context.Context, dest any, query db.Query, args ...any) error {
	rows, err := p.QueryContext(ctx, query, args...)
	if err != nil {
		return err
	}

	return pgxscan.ScanOne(dest, rows)
}

// ScanAllContext scans all rows from a query
func (p *pg) ScanAllContext(ctx context.Context, dest any, query db.Query, args ...any) error {
	rows, err := p.QueryContext(ctx, query, args...)
	if err != nil {
		return err
	}

	return pgxscan.ScanAll(dest, rows)
}

// QueryContext return all row from a query
func (p *pg) QueryContext(ctx context.Context, query db.Query, args ...any) (pgx.Rows, error) {
	tx, ok := ctx.Value(TxKey).(pgx.Tx)
	if ok {
		return tx.Query(ctx, query.RawQuery, args...)
	}
	return p.dbc.Query(ctx, query.RawQuery, args...)
}

// QueryRowContext return a single row from a query
func (p *pg) QueryRowContext(ctx context.Context, query db.Query, args ...any) pgx.Row {
	tx, ok := ctx.Value(TxKey).(pgx.Tx)
	if ok {
		return tx.QueryRow(ctx, query.RawQuery, args...)
	}
	return p.dbc.QueryRow(ctx, query.RawQuery, args...)
}

// ExecContext execute a query
func (p *pg) ExecContext(ctx context.Context, query db.Query, args ...any) (pgconn.CommandTag, error) {
	tx, ok := ctx.Value(TxKey).(pgx.Tx)
	if ok {
		return tx.Exec(ctx, query.RawQuery, args...)
	}
	return p.dbc.Exec(ctx, query.RawQuery, args...)
}

//

// BeginTx starts a transaction
func (p *pg) BeginTx(ctx context.Context, txOptions pgx.TxOptions) (pgx.Tx, error) {
	return p.dbc.BeginTx(ctx, txOptions)
}

// Ping pings the database
func (p *pg) Ping(ctx context.Context) error {
	return p.dbc.Ping(ctx)
}

// Close closes connection to the database
func (p *pg) Close() {
	p.dbc.Close()
}

// MakeContextTx creates a context with a transaction
func MakeContextTx(ctx context.Context, tx pgx.Tx) context.Context {
	return context.WithValue(ctx, TxKey, tx)
}
