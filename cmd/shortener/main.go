package main

import (
	"context"
	"log"
	"net/http"

	"github.com/DNA-Z/url-shortener/internal/config"
	"github.com/DNA-Z/url-shortener/internal/config/db"
	"github.com/DNA-Z/url-shortener/internal/handler"
	"github.com/DNA-Z/url-shortener/internal/infrastructure"
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

	dbConfig := db.DBConfigInit()
	ctx := context.Background()
	database, err := infrastructure.DbConnect(ctx, dbConfig)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer database.Close()

	configure := config.NewOptions()
	configure.OptionsInit()

	consumer, err := infrastructure.NewConsumer(configure.FileStoragePath)
	if err != nil {
		logger.Fatal("Error creating consumer", zap.Error(err))
	}
	producer, err := infrastructure.NewURLProducer(configure.FileStoragePath)
	if err != nil {
		logger.Fatal("Error creating producer", zap.Error(err))
	}

	urlService := service.NewURL(consumer, producer)
	urlHandler := handler.NewURLHandler(urlService, configure.ServerAddress, configure.BaseURL)
	pingHandler := handler.NewDBPingHandler(database)

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
