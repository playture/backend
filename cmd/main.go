package main

import (
	"context"
	"github.com/playture/backend/internal/infrastructure/godotenv"
	"github.com/playture/backend/internal/infrastructure/postgresql"
	"github.com/playture/backend/internal/infrastructure/redis"
	"log"
	"log/slog"
	"os"
	"time"
)

func main() {
	env := godotenv.NewEnv()
	logger := initSlogLogger()
	logger.Info("service started")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	pg := postgresql.NewPostgres(env)
	err := pg.Setup(ctx)
	if err != nil {
		log.Fatalf("postgresql error %s\n", err)
	}
	defer pg.Close()

	rdis := redis.NewRedis(env)
	err = rdis.Setup(ctx)
	if err != nil {
		log.Fatalf("redis error %s\n", err)
	}
	defer rdis.Close()
	boot := wireApp(env, logger, pg, rdis)
	boot.Boot()
}

func initSlogLogger() *slog.Logger {
	slogHandlerOptions := &slog.HandlerOptions{
		AddSource: true,
		Level:     slog.LevelDebug,
	}
	logger := slog.New(slog.NewTextHandler(os.Stdout, slogHandlerOptions))
	slog.SetDefault(logger)

	return logger
}
