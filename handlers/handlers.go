package handlers

import (
	"context"

	"github.com/a-h/templ"
	"github.com/a-h/templ-examples/hello-world/database/repository"
	"github.com/labstack/echo/v4"
)

type Handlers struct {
	echo *echo.Echo
	repo *repository.Queries
	ctx  context.Context
}

func New(ctx context.Context, repo *repository.Queries, echo *echo.Echo) *Handlers {
	return &Handlers{
		echo: echo,
		repo: repo,
		ctx:  ctx,
	}
}

func (h *Handlers) render(ctx echo.Context, component templ.Component) error {
	return component.Render(ctx.Request().Context(), ctx.Response())
}

func (h *Handlers) registerRoutes() {
	h.homeHandlers()
	h.productsHandlers()
}

func (h *Handlers) StartServer(address string) {
	h.echo.Static("/assets", "public/assets")

	h.registerRoutes()

	h.echo.Logger.Fatal(h.echo.Start(address))
}

func (h *Handlers) Shutdown() error {
	return h.echo.Close()
}
