//go:build wireinject
// +build wireinject

package main

import (
	"github.com/google/wire"
	"github.com/playture/backend/internal/infrastructure/godotenv"
	"github.com/playture/backend/internal/infrastructure/postgresql"
	"github.com/playture/backend/internal/infrastructure/redis"

	"log/slog"
)

func wireApp(
	env *godotenv.Env,
	logger *slog.Logger,
	postgresql *postgresql.Postgres,
	rdis *redis.Redis,
) *Boot {
	wire.Build(
		wire.NewSet(NewBoot),
	)
	return &Boot{}
}
