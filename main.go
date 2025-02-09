package main

import (
	"context"

	"github.com/a-h/templ"
	"github.com/a-h/templ-examples/hello-world/database"
	"github.com/a-h/templ-examples/hello-world/database/repository"
	"github.com/a-h/templ-examples/hello-world/handlers"
	"github.com/labstack/echo/v4"
	"go.uber.org/fx"
)

func render(ctx echo.Context, component templ.Component) error {
	return component.Render(ctx.Request().Context(), ctx.Response())
}

func getDB(c echo.Context) *repository.Queries {
	return c.Get("repo").(*repository.Queries)
}

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
	e *echo.Echo,
) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go func() {
				e.Logger.Fatal(e.Start(":3000"))
			}()

			return nil
		},
		OnStop: func(ctx context.Context) error {
			return e.Close()
		},
	})
}
