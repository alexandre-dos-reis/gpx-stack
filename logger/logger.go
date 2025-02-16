package logger

import (
	"context"
	"fmt"
	"os"

	"go.uber.org/fx"
	"go.uber.org/zap"
)

func New(lc fx.Lifecycle) *zap.Logger {
	logger, err := zap.NewDevelopment()

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			if err != nil {
				fmt.Fprintf(os.Stderr, "Unable to load zap logger: %v\n", err)
				os.Exit(1)
			}
			return nil
		},
	})

	return logger
}
