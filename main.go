package main

import (
	"fmt"
	"net/http"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		Hello(Props{name: "Alex", id: "sldkfjl-sdfsdfsdf-fsdfsdf-fdsdsf"}).Render(r.Context(), w)
	})

	port := ":3000"
	fmt.Printf("Listening on http://localhost%s", port)
	http.ListenAndServe(port, nil)
}
