package main

import (
	"context"
	"fmt"
	"os"

	"github.com/a-h/templ"
	"github.com/a-h/templ-examples/hello-world/db/repository"
	"github.com/a-h/templ-examples/hello-world/views/pages"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
)

func render(ctx echo.Context, component templ.Component) error {
	return component.Render(ctx.Request().Context(), ctx.Response())
}

func getDB(c echo.Context) *pgxpool.Pool {
	return c.Get("db").(*pgxpool.Pool)
}

func main() {
	ctx := context.Background()

	db, err := pgxpool.New(ctx, os.Getenv("DATABASE_URL"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer db.Close()

	e := echo.New()

	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Set("db", db)
			return next(c)
		}
	})

	e.GET("/", func(c echo.Context) error {
		return render(c, pages.HomePage())
	})

	e.GET("/products", func(c echo.Context) error {
		db := getDB(c)
		repo := repository.New(db)
		products, _ := repo.FindAllProducts(ctx)
		return render(c, pages.ProductsPage(
			pages.ProductsPageProps{Products: products},
		))
	})

	e.Logger.Fatal(e.Start(":3000"))
}
