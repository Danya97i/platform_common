package transaction

import (
	"context"

	"github.com/jackc/pgx/v4"
	"github.com/pkg/errors"

	"github.com/Danya97i/platform_common/pkg/db"
	"github.com/Danya97i/platform_common/pkg/db/pg"
)

type manager struct {
	db db.Transactor
}

// NewTransactionManager creates new transaction manager
func NewTransactionManager(db db.Transactor) db.TxManager {
	return &manager{
		db: db,
	}
}

func (m *manager) transaction(ctx context.Context, opts pgx.TxOptions, fn db.Handler) (err error) {
	tx, ok := ctx.Value(pg.TxKey).(pgx.Tx)
	if ok {
		return fn(ctx)
	}

	tx, err = m.db.BeginTx(ctx, opts)
	if err != nil {
		return errors.Wrap(err, "can't begin transaction")
	}

	ctx = pg.MakeContextTx(ctx, tx)

	defer func() {
		if r := recover(); r != nil {
			err = errors.Errorf("panic recovered: %v", r)
		}

		if err != nil {
			if rollbackErr := tx.Rollback(ctx); rollbackErr != nil {
				err = errors.Wrapf(err, "rollbackErr: %v", rollbackErr)
			}
			return
		}

		if err == nil {
			err = tx.Commit(ctx)
			if err != nil {
				err = errors.Wrap(err, "tx commit failed")
			}
		}
	}()

	if err = fn(ctx); err != nil {
		err = errors.Wrap(err, "failed executing code inside transaction")
	}

	return err
}

// ReadUncommited – executes code inside transaction with read committed isolation level
func (m *manager) ReadCommited(ctx context.Context, f db.Handler) error {
	txOps := pgx.TxOptions{IsoLevel: pgx.ReadCommitted}
	return m.transaction(ctx, txOps, f)
}
