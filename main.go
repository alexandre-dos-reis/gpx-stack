package main

import (
	"log"
	"net/http"
)

func main() {
	router := http.NewServeMux()

	router.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		Hello(
			HelloProps{name: "Bob", id: "sldkfjl-sdfsdfsdf-fsdfsdf-fdsdsf"},
		).Render(r.Context(), w)
	})

	server := http.Server{
		Addr:    ":3000",
		Handler: router,
	}

	log.Println("Listening on http://localhost:3000")
	server.ListenAndServe()
}
