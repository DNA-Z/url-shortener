package handler

import (
	"context"
	"database/sql"
	"net/http"
	"time"
)

type DBPingHandler struct {
	db            *sql.DB
	serverAddress string
	baseURL       string
}

func NewDBPingHandler(db *sql.DB, serverAddress string, baseURL string) *DBPingHandler {
	return &DBPingHandler{
		db:            db,
		serverAddress: serverAddress,
		baseURL:       baseURL,
	}
}

func (h *DBPingHandler) GetDbPing(res http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		http.Error(res, "Only GET requests are allowed!", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := h.db.PingContext(ctx); err != nil {
		http.Error(res, "Database connection failed", http.StatusInternalServerError)
		return
	}

	res.WriteHeader(http.StatusOK)
}
