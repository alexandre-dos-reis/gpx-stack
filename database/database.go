package database

import (
	"context"
	"fmt"
	"os"

	"github.com/a-h/templ-examples/hello-world/database/repository"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/fx"
)

func NewPostgresPool(lc fx.Lifecycle, ctx context.Context) repository.DBTX {
	pool, err := pgxpool.New(ctx, os.Getenv("DATABASE_URL"))

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
