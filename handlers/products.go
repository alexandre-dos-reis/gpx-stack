package handlers

import (
	"github.com/a-h/templ-examples/hello-world/views/pages"
	"github.com/labstack/echo/v4"
)

func (h *Handlers) Products() {
	h.echo.GET("/products", func(c echo.Context) error {
		products, _ := h.repo.FindAllProducts(h.ctx)
		return h.render(c, pages.ProductsPage(pages.ProductsPageProps{Products: products}))
	})
}
