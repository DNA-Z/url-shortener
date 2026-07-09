package handler_test

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"

	"github.com/DNA-Z/url-shortener/internal/audit"
	"github.com/DNA-Z/url-shortener/internal/config"
	"github.com/DNA-Z/url-shortener/internal/handler"
	"github.com/DNA-Z/url-shortener/internal/service"
	"github.com/DNA-Z/url-shortener/internal/storage"
	"github.com/go-chi/chi/v5"
)

// Пример создания короткого URL через POST /.
func ExampleURLHandler_ShortenerPost() {
	// Создаем тестовый сервер
	cfg := &config.Options{
		ServerAddress: "localhost:8080",
		BaseURL:       "http://localhost:8080/",
	}
	store := storage.NewMemoryStorage()
	urlService := &service.URLStorage{Storage: store}
	publisher := audit.NewPublisher()

	handler := handler.NewURLHandler(urlService, cfg.ServerAddress, cfg.BaseURL, publisher)

	r := chi.NewRouter()
	r.Post("/", handler.ShortenerPost)

	// Создаем тестовый запрос
	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBufferString("https://example.com"))
	req.Header.Set("Content-Type", "text/plain")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	fmt.Printf("Status: %d\n", w.Code)
	fmt.Printf("Response: %s\n", w.Body.String())

	// Output:
	// Status: 201
	// Response: http://localhost:8080/
}

// Пример создания короткого URL через POST /api/shorten.
func ExampleURLHandler_ShortenURLPost() {
	cfg := &config.Options{
		ServerAddress: "localhost:8080",
		BaseURL:       "http://localhost:8080/",
	}
	store := storage.NewMemoryStorage()
	urlService := &service.URLStorage{Storage: store}
	publisher := audit.NewPublisher()

	handler := handler.NewURLHandler(urlService, cfg.ServerAddress, cfg.BaseURL, publisher)

	r := chi.NewRouter()
	r.Post("/api/shorten", handler.ShortenURLPost)

	// Создаем JSON запрос
	body := map[string]string{"url": "https://example.com"}
	jsonBody, _ := json.Marshal(body)

	req := httptest.NewRequest(http.MethodPost, "/api/shorten", bytes.NewReader(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	fmt.Printf("Status: %d\n", w.Code)

	// Output:
	// Status: 201
}

// Пример пакетного создания коротких URL.
func ExampleURLHandler_ShortenBatchPost() {
	cfg := &config.Options{
		ServerAddress: "localhost:8080",
		BaseURL:       "http://localhost:8080/",
	}
	store := storage.NewMemoryStorage()
	urlService := &service.URLStorage{Storage: store}
	publisher := audit.NewPublisher()

	handler := handler.NewURLHandler(urlService, cfg.ServerAddress, cfg.BaseURL, publisher)

	r := chi.NewRouter()
	r.Post("/api/shorten/batch", handler.ShortenBatchPost)

	// Создаем batch запрос
	batch := []map[string]string{
		{"correlation_id": "1", "original_url": "https://example.com/1"},
		{"correlation_id": "2", "original_url": "https://example.com/2"},
	}
	jsonBody, _ := json.Marshal(batch)

	req := httptest.NewRequest(http.MethodPost, "/api/shorten/batch", bytes.NewReader(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	fmt.Printf("Status: %d\n", w.Code)

	// Output:
	// Status: 201
}

// Пример получения URL пользователя.
func ExampleURLHandler_GetUserURLs() {
	cfg := &config.Options{
		ServerAddress: "localhost:8080",
		BaseURL:       "http://localhost:8080/",
	}
	store := storage.NewMemoryStorage()
	urlService := &service.URLStorage{Storage: store}
	publisher := audit.NewPublisher()

	handler := handler.NewURLHandler(urlService, cfg.ServerAddress, cfg.BaseURL, publisher)

	r := chi.NewRouter()
	r.Get("/api/user/urls", handler.GetUserURLs)

	req := httptest.NewRequest(http.MethodGet, "/api/user/urls", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	fmt.Printf("Status: %d\n", w.Code)

	// Output:
	// Status: 204
}

// Пример удаления URL пользователя.
func ExampleURLHandler_DeleteUserURLs() {
	cfg := &config.Options{
		ServerAddress: "localhost:8080",
		BaseURL:       "http://localhost:8080/",
	}
	store := storage.NewMemoryStorage()
	urlService := &service.URLStorage{Storage: store}
	publisher := audit.NewPublisher()

	handler := handler.NewURLHandler(urlService, cfg.ServerAddress, cfg.BaseURL, publisher)

	r := chi.NewRouter()
	r.Delete("/api/user/urls", handler.DeleteUserURLs)

	// Создаем запрос с массивом ID
	ids := []string{"abc123", "def456"}
	jsonBody, _ := json.Marshal(ids)

	req := httptest.NewRequest(http.MethodDelete, "/api/user/urls", bytes.NewReader(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	fmt.Printf("Status: %d\n", w.Code)

	// Output:
	// Status: 202
}

// Пример проверки подключения к БД.
func ExampleDBPingHandler_GetDbPing() {
	cfg := &config.Options{
		ConnectionString: "postgres://user:pass@localhost/db",
	}

	// В реальном коде здесь было бы реальное подключение
	// Для примера используем nil (но это вызовет ошибку)
	db := &sql.DB{}
	handler := handler.NewDBPingHandler(cfg, db)

	r := chi.NewRouter()
	r.Get("/ping", handler.GetDbPing)

	req := httptest.NewRequest(http.MethodGet, "/ping", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	fmt.Printf("Status: %d\n", w.Code)

	// Output:
	// Status: 500
}
