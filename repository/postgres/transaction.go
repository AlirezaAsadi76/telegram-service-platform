package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type TransactionProvider interface {
	Execute(ctx context.Context, fn func(tx pgx.Tx) error) error
}

type transactionProvider struct {
	pool *pgxpool.Pool
}

func NewTransactionProvider(pool *pgxpool.Pool) TransactionProvider {
	return &transactionProvider{
		pool: pool,
	}
}

func (p *transactionProvider) Execute(ctx context.Context, fn func(tx pgx.Tx) error) (err error) {
	tx, err := p.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf(
			"begin transaction: %w",
			err,
		)
	}

	defer func() {
		if rollbackErr := tx.Rollback(ctx); err == nil &&
			rollbackErr != nil &&
			!errors.Is(rollbackErr, pgx.ErrTxClosed) {
			err = fmt.Errorf(
				"rollback transaction: %w",
				rollbackErr,
			)
		}
	}()

	if err = fn(tx); err != nil {
		return err
	}

	if err = tx.Commit(ctx); err != nil {
		return fmt.Errorf(
			"commit transaction: %w",
			err,
		)
	}

	return nil
}
