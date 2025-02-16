package database

import (
	"context"
	"fmt"
	"os"

	"github.com/a-h/templ-examples/hello-world/database/repository"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/tracelog"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Custom pgx logger using Zap
type PgxLogger struct {
	logger *zap.Logger
}

// Implements tracelog.Logger
func (logger *PgxLogger) Log(
	ctx context.Context,
	level tracelog.LogLevel,
	msg string,
	data map[string]any,
) {
	// Map tracelog.LogLevel to Zap log levels
	var zapLevel zapcore.Level
	switch level {
	case tracelog.LogLevelTrace, tracelog.LogLevelDebug:
		zapLevel = zap.DebugLevel
	case tracelog.LogLevelInfo:
		zapLevel = zap.InfoLevel
	case tracelog.LogLevelWarn:
		zapLevel = zap.WarnLevel
	case tracelog.LogLevelError:
		zapLevel = zap.ErrorLevel
	default:
		zapLevel = zap.InfoLevel
	}

	// Convert data map into Zap fields
	fields := []zap.Field{}
	for k, v := range data {
		fields = append(fields, zap.Any(k, v))
	}

	// Log with Zap
	logger.logger.Log(zapLevel, msg, fields...)
}

func NewPgPool(lc fx.Lifecycle, ctx context.Context, logger *zap.Logger) *pgxpool.Pool {
	config, _ := pgxpool.ParseConfig(os.Getenv("DB_URL"))
	config.ConnConfig.Tracer = &tracelog.TraceLog{
		Logger:   &PgxLogger{logger: logger},
		LogLevel: tracelog.LogLevelTrace, // Ensure we log everything
		Config:   &tracelog.TraceLogConfig{},
	}

	pool, err := pgxpool.NewWithConfig(
		ctx,
		config,
	)

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

	pool.Config().ConnConfig.Tracer = &tracelog.TraceLog{
		Logger:   &PgxLogger{logger: logger},
		LogLevel: tracelog.LogLevelDebug,
		Config:   &tracelog.TraceLogConfig{},
	}

	return pool
}

func NewRepositoryPool(pool *pgxpool.Pool) repository.DBTX {
	return pool
}
