package postgresql

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/playture/backend/internal/infrastructure/godotenv"
	"github.com/playture/backend/utils"
)

type Postgres struct {
	PrimaryConn *pgxpool.Pool
	Env         *godotenv.Env
}

func NewPostgres(env *godotenv.Env) *Postgres {
	return &Postgres{
		Env: env,
	}
}

func (p *Postgres) Setup(ctx context.Context) error {
	if p.PrimaryConn != nil {
		p.PrimaryConn.Close()
	}
	var err error

	p.PrimaryConn, err = createConnection(ctx, p.Env.DatabaseURL, 10, 5)
	if err != nil {
		return utils.WrapError("failed to setup primary connection: ", err)
	}

	return nil
}

func (p *Postgres) HealthCheck(ctx context.Context) error {
	if p.PrimaryConn == nil {
		return errors.New("one or both PostgreSQL connection pools are not initialized")
	}

	query := `SELECT 1;`

	row := p.PrimaryConn.QueryRow(ctx, query)
	var result int
	if err := row.Scan(&result); err != nil {
		return errors.New("health check query failed on primary connection: " + err.Error())
	}
	if result != 1 {
		return errors.New("unexpected result from health check query on primary connection")
	}

	return nil
}

func (p *Postgres) Close() error {
	if p.PrimaryConn != nil {
		p.PrimaryConn.Close()
	}
	return nil
}

func createConnection(ctx context.Context, connString string, maxConns, minConns int) (*pgxpool.Pool, error) {
	connConfig, err := pgxpool.ParseConfig(connString)
	if err != nil {
		return nil, utils.WrapError("failed to parse connection string: ", err)
	}
	connConfig.MaxConns = int32(maxConns)
	connConfig.MinConns = int32(minConns)

	pool, err := pgxpool.NewWithConfig(ctx, connConfig)
	if err != nil {
		return nil, utils.WrapError("failed to create connection pool: ", err)
	}

	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, utils.WrapError("failed to ping PostgresSQL: ", err)
	}

	return pool, nil
}
