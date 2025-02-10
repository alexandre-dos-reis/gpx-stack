package handlers

import (
	"github.com/a-h/templ-examples/hello-world/views/pages"
	"github.com/labstack/echo/v4"
)

func (h *Handlers) homeHandlers() {
	h.echo.GET("/", func(c echo.Context) error {
		return h.render(c, pages.HomePage())
	})
}
