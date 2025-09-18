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
	Create(ctx context.Context, job *entity.Job, tx pgx.Tx) (string, error) // return id
	FindByField(ctx context.Context, field string, value interface{}, tx pgx.Tx) (*entity.Job, error)
	List(ctx context.Context, statuses []entity.JobStatus, orderBy string, ascending bool, limit, page int, tx pgx.Tx) ([]*entity.Job, error)
	Delete(ctx context.Context, id string, tx pgx.Tx) error
	Update(ctx context.Context, job *entity.Job, tx pgx.Tx) error
}
