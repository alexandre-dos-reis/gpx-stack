package handlers

import (
	"context"

	"github.com/a-h/templ-examples/hello-world/database/repository"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
)

type Handlers struct {
	echo *echo.Echo
	repo *repository.Queries
	ctx  context.Context
	db   *pgxpool.Pool
}

func New(
	ctx context.Context,
	repo *repository.Queries,
	echo *echo.Echo,
	db *pgxpool.Pool,
) *Handlers {
	return &Handlers{
		echo: echo,
		repo: repo,
		ctx:  ctx,
		db:   db,
	}
}

func (h *Handlers) registerFrontRoutes() {
	h.echo.Static("/assets", "public/assets")
	h.homeHandlers()
	h.productsHandlers()
}

func (h *Handlers) registerAdminRoutes() {
	// h.echo.Static("/*", "public/assets/admin")
	h.adminHandlers()
}

func (h *Handlers) StartServer(address string) {
	h.registerFrontRoutes()
	h.registerAdminRoutes()

	h.echo.Logger.Fatal(h.echo.Start(address))
}

func (h *Handlers) Shutdown() error {
	return h.echo.Close()
}
