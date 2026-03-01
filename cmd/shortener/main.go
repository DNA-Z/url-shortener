package main

import (
	"net/http"

	"github.com/DNA-Z/url-shortener/internal/handler"
)

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("POST /", handler.ShortenerPost)
	mux.HandleFunc("GET /", handler.GetByIdGet)

	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		panic(err)
	}
}
