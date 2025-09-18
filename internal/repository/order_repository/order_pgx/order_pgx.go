package order_pgx

import (
	"context"
	"github.com/playture/backend/internal/entity"
	"github.com/playture/backend/internal/infrastructure/postgresql"
	orderRepository "github.com/playture/backend/internal/repository/order_repository"
	"log/slog"
)

type OrderPgx struct {
	logger   *slog.Logger
	posrgres *postgresql.Postgres
}

func NewOrderPgx(
	logger *slog.Logger,
	posrgres *postgresql.Postgres,
) *OrderPgx {
	return &OrderPgx{
		logger:   logger.With("layer", "repostitory"),
		posrgres: posrgres,
	}
}

func (o *OrderPgx) Create(ctx context.Context, req *entity.Order) (string, error) {
	lg := o.logger.With("method", "Create")

	var id string
	lg.Info(id)
	return "", orderRepository.ErrOrderNotFound
}
