// Package main - сервис сокращения URL.
//
// Сервис предоставляет API для создания коротких ссылок, их хранения и перенаправления.
// Поддерживается несколько типов хранилищ: in-memory, файловое и PostgreSQL
package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"net/http/pprof"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/DNA-Z/url-shortener/internal/audit"
	"github.com/DNA-Z/url-shortener/internal/auth"
	"github.com/DNA-Z/url-shortener/internal/config"
	"github.com/DNA-Z/url-shortener/internal/handler"
	"github.com/DNA-Z/url-shortener/internal/middleware"
	"github.com/DNA-Z/url-shortener/internal/service"
	"github.com/DNA-Z/url-shortener/internal/storage"
	"github.com/go-chi/chi/v5"

	"go.uber.org/zap"
	"golang.org/x/crypto/acme/autocert"
)

var (
	buildVersion string
	buildDate    string
	buildCommit  string
)

func main() {
	printBuildInfo()

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

	log.Printf("Сервер запущен на %s\n", cfg.ServerAddress)

	startServer(cfg, r)
}

func startServer(cfg *config.Options, handler http.Handler) {
	server := &http.Server{
		Addr:    cfg.ServerAddress,
		Handler: handler,
	}

	idleConnsClosed := make(chan struct{})
	sigint := make(chan os.Signal, 1)

	signal.Notify(sigint, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	go func() {
		<-sigint
		log.Println("Получен сигнал завершения, начинаем graceful shutdown...")

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		if err := server.Shutdown(ctx); err != nil {
			log.Printf("Ошибка при graceful shutdown: %v", err)
		}

		log.Println("HTTP сервер завершил работу gracefully")

		close(idleConnsClosed)
	}()

	var err error
	if cfg.EnableHTTPS {
		log.Printf("Запуск HTTPS сервера на %s\n", cfg.ServerAddress)
		log.Println("HTTPS включен с автоматическими сертификатами Let's Encrypt")

		manager := &autocert.Manager{
			Cache:  autocert.DirCache("cache-dir"),
			Prompt: autocert.AcceptTOS,
		}

		server.TLSConfig = manager.TLSConfig()
		err = server.ListenAndServeTLS("", "")
	} else {
		log.Printf("Запуск HTTP сервера на %s\n", cfg.ServerAddress)
		err = server.ListenAndServe()
	}

	if err != nil && err != http.ErrServerClosed {
		log.Fatalf("Ошибка запуска HTTP сервера: %v", err)
	}

	<-idleConnsClosed
	log.Println("Сервер полностью остановлен, ресурсы освобождены")
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

func printBuildInfo() {
	version := buildVersion
	if version == "" {
		version = "N/A"
	}

	date := buildDate
	if date == "" {
		date = "N/A"
	}

	commit := buildCommit
	if commit == "" {
		commit = "N/A"
	}

	fmt.Printf("Build version: %s\n", version)
	fmt.Printf("Build date: %s\n", date)
	fmt.Printf("Build commit: %s\n", commit)
}
