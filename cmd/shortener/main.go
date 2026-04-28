package main

import (
	"context"
	"database/sql"
	"log"
	"net/http"

	"github.com/DNA-Z/url-shortener/internal/config"
	"github.com/DNA-Z/url-shortener/internal/handler"
	"github.com/DNA-Z/url-shortener/internal/infrastructure"
	"github.com/DNA-Z/url-shortener/internal/middleware"
	"github.com/DNA-Z/url-shortener/internal/service"
	"github.com/go-chi/chi/v5"

	"go.uber.org/zap"
)

func main() {
	logger := getLogger()
	configure := config.NewOptions()
	configure.OptionsInit()
	database := getDB(configure)
	consumer, producer := getBroker(getLogger(), configure.FileStoragePath)

	urlService := service.NewURL(consumer, producer)
	urlHandler := handler.NewURLHandler(urlService, configure.ServerAddress, configure.BaseURL)
	pingHandler := handler.NewDBPingHandler(database, configure.ServerAddress, configure.BaseURL)

	middleware.InitLogger(logger)

	r := chi.NewRouter()
	r.Use(middleware.LoggerMiddleware)
	r.Use(middleware.GzipMiddleware)
	r.Get("/ping", pingHandler.GetDbPing)
	r.Get("/{id}", urlHandler.GetByIDGet)
	r.Post("/", urlHandler.ShortenerPost)
	r.Post("/api/shorten", urlHandler.ShortenURLPost)

	log.Printf("Сервер запущен на %s\n", configure.ServerAddress)
	log.Fatal(http.ListenAndServe(configure.ServerAddress, r))
}

func getLogger() *zap.Logger {
	logger, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}
	defer logger.Sync()

	return logger
}

func getDB(cfg *config.Options) *sql.DB {
	ctx := context.Background()
	database, err := infrastructure.DbConnect(ctx, cfg)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	return database
}

func getBroker(log *zap.Logger, fspath string) (*infrastructure.URLConsumer, *infrastructure.URLProducer) {
	consumer, err := infrastructure.NewConsumer(fspath)
	if err != nil {
		log.Fatal("Error creating consumer", zap.Error(err))
	}
	producer, err := infrastructure.NewURLProducer(fspath)
	if err != nil {
		log.Fatal("Error creating producer", zap.Error(err))
	}

	return consumer, producer
}
