package storage

import (
	"context"
	"database/sql"
	"embed"
	"fmt"
	"log"

	"github.com/DNA-Z/url-shortener/internal/model"
	"github.com/google/uuid"
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

func (d *DBStorage) Get(shortURL string) (string, error) {
	query, err := sqlFiles.ReadFile("queries/get_original_url.sql")
	if err != nil {
		log.Printf("failed to read get_original_url.sql: %v", err)
		return "", err
	}

	var originalURL string

	stmt, err := d.db.PrepareContext(context.Background(), string(query))
	if err != nil {
		log.Printf("failed to get() prepare statement: %v", err)
		return "", err
	}
	defer stmt.Close()

	err = stmt.QueryRowContext(context.Background(), shortURL).Scan(&originalURL)
	if err != nil {
		log.Printf("failed to get() execute statement: %v", err)
		return "", err
	}

	log.Printf("originalURL: %s", originalURL)
	return originalURL, nil
}

func (d *DBStorage) LoadAll() (map[string]string, error) {
	query, err := sqlFiles.ReadFile("queries/load_all_urls.sql")
	if err != nil {
		log.Printf("failed to read load_all_urls.sql: %v", err)
		return nil, err
	}

	stmt, err := d.db.PrepareContext(context.Background(), string(query))
	if err != nil {
		log.Printf("failed to load() prepare statement: %v", err)
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.Query()
	if err != nil {
		log.Printf("failed to load() execute statement: %v", err)
		return nil, err
	}
	defer rows.Close()

	data := make(map[string]string)
	for rows.Next() {
		var short, original string
		if err := rows.Scan(&short, &original); err != nil {
			log.Printf("failed to load() scan: %v", err)
			return nil, err
		}
		data[short] = original
	}

	if err = rows.Err(); err != nil {
		log.Printf("failed to load() rows: %v", err)
		return nil, err
	}

	return data, nil
}

//func (d *DBStorage) GetUserURLs(userID uuid.UUID) ([]dto.UserURLsResponseDto, error) {
//	query, err := sqlFiles.ReadFile("queries/get_user_urls.sql")
//	if err != nil {
//		log.Printf("failed to read get_user_urls.sql: %v", err)
//		return nil, err
//	}
//
//	stmt, err := d.db.PrepareContext(context.Background(), string(query))
//	if err != nil {
//		log.Printf("failed to GetUserURLs() prepare statement: %v", err)
//		return nil, err
//	}
//	defer stmt.Close()
//
//	rows, err := stmt.QueryContext(context.Background(), userID)
//	if err != nil {
//		log.Printf("failed to GetUserURLs() execute statement: %v", err)
//		return nil, err
//	}
//	defer rows.Close()
//
//	var result []dto.UserURLsResponseDto
//	for rows.Next() {
//		var short, original string
//		if err := rows.Scan(&short, &original); err != nil {
//			log.Printf("failed to scan user urls: %v", err)
//			return nil, err
//		}
//		result = append(result, dto.UserURLsResponseDto{
//			ShortURL:    short,
//			OriginalURL: original,
//		})
//	}
//
//	if err = rows.Err(); err != nil {
//		log.Printf("error iterating over user urls rows: %v", err)
//		return nil, err
//	}
//
//	return result, nil
//}

func (d *DBStorage) Save(url *model.URL) error {
	query, err := sqlFiles.ReadFile("queries/insert_url.sql")
	if err != nil {
		log.Printf("failed to read insert_url.sql: %v", err)
		return err
	}

	stmt, err := d.db.PrepareContext(context.Background(), string(query))
	if err != nil {
		log.Printf("failed to save() prepare statement: %v", err)
		return err
	}
	defer stmt.Close()

	return d.write(stmt, url)
}

func (d *DBStorage) Saves(urls []model.URL) error {
	tx, err := d.db.Begin()
	if err != nil {
		log.Printf("failed to begin transaction: %v", err)
		return err
	}

	query, err := sqlFiles.ReadFile("queries/insert_url.sql")
	if err != nil {
		tx.Rollback()
		log.Printf("failed to read insert_url.sql: %v", err)
		return err
	}

	stmt, err := tx.PrepareContext(context.Background(), string(query))
	if err != nil {
		tx.Rollback()
		log.Printf("failed to prepare statement in transaction: %v", err)
		return err
	}
	defer stmt.Close()

	for i := range urls {
		err := d.write(stmt, &urls[i])
		if err != nil {
			tx.Rollback()
			log.Printf("failed to save URL in batch: %v", err)
			return err
		}
	}

	return tx.Commit()
}

func (d *DBStorage) write(stmt *sql.Stmt, url *model.URL) error {
	_, err := stmt.Exec(
		url.UUID,
		uuid.Nil,
		url.ShortURL,
		url.OriginalURL,
	)
	if err != nil {
		log.Printf("failed to write() execute statement: %v", err)
		return err
	}

	log.Printf("saved url: %s -> %s", url.ShortURL, url.OriginalURL)
	return nil
}
