package pg

import (
	"context"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/pkg/errors"

	"github.com/Danya97i/chat-server/internal/client/db"
)

type pgClient struct {
	masterDBC db.DB
}

// NewPGClient creates a new PG client.
func NewPGClient(ctx context.Context, dsn string) (db.Client, error) {
	dbc, err := pgxpool.Connect(ctx, dsn)
	if err != nil {
		return nil, errors.Errorf("failed to connect to db: %v", err)
	}

	return &pgClient{
		masterDBC: &pg{dbc},
	}, nil
}

// DB returns the underlying DB connection.
func (c *pgClient) DB() db.DB {
	return c.masterDBC
}

// Close closes the underlying DB connection.
func (c *pgClient) Close() error {
	if c.masterDBC != nil {
		c.masterDBC.Close()
	}

	return nil
}
