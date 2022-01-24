package db

import (
	"context"
	"database/sql"

	"github.com/B3rs/gork/jobs"
)

func NewTxWrapper(db *sql.DB) TxWrapper {
	return TxWrapper{
		db: db,
	}
}

type TxWrapper struct {
	db *sql.DB
}

func (w TxWrapper) WrapTx(ctx context.Context, f func(ctx context.Context, tx *sql.Tx) (interface{}, error)) (interface{}, error) {
	tx, err := w.db.Begin()
	if err != nil {
		return nil, err
	}

	res, err := f(ctx, tx)
	switch err {
	case nil:
	case jobs.ErrJobNotFound:
		_ = tx.Commit()
		return nil, err
	default:
		_ = tx.Rollback()
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return res, nil
}
