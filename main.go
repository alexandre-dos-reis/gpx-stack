package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/a-h/templ-examples/hello-world/db/repository"
	"github.com/a-h/templ-examples/hello-world/views"
	"github.com/jackc/pgx/v5"
)

func main() {
	ctx := context.Background()

	conn, err := pgx.Connect(ctx, os.Getenv("DATABASE_URL"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close(ctx)

	repo := repository.New(conn)
	products, _ := repo.FindAllProducts(ctx)

	router := http.NewServeMux()

	router.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		views.HomePage(
			views.HomePageProps{Products: products},
		).Render(r.Context(), w)
	})

	router.HandleFunc("GET /products", func(w http.ResponseWriter, r *http.Request) {
		views.HomePage(
			views.HomePageProps{Products: products},
		).Render(r.Context(), w)
	})

	port := ":3000"

	server := http.Server{
		Addr:    port,
		Handler: router,
	}

	log.Println("Listening on http://localhost" + port)
	server.ListenAndServe()
}
