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
	"net"
	"net/http"
	"net/http/pprof"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/DNA-Z/url-shortener/api/proto"
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
	"google.golang.org/grpc"
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

	log.Printf("Конфигурация:")
	log.Printf("  HTTP адрес: %s", cfg.ServerAddress)
	log.Printf("  gRPC адрес: %s", cfg.GRPCServerAddress)
	log.Printf("  gRPC включен: %v", cfg.EnableGRPC)
	log.Printf("  Base URL: %s", cfg.BaseURL)
	log.Printf("  HTTPS: %v", cfg.EnableHTTPS)

	auth.SetJWTSecretKey(cfg.SecretKey)
	middleware.InitAuthMiddleware(auth.IsAuthEnabled())

	urlService := getService(cfg)
	sqlDB := getDB(cfg)

	auditPublisher := getAuditPublisher(cfg)
	defer auditPublisher.Close()

	urlHandler := handler.NewURLHandler(urlService, cfg, auditPublisher)
	pingHandler := handler.NewDBPingHandler(cfg, sqlDB)

	grpcHandler := handler.NewGRPCHandler(urlService, cfg.BaseURL)

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
		r.Get("/api/internal/stats", urlHandler.StatsGet)
		r.Post("/", urlHandler.ShortenerPost)
		r.Post("/api/shorten", urlHandler.ShortenURLPost)
		r.Post("/api/shorten/batch", urlHandler.ShortenBatchPost)
		r.Delete("/api/user/urls", urlHandler.DeleteUserURLs)
	})

	httpServer := &http.Server{
		Addr:    cfg.ServerAddress,
		Handler: r,
	}

	var grpcServer *grpc.Server
	var grpcListener net.Listener
	if cfg.EnableGRPC {
		grpcServer = grpc.NewServer()
		shortener.RegisterShortenerServiceServer(grpcServer, grpcHandler)

		var err error
		grpcListener, err = net.Listen("tcp", cfg.GRPCServerAddress)
		if err != nil {
			log.Fatalf("Failed to listen on gRPC port %s: %v", cfg.GRPCServerAddress, err)
		}
	}

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		startHTTPServer(httpServer, cfg)
	}()

	if cfg.EnableGRPC {
		wg.Add(1)
		go func() {
			defer wg.Done()
			startGRPCServer(grpcServer, grpcListener)
		}()
	}

	waitForShutdown(httpServer, grpcServer, &wg)
}

func startHTTPServer(server *http.Server, cfg *config.Options) {
	log.Printf("HTTP сервер запущен на %s", cfg.ServerAddress)

	var err error
	if cfg.EnableHTTPS {
		log.Println("HTTPS включен с автоматическими сертификатами Let's Encrypt")
		manager := &autocert.Manager{
			Cache:  autocert.DirCache("cache-dir"),
			Prompt: autocert.AcceptTOS,
		}
		server.TLSConfig = manager.TLSConfig()
		err = server.ListenAndServeTLS("", "")
	} else {
		err = server.ListenAndServe()
	}

	if err != nil && err != http.ErrServerClosed {
		log.Fatalf("Ошибка запуска HTTP сервера: %v", err)
	}
}

func startGRPCServer(server *grpc.Server, listener net.Listener) {
	log.Printf("gRPC сервер запущен на %s", listener.Addr().String())

	if err := server.Serve(listener); err != nil {
		log.Fatalf("Failed to start gRPC server: %v", err)
	}
}

func waitForShutdown(httpServer *http.Server, grpcServer *grpc.Server, wg *sync.WaitGroup) {
	sigint := make(chan os.Signal, 1)
	signal.Notify(sigint, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	<-sigint
	log.Println("Получен сигнал завершения, начинаем graceful shutdown...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Graceful shutdown HTTP сервера
	log.Println("Останавливаем HTTP сервер...")
	if err := httpServer.Shutdown(shutdownCtx); err != nil {
		log.Printf("Ошибка при остановке HTTP сервера: %v", err)
	} else {
		log.Println("HTTP сервер остановлен gracefully")
	}

	// Graceful stop gRPC сервера
	if grpcServer != nil {
		log.Println("Останавливаем gRPC сервер...")

		done := make(chan struct{})
		go func() {
			grpcServer.GracefulStop()
			close(done)
		}()

		select {
		case <-done:
			log.Println("gRPC сервер остановлен gracefully")
		case <-shutdownCtx.Done():
			log.Println("Таймаут при остановке gRPC сервера, принудительное завершение")
			grpcServer.Stop()
		}
	}

	log.Println("Ожидаем завершения всех горутин...")
	wg.Wait()

	log.Println("Сервер завершил работу gracefully")
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
