package handler

import (
	"log"
	"net/http"
	"strings"

	"github.com/DNA-Z/url-shortener/internal/audit"
	"github.com/DNA-Z/url-shortener/internal/middleware"
)

// GetByIDGet обрабатывает GET /{id} запросы для перенаправления по короткому URL.
//
// Выполняет HTTP-редирект (307 Temporary Redirect) на оригинальный URL.
//
// Пример запроса:
//
//	GET /abc123 HTTP/1.1
//
// Пример ответа:
//
//	HTTP/1.1 307 Temporary Redirect
//	Location: https://practicum.ru
//
// Возможные статусы:
//   - 307 Temporary Redirect - успешное перенаправление
//   - 400 Bad Request - ID не указан
//   - 404 Not Found - короткий URL не найден
//   - 410 Gone - URL был удален
func (h *URLHandler) GetByIDGet(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Only GET requests are allowed!", http.StatusBadRequest)
		return
	}

	id := r.PathValue("id")
	log.Printf("found short URL: %v", id)

	var baseAddress = h.baseURL
	if !strings.HasSuffix(baseAddress, "/") {
		baseAddress += "/"
	}
	log.Printf("base address: %v", h.baseURL)

	if id == "" {
		http.Error(w, "ID not provided", http.StatusBadRequest)
		return
	}

	url, err := h.urlService.GetOriginURLByShortURL(id)
	if err != nil {
		http.Error(w, "URL not found", http.StatusBadRequest)
		return
	}
	if url.IsDeleted {
		http.Error(w, "URL is gone", http.StatusGone)
	}

	result := url.OriginalURL

	w.Header().Set("Location", result)
	w.WriteHeader(http.StatusTemporaryRedirect)

	if h.publisher != nil {
		userID, err := middleware.GetUserIDFromContext(r.Context())
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		event := audit.NewEvent(audit.Follow, userID, url.OriginalURL)
		h.publisher.Publish(event)
	}
}
