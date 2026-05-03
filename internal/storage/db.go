package storage

import (
	"context"
	"database/sql"
	"embed"
	"fmt"

	"github.com/DNA-Z/url-shortener/internal/model"
	_ "github.com/lib/pq"
)

//go:embed queries/*.sql
var sqlFiles embed.FS

type DBStorage struct {
	db *sql.DB
}

func NewDBStorage(connectionString string) (*DBStorage, error) {
	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		return nil, fmt.Errorf("sql.Open failed: %w", err)
	}

	if err := db.Ping(); err != nil {
		db.Close()
		return nil, fmt.Errorf("db.Ping failed: %w", err)
	}

	query, err := sqlFiles.ReadFile("queries/create_table.sql")
	if err != nil {
		return nil, fmt.Errorf("failed to read create_table.sql: %w", err)
	}

	_, err = db.ExecContext(context.Background(), string(query))
	if err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to exec create_table.sql: %w", err)
	}

	return &DBStorage{db: db}, nil
}

func (d *DBStorage) Save(url *model.URLDto) error {
	query, err := sqlFiles.ReadFile("queries/insert_url.sql")
	if err != nil {
		return err
	}

	stmt, err := d.db.PrepareContext(context.Background(), string(query))
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(
		url.UUID,
		url.ShortURL,
		url.OriginalURL,
	)
	return err
}

func (d *DBStorage) Saves(urls []model.URLDto) error {
	tx, err := d.db.Begin()
	if err != nil {
		return err
	}
	for i := range urls {
		err := d.Save(&urls[i])
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit()
}

func (d *DBStorage) Get(shortURL string) (string, error) {
	query, err := sqlFiles.ReadFile("queries/get_original_url.sql")
	if err != nil {
		return "", err
	}

	var originalURL string

	stmt, err := d.db.PrepareContext(context.Background(), string(query))
	if err != nil {
		return "", err
	}
	defer stmt.Close()

	err = stmt.QueryRowContext(context.Background(), shortURL).Scan(&originalURL)
	if err != nil {
		return "", err
	}

	return originalURL, nil
}

func (d *DBStorage) LoadAll() (map[string]string, error) {
	query, err := sqlFiles.ReadFile("queries/load_all_urls.sql")
	if err != nil {
		return nil, err
	}

	stmt, err := d.db.PrepareContext(context.Background(), string(query))
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.Query()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	data := make(map[string]string)
	for rows.Next() {
		var short, original string
		if err := rows.Scan(&short, &original); err != nil {
			return nil, err
		}
		data[short] = original
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return data, nil
}
