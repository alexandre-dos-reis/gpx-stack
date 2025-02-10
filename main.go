package main

import (
	"context"

	"github.com/a-h/templ-examples/hello-world/database"
	"github.com/a-h/templ-examples/hello-world/database/repository"
	"github.com/a-h/templ-examples/hello-world/handlers"
	"github.com/labstack/echo/v4"
	"go.uber.org/fx"
)

func main() {
	app := fx.New(
		fx.Provide(
			context.Background,
			database.NewPostgresPool,
			repository.New,
			echo.New,
			handlers.New,
		),
		fx.Invoke(invokeServer),
	)

	app.Run()
}

func invokeServer(
	lc fx.Lifecycle,
	h *handlers.Handlers,
) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go func() {
				h.StartServer(":3000")
			}()

			return nil
		},
		OnStop: func(ctx context.Context) error {
			return h.Shutdown()
		},
	})
}
