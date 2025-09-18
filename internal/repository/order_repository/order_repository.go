package orderRepository

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v5"
	"github.com/playture/backend/internal/entity"
)

var (
	ErrOrderNotFound = errors.New("order not found")
)

type Repository interface {
	Create(ctx context.Context, card *entity.Order, tx pgx.Tx) (string, error)
	FindByField(ctx context.Context, field string, value interface{}, tx pgx.Tx) (*entity.Order, error)
	Delete(ctx context.Context, id string, tx pgx.Tx) error
	Update(ctx context.Context, card *entity.Order, tx pgx.Tx) error
}
