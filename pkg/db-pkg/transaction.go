package dbpkg

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type transactionManager struct {
	pool   *pgxpool.Pool
	logger *slog.Logger
}

func NewTransactionManager(pool *pgxpool.Pool, logger *slog.Logger) TransactionManager {
	if logger == nil {
		logger = slog.Default()
	}
	return &transactionManager{
		pool:   pool,
		logger: logger,
	}
}

func GetTransactionManager() TransactionManager {
	return NewTransactionManager(DB, slog.Default())
}

func (tm *transactionManager) WithTransaction(ctx context.Context, fn func(ctx context.Context) error) error {
	return tm.withTx(ctx, pgx.TxOptions{}, fn)
}

func (tm *transactionManager) WithReadOnlyTransaction(ctx context.Context, fn func(ctx context.Context) error) error {
	return tm.withTx(ctx, pgx.TxOptions{AccessMode: pgx.ReadOnly}, fn)
}

func (tm *transactionManager) withTx(ctx context.Context, opts pgx.TxOptions, fn func(ctx context.Context) error) error {
	tx, err := tm.pool.BeginTx(ctx, opts)
	if err != nil {
		tm.logger.ErrorContext(ctx, "failed to begin transaction", "error", err)
		return fmt.Errorf("begin transaction: %w", err)
	}

	defer func() {
		if p := recover(); p != nil {
			if rbErr := tx.Rollback(ctx); rbErr != nil {
				tm.logger.ErrorContext(ctx, "failed to rollback transaction on panic", "error", rbErr, "panic", p)
			}
			panic(p)
		}
	}()

	// Inject transaction into context
	txCtx := context.WithValue(ctx, "tx", tx)

	if err := fn(txCtx); err != nil {
		if rbErr := tx.Rollback(ctx); rbErr != nil {
			tm.logger.ErrorContext(ctx, "failed to rollback transaction", "rollback_error", rbErr, "original_error", err)
			return fmt.Errorf("rollback failed: %v, original error: %w", rbErr, err)
		}
		tm.logger.DebugContext(ctx, "transaction rolled back due to error", "error", err)
		return err
	}

	if err := tx.Commit(ctx); err != nil {
		tm.logger.ErrorContext(ctx, "failed to commit transaction", "error", err)
		return fmt.Errorf("commit failed: %w", err)
	}

	tm.logger.DebugContext(ctx, "transaction committed successfully")
	return nil
}

func GetDBFromContext(ctx context.Context) DBTX {
	if tx, ok := ctx.Value("tx").(pgx.Tx); ok {
		return tx
	}
	return DB
}
