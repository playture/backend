package jobRepository

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v5"
	"github.com/playture/backend/internal/entity"
)

var (
	ErrJobNotFound = errors.New("job not found")
)

type Repository interface {
	Create(ctx context.Context, card *entity.Job, tx pgx.Tx) error
	FindByField(ctx context.Context, field string, value interface{}, tx pgx.Tx) (*entity.Job, error)
	Delete(ctx context.Context, id string, tx pgx.Tx) error
	Update(ctx context.Context, card *entity.Job, tx pgx.Tx) error
}
