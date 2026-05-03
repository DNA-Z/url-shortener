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
	logger := getLogger()
	cfg := config.NewOptions()
	cfg.OptionsInit()
	urlService := getService(cfg)

	urlHandler := handler.NewURLHandler(urlService, cfg.ServerAddress, cfg.BaseURL)
	pingHandler := handler.NewDBPingHandler(cfg)

	middleware.InitLogger(logger)

	r := chi.NewRouter()
	r.Use(middleware.LoggerMiddleware)
	r.Use(middleware.GzipMiddleware)
	r.Get("/ping", pingHandler.GetDbPing)
	r.Get("/{id}", urlHandler.GetByIDGet)
	r.Post("/", urlHandler.ShortenerPost)
	r.Post("/api/shorten", urlHandler.ShortenURLPost)
	r.Post("/api/shorten/batch", urlHandler.ShortenBatchPost)

	log.Printf("Сервер запущен на %s\n", cfg.ServerAddress)
	log.Fatal(http.ListenAndServe(cfg.ServerAddress, r))
}

func getLogger() *zap.Logger {
	logger, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}
	defer logger.Sync()

	return logger
}

func getService(cfg *config.Options) *service.URL {
	urlService, err := service.NewURL(cfg)
	if err != nil {
		log.Fatal("Failed to initialize storage: ", err)
	}
	return urlService
}
