package handlers

import (
	"github.com/a-h/templ-examples/hello-world/views/pages"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/labstack/echo/v4"
)

func (h *Handlers) Products() {
	h.echo.GET("/products", func(c echo.Context) error {
		products, _ := h.repo.FindAllProducts(h.ctx)
		return h.render(c, pages.ProductsPage(pages.ProductsPageProps{Products: products}))
	})
	h.echo.GET("/products/:slug", func(c echo.Context) error {
		slug := pgtype.Text{String: c.Param("slug"), Valid: true}
		product, _ := h.repo.FindOneProductBySlug(h.ctx, slug)
		return h.render(c, pages.ProductPage(pages.ProductPageProps{Product: product}))
	})
}
