package uow

import (
	"context"
	"github.com/jackc/pgx/v5"
	"github.com/playture/backend/internal/infrastructure/postgresql"
	"github.com/playture/backend/utils"
	"time"
)

type TransactionFN func(ctx context.Context, tx pgx.Tx) (interface{}, error)

type IUOW interface {
	Do(ctx context.Context, fn TransactionFN, timeout time.Duration) (interface{}, error)
}

type UOW struct {
	pg *postgresql.Postgres
}

func NewUOW(conn *postgresql.Postgres) IUOW {
	return &UOW{
		pg: conn,
	}
}
func (uow *UOW) Do(ctx context.Context, fn TransactionFN, timeout time.Duration) (interface{}, error) {
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	conn, err := uow.pg.PrimaryConn.Acquire(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Release()

	tx, err := conn.Begin(ctx)
	if err != nil {
		return nil, err
	}

	var result interface{}
	result, err = fn(ctx, tx)
	if err != nil {
		if rollbackErr := tx.Rollback(ctx); rollbackErr != nil {
			return nil, utils.WrapError("transaction rollback failed", rollbackErr, err)
		}
		return nil, err
	}

	if commitErr := tx.Commit(ctx); commitErr != nil {
		return nil, utils.WrapError("transaction commit failed", commitErr)
	}

	return result, nil
}
