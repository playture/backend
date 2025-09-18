package service

import (
	"context"
	"github.com/playture/backend/internal/dto"
	"github.com/playture/backend/internal/entity"
	jobRepository "github.com/playture/backend/internal/repository/job_repository"
	orderRepository "github.com/playture/backend/internal/repository/order_repository"
	"log/slog"
)

type Job interface {
	CreateJob(ctx context.Context, req dto.CreateJobReq) (dto.CreateJobRes, error) // from api
	GetJob(ctx context.Context, id string) (entity.Job, error)                     // from api
	ProcessJob(ctx context.Context, req entity.Job)                                // worker pool
}

type job struct {
	logger    *slog.Logger
	jobRepo   jobRepository.Repository
	orderRepo orderRepository.Repository
}

func NewJob(logger *slog.Logger,
	jobRepo jobRepository.Repository,
	orderRepo orderRepository.Repository,
) Job {
	return &job{
		logger:    logger.With("layer", "servuce"),
		jobRepo:   jobRepo,
		orderRepo: orderRepo,
	}
}

func (j *job) CreateJob(ctx context.Context, req dto.CreateJobReq) (dto.CreateJobRes, error) {
	lg := j.logger.With("method", "CreateJob")
	lg.Info("create job")
	return dto.CreateJobRes{}, nil
}
func (j *job) GetJob(ctx context.Context, id string) (entity.Job, error) {
	lg := j.logger.With("method", "GetJob")
	lg.Info("get job")
	return entity.Job{}, nil
}
func (j *job) ProcessJob(ctx context.Context, req entity.Job) {
	lg := j.logger.With("method", "ProcessJob")
	lg.Info("process job")

}
