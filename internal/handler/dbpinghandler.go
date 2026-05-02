package handler

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"time"

	"github.com/DNA-Z/url-shortener/internal/config"
	"github.com/DNA-Z/url-shortener/internal/storage"
)

var defaultConnnectionStr string = "postgres://postgres:postgres@postgres:5432/praktikum?sslmode=disable"

type DBPingHandler struct {
	dbConnectionString string
	serverAddress      string
	baseURL            string
}

func NewDBPingHandler(cfg *config.Options) *DBPingHandler {
	return &DBPingHandler{
		dbConnectionString: cfg.ConnectionString,
		serverAddress:      cfg.ServerAddress,
		baseURL:            cfg.BaseURL,
	}
}

func (h *DBPingHandler) GetDbPing(res http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		http.Error(res, "Only GET requests are allowed!", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	db := getDB(h.dbConnectionString)
	if err := db.PingContext(ctx); err != nil {
		http.Error(res, "Database connection failed", http.StatusInternalServerError)
		return
	}

	res.WriteHeader(http.StatusOK)
}

func getDB(connectionString string) *sql.DB {
	if connectionString == "" {
		connectionString = defaultConnnectionStr
	}
	ctx := context.Background()
	database, err := storage.DbConnect(ctx, connectionString)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	return database
}
