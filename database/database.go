package database

import (
	"context"
	"fmt"
	"os"

	"github.com/a-h/templ-examples/hello-world/database/repository"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/fx"
)

func NewPgPool(lc fx.Lifecycle, ctx context.Context) *pgxpool.Pool {
	pool, err := pgxpool.New(context.Background(), os.Getenv("DATABASE_URL"))

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			if err != nil {
				fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
				os.Exit(1)
			}
			return nil
		},
		OnStop: func(ctx context.Context) error {
			pool.Close()
			return nil
		},
	})

	return pool
}

func NewRepositoryPool(pool *pgxpool.Pool) repository.DBTX {
	return pool
}
