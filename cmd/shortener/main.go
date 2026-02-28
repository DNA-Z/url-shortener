package main

import (
	"net/http"

	"github.com/DNA-Z/url-shortener/internal/handler/getbyidhandler"
	"github.com/DNA-Z/url-shortener/internal/handler/shortenerhandler"
)

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("POST /", shortenerhandler.ShortenerPost)
	mux.HandleFunc("GET /", getbyidhandler.GetByIdGet)

	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		panic(err)
	}
}
