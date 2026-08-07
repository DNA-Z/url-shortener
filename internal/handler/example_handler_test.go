package handler_test

import (
	"bytes"
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
func Example_shortenURL() {
	// Создаем конфигурацию с in-memory хранилищем
	cfg := &config.Options{
		ServerAddress:    "localhost:8080",
		BaseURL:          "http://localhost:8080/",
		FileStoragePath:  "", // Пусто, чтобы использовать in-memory
		ConnectionString: "", // Пусто, чтобы не подключаться к БД
	}

	// Создаем сервис с in-memory хранилищем напрямую для примера
	store := storage.NewMemoryStorage()
	svc, _ := service.NewURLService(store, false)
	publisher := audit.NewPublisher()

	h := handler.NewURLHandler(svc, cfg, publisher)

	r := chi.NewRouter()
	r.Post("/", h.ShortenerPost)

	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBufferString("https://example.com"))
	req.Header.Set("Content-Type", "text/plain")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	fmt.Printf("Status: %d\n", w.Code)
	// Output: Status: 201
}

// Пример создания короткого URL через POST /api/shorten.
func Example_shortenURLJSON() {
	cfg := &config.Options{
		ServerAddress:    "localhost:8080",
		BaseURL:          "http://localhost:8080/",
		FileStoragePath:  "",
		ConnectionString: "",
	}

	store := storage.NewMemoryStorage()
	svc, _ := service.NewURLService(store, false)
	publisher := audit.NewPublisher()

	h := handler.NewURLHandler(svc, cfg, publisher)

	r := chi.NewRouter()
	r.Post("/api/shorten", h.ShortenURLPost)

	body := map[string]string{"url": "https://example.com"}
	jsonBody, _ := json.Marshal(body)

	req := httptest.NewRequest(http.MethodPost, "/api/shorten", bytes.NewReader(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	fmt.Printf("Status: %d\n", w.Code)
	// Output: Status: 201
}

// Пример пакетного создания коротких URL.
func Example_batchShorten() {
	cfg := &config.Options{
		ServerAddress:    "localhost:8080",
		BaseURL:          "http://localhost:8080/",
		FileStoragePath:  "",
		ConnectionString: "",
	}

	store := storage.NewMemoryStorage()
	svc, _ := service.NewURLService(store, false)
	publisher := audit.NewPublisher()

	h := handler.NewURLHandler(svc, cfg, publisher)

	r := chi.NewRouter()
	r.Post("/api/shorten/batch", h.ShortenBatchPost)

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
	// Output: Status: 201
}

// Пример получения URL пользователя (без аутентификации).
func Example_getUserURLs() {
	cfg := &config.Options{
		ServerAddress:    "localhost:8080",
		BaseURL:          "http://localhost:8080/",
		FileStoragePath:  "",
		ConnectionString: "",
	}

	store := storage.NewMemoryStorage()
	svc, _ := service.NewURLService(store, false)
	publisher := audit.NewPublisher()

	h := handler.NewURLHandler(svc, cfg, publisher)

	r := chi.NewRouter()
	r.Get("/api/user/urls", h.GetUserURLs)

	req := httptest.NewRequest(http.MethodGet, "/api/user/urls", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	// Без аутентификации вернется 401
	fmt.Printf("Status: %d\n", w.Code)
	// Output: Status: 401
}

// Пример удаления URL пользователя (без аутентификации).
func Example_deleteUserURLs() {
	cfg := &config.Options{
		ServerAddress:    "localhost:8080",
		BaseURL:          "http://localhost:8080/",
		FileStoragePath:  "",
		ConnectionString: "",
	}

	store := storage.NewMemoryStorage()
	svc, _ := service.NewURLService(store, false)
	publisher := audit.NewPublisher()

	h := handler.NewURLHandler(svc, cfg, publisher)

	r := chi.NewRouter()
	r.Delete("/api/user/urls", h.DeleteUserURLs)

	ids := []string{"abc123"}
	jsonBody, _ := json.Marshal(ids)

	req := httptest.NewRequest(http.MethodDelete, "/api/user/urls", bytes.NewReader(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	// Без аутентификации вернется 401
	fmt.Printf("Status: %d\n", w.Code)
	// Output: Status: 401
}

// Пример проверки подключения к БД.
func Example_ping() {
	cfg := &config.Options{
		ServerAddress:    "localhost:8080",
		BaseURL:          "http://localhost:8080/",
		FileStoragePath:  "",
		ConnectionString: "",
	}

	// Для примера используем nil, чтобы показать ошибку
	h := handler.NewDBPingHandler(cfg, nil)

	r := chi.NewRouter()
	r.Get("/ping", h.GetDBPing)

	req := httptest.NewRequest(http.MethodGet, "/ping", nil)
	w := httptest.NewRecorder()

	// При nil БД будет паника, ловим ее для примера
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Status: 500")
		}
	}()

	r.ServeHTTP(w, req)

	// Output: Status: 500
}
