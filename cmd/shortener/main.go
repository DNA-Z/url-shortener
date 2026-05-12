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
	defer logger.Sync()

	cfg := config.NewOptions()
	cfg.OptionsInit()

	//auth.SetJWTSecretKey(cfg.SecretKey)
	//middleware.InitAuthMiddleware(auth.IsAuthEnabled())

	urlService := getService(cfg)

	urlHandler := handler.NewURLHandler(urlService, cfg.ServerAddress, cfg.BaseURL)
	pingHandler := handler.NewDBPingHandler(cfg)

	middleware.InitLogger(logger)

	r := chi.NewRouter()
	r.Use(middleware.LoggerMiddleware)
	r.Use(middleware.GzipMiddleware)
	//r.Use(middleware.AuthMiddleware)

	r.Get("/ping", pingHandler.GetDbPing)
	r.Get("/{id}", urlHandler.GetByIDGet)
	r.Get("/api/user/urls", urlHandler.GetUserURLs)
	r.Post("/", urlHandler.ShortenerPost)
	r.Post("/api/shorten", urlHandler.ShortenURLPost)
	r.Post("/api/shorten/batch", urlHandler.ShortenBatchPost)

	//r.Group(func(r chi.Router) {
	//	r.Use(middleware.AuthMiddleware)
	//	r.Get("/api/user/urls", urlHandler.GetUserURLs)
	//})

	log.Printf("Сервер запущен на %s\n", cfg.ServerAddress)
	log.Fatal(http.ListenAndServe(cfg.ServerAddress, r))
}

func getLogger() *zap.Logger {
	logger, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}

	return logger
}

func getService(cfg *config.Options) *service.URLStorages {
	urlService, err := service.NewURL(cfg)
	if err != nil {
		log.Fatal("Failed to initialize storage: ", err)
	}
	return urlService
}
