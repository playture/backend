package main

import (
	"github.com/playture/backend/internal/infrastructure/godotenv"
	"github.com/playture/backend/internal/infrastructure/postgresql"
	"github.com/playture/backend/internal/infrastructure/redis"
	"log/slog"
)

type Boot struct {
	env        *godotenv.Env
	logger     *slog.Logger
	postgresql *postgresql.Postgres
	rdis       *redis.Redis
}

func NewBoot(
	e *godotenv.Env,
	lg *slog.Logger,
	rd *redis.Redis,
	pg *postgresql.Postgres,
) *Boot {
	return &Boot{
		env:        e,
		logger:     lg.With("module", "boot"),
		postgresql: pg,
		rdis:       rd,
	}
}

func (b *Boot) Boot() {
	lg := b.logger.With("method", "Boot")
	lg.Info("it's running")

}
