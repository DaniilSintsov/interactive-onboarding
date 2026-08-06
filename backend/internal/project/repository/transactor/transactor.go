package transactor

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
)

type (
	DB interface {
		Begin(ctx context.Context) (pgx.Tx, error)
	}

	Transactor interface {
		WithTx(ctx context.Context, f func(ctx context.Context) error) (err error)
	}
)

var _ Transactor = (*transactorImpl)(nil)

type transactorImpl struct {
	db DB
}

func NewTransactor(db DB) *transactorImpl {
	return &transactorImpl{
		db: db,
	}
}
func (t *transactorImpl) WithTx(
	ctx context.Context,
	f func(ctx context.Context) error,
) (err error) {
	if _, err = ExtractTx(ctx); err == nil {
		return f(ctx)
	}

	tx, err := t.db.Begin(ctx)

	if err != nil {
		return fmt.Errorf("transactor - begin tx: %w", err)
	}

	defer func() {
		_ = tx.Rollback(ctx)
	}()

	ctx = injectTx(ctx, tx)

	err = f(ctx)

	if err != nil {
		return fmt.Errorf("transactor - execute in tx: %w", err)
	}

	if commitErr := tx.Commit(ctx); commitErr != nil {
		return fmt.Errorf("transactor - commit tx: %w", commitErr)
	}

	return nil
}

type txKey struct{}

var ErrTxNotFound = errors.New("tx not found in context")

func ExtractTx(ctx context.Context) (pgx.Tx, error) {
	tx, ok := ctx.Value(txKey{}).(pgx.Tx)

	if !ok {
		return nil, ErrTxNotFound
	}

	return tx, nil
}

func injectTx(ctx context.Context, tx pgx.Tx) context.Context {
	return context.WithValue(ctx, txKey{}, tx)
}
