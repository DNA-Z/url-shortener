// Package main - сервис сокращения URL.
//
// Сервис предоставляет API для создания коротких ссылок, их хранения и перенаправления.
// Поддерживается несколько типов хранилищ: in-memory, файловое и PostgreSQL
package main

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"net/http/pprof"
	_ "net/http/pprof"

	"github.com/DNA-Z/url-shortener/internal/audit"
	"github.com/DNA-Z/url-shortener/internal/auth"
	"github.com/DNA-Z/url-shortener/internal/config"
	"github.com/DNA-Z/url-shortener/internal/handler"
	"github.com/DNA-Z/url-shortener/internal/middleware"
	"github.com/DNA-Z/url-shortener/internal/service"
	"github.com/DNA-Z/url-shortener/internal/storage"
	"github.com/go-chi/chi/v5"

	"go.uber.org/zap"
)

func main() {
	logger := getLogger()
	defer logger.Sync()

	cfg := config.NewOptions()
	cfg.OptionsInit()

	auth.SetJWTSecretKey(cfg.SecretKey)
	middleware.InitAuthMiddleware(auth.IsAuthEnabled())

	urlService := getService(cfg)
	sqlDB := getDB(cfg)

	auditPublisher := getAuditPublisher(cfg)
	defer auditPublisher.Close()

	urlHandler := handler.NewURLHandler(urlService, cfg.ServerAddress, cfg.BaseURL, auditPublisher)
	pingHandler := handler.NewDBPingHandler(cfg, sqlDB)

	middleware.InitLogger(logger)

	r := chi.NewRouter()
	r.Use(middleware.LoggerMiddleware)
	r.Use(middleware.GzipMiddleware)
	r.Use(middleware.AuditMiddleware(auditPublisher))

	r.Route("/debug/pprof", func(r chi.Router) {
		r.Get("/", pprof.Index)
		r.Get("/cmdline", pprof.Cmdline)
		r.Get("/profile", pprof.Profile)
		r.Get("/symbol", pprof.Symbol)
		r.Get("/trace", pprof.Trace)
		r.Get("/{name}", pprof.Index)
	})

	r.Get("/ping", pingHandler.GetDBPing)

	r.Group(func(r chi.Router) {
		r.Use(middleware.AuthMiddleware)
		r.Get("/{id}", urlHandler.GetByIDGet)
		r.Get("/api/user/urls", urlHandler.GetUserURLs)
		r.Post("/", urlHandler.ShortenerPost)
		r.Post("/api/shorten", urlHandler.ShortenURLPost)
		r.Post("/api/shorten/batch", urlHandler.ShortenBatchPost)
		r.Delete("/api/user/urls", urlHandler.DeleteUserURLs)
	})

	//go func() {
	//	log.Println("Starting pprof on :6060")
	//	if err := http.ListenAndServe(":6060", nil); err != nil {
	//		log.Printf("pprof server error: %v", err)
	//	}
	//}()

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

func getService(cfg *config.Options) *service.URLStorage {
	urlService, err := service.NewURL(cfg)
	if err != nil {
		log.Fatal("Failed to initialize storage: ", err)
	}
	return urlService
}

func getDB(cfg *config.Options) *sql.DB {
	ctx := context.Background()
	database, err := storage.DBConnect(ctx, "postgres://postgres:postgres@postgres:5432/praktikum?sslmode=disable")
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	return database
}

func getAuditPublisher(cfg *config.Options) audit.IPublisher {
	publisher := audit.NewPublisher()

	if cfg.AuditFile != "" {
		fileObserver, err := audit.NewFileObserver(cfg.AuditFile)
		if err != nil {
			log.Printf("Failed to create file observer: %v", err)
		} else {
			publisher.Register(fileObserver)
			log.Printf("Audit file observer registered: %s", cfg.AuditFile)
		}
	}

	if cfg.AuditURL != "" {
		httpObserver := audit.NewHTTPObserver(cfg.AuditURL)
		publisher.Register(httpObserver)
		log.Printf("Audit HTTP observer registered: %s", cfg.AuditURL)
	}

	return publisher
}
