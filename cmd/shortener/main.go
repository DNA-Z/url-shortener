package main

import (
	"log"
	"net/http"

	"github.com/DNA-Z/url-shortener/internal/handler"
	"github.com/DNA-Z/url-shortener/internal/service"
	"github.com/go-chi/chi/v5"
)

func main() {
	urlService := service.NewURL()
	urlHandler := handler.NewURLHandler(urlService)

	r := chi.NewRouter()
	r.Get("/{id}", urlHandler.GetByIDGet)
	r.Post("/", urlHandler.ShortenerPost)

	log.Fatal(http.ListenAndServe(":8080", r))
}
