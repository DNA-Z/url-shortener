package storage

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func DBConnect(ctx context.Context, connectionString string) (*sql.DB, error) {

	db, err := sql.Open("pgx", connectionString)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
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
