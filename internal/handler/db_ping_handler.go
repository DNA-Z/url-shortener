package handler

import (
	"context"
	"database/sql"
	"net/http"
	"time"

	"github.com/DNA-Z/url-shortener/internal/config"
)

type DBPingHandler struct {
	dbConnectionString string
	serverAddress      string
	baseURL            string
	db                 *sql.DB
}

func NewDBPingHandler(cfg *config.Options, db *sql.DB) *DBPingHandler {
	return &DBPingHandler{
		dbConnectionString: cfg.ConnectionString,
		serverAddress:      cfg.ServerAddress,
		baseURL:            cfg.BaseURL,
		db:                 db,
	}
}

func (h *DBPingHandler) GetDbPing(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Only GET requests are allowed!", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := h.db.PingContext(ctx); err != nil {
		http.Error(w, "Database connection failed", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
