package handlers

import (
	"net/http"

	Ra "github.com/a-h/templ-examples/hello-world/handlers/ReactAdmin"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func (h *Handlers) adminHandlers() {
	// h.echo.Use(middleware.Logger())
	// useless if admin is served by the same server...
	h.echo.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{
			http.MethodGet,
			http.MethodHead,
			http.MethodPut,
			http.MethodPatch,
			http.MethodPost,
			http.MethodDelete,
		},
	}))
	h.echo.Debug = true

	g := h.echo.Group("/admin")

	Ra.AllowedTables = map[string]struct{ ColumnsAllowed []string }{
		"products": {ColumnsAllowed: []string{"*"}},
	}

	g.PUT("/:resource/:id", func(c echo.Context) error {
		return Ra.UpdateHandler(c, h.db, h.ctx)
	})
	g.GET("/:resource/:id", func(c echo.Context) error {
		return Ra.GetOneHandler(c, h.db, h.ctx)
	})
	g.GET("/:resource", func(c echo.Context) error {
		return Ra.GetListHandler(c, h.db, h.ctx)
	})
}
