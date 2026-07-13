package handler

import (
	"context"
	"database/sql"
	"net/http"
	"time"

	"github.com/DNA-Z/url-shortener/internal/config"
)

// DBPingHandler обрабатывает запросы для проверки подключения к базе данных.
type DBPingHandler struct {
	dbConnectionString string
	serverAddress      string
	baseURL            string
	db                 *sql.DB
}

// NewDBPingHandler создает новый экземпляр DBPingHandler.
func NewDBPingHandler(cfg *config.Options, db *sql.DB) *DBPingHandler {
	return &DBPingHandler{
		dbConnectionString: cfg.ConnectionString,
		serverAddress:      cfg.ServerAddress,
		baseURL:            cfg.BaseURL,
		db:                 db,
	}
}

// GetDBPing обрабатывает GET /ping запросы для проверки работоспособности БД.
//
// Пример запроса:
// GET /ping HTTP/1.1
//
// Пример ответа (HTTP 200 OK):
// (пустое тело)
//
// Пример ответа (HTTP 500 Internal Server Error):
// Database connection failed
//
// Возможные статусы:
//   - 200 OK - подключение к БД успешно
//   - 400 Bad Request - неверный метод запроса
//   - 500 Internal Server Error - ошибка подключения к БД
func (h *DBPingHandler) GetDBPing(w http.ResponseWriter, r *http.Request) {
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
