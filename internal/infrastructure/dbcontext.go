package infrastructure

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/DNA-Z/url-shortener/internal/config/db"
	_ "github.com/jackc/pgx/v5/stdlib"
)

func DbConnect(ctx context.Context, cfg *db.DBConfig) (*sql.DB, error) {
	strConn := cfg.GetConnectionString()

	db, err := sql.Open("pgx", strConn)
	if err != nil {
		log.Fatal("Failed to connect to database: %v", err)
	}

	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(time.Hour)

	ctxWithTimeout, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err := db.PingContext(ctxWithTimeout); err != nil {
		db.Close()
		return nil, fmt.Errorf("Failed to ping database: %w", err)
	}

	log.Println("Database connected successfully")
	return db, nil
}
