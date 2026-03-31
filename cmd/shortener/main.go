package main

import (
	"log"
	"net/http"

	"github.com/DNA-Z/url-shortener/internal/config"
	"github.com/DNA-Z/url-shortener/internal/handler"
	"github.com/DNA-Z/url-shortener/internal/middleware"
	"github.com/DNA-Z/url-shortener/internal/service"
	"github.com/go-chi/chi/v5"

	"go.uber.org/zap"
)

func main() {
	logger, err := zap.NewDevelopment()

	if err != nil {
		panic(err)
	}
	defer logger.Sync()

	middleware.InitLogger(logger)

	configure := config.NewOptions()
	configure.OptionsInit()

	urlService := service.NewURL()
	urlHandler := handler.NewURLHandler(urlService, configure.ServerAddress, configure.BaseURL)

	r := chi.NewRouter()
	r.Use(middleware.LoggerMiddleware)
	r.Use(middleware.GzipMiddleware)
	r.Get("/{id}", urlHandler.GetByIDGet)
	r.Post("/", urlHandler.ShortenerPost)
	r.Post("/api/shorten", urlHandler.ShortenURLPost)

	log.Printf("Сервер запущен на %s\n", configure.ServerAddress)
	log.Fatal(http.ListenAndServe(configure.ServerAddress, r))
}
