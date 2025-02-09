package handlers

import (
	"github.com/a-h/templ"
	"github.com/a-h/templ-examples/hello-world/database/repository"
	"github.com/a-h/templ-examples/hello-world/views/pages"
	"github.com/labstack/echo/v4"
)

type Handlers interface{}

type handlers struct {
	repo *repository.Queries
	echo *echo.Echo
}

func New(repo *repository.Queries, echo *echo.Echo) Handlers {
	return &handlers{
		repo: repo,
		echo: echo,
	}
}

func render(ctx echo.Context, component templ.Component) error {
	return component.Render(ctx.Request().Context(), ctx.Response())
}

func (h *handlers) home() {
	h.echo.GET("/", func(c echo.Context) error {
		return render(c, pages.HomePage())
	})
}
