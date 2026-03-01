package main

import (
	"net/http"

	"github.com/DNA-Z/url-shortener/internal/handler"
	"github.com/DNA-Z/url-shortener/internal/service"
)

func main() {
	urlService := service.NewUrl()
	urlHandler := handler.NewUrlHandler(urlService)

	mux := http.NewServeMux()

	mux.HandleFunc("POST /", urlHandler.ShortenerPost)
	mux.HandleFunc("GET /{id}", urlHandler.GetByIdGet)

	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		panic(err)
	}
}
