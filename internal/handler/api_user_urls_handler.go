package handler

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/DNA-Z/url-shortener/internal/dto"
	"github.com/DNA-Z/url-shortener/internal/middleware"
)

// GetUserURLs обрабатывает GET /api/user/urls запросы для получения всех URL пользователя.
// Возвращает список всех не удаленных URL, созданных текущим пользователем.
//
// Пример запроса:
//
//	GET /api/user/urls HTTP/1.1
//	Authorization: Bearer <token>
//
// Пример ответа (HTTP 200 OK):
//
//	[
//	    {
//	        "short_url": "http://localhost:8080/abc123",
//	        "original_url": "https://example.com/very/long/url/1"
//	    },
//	    {
//	        "short_url": "http://localhost:8080/def456",
//	        "original_url": "https://example.com/very/long/url/2"
//	    }
//	]
//
// Возможные статусы:
//   - 200 OK - успешный ответ со списком URL
//   - 204 No Content - у пользователя нет URL
//   - 401 Unauthorized - требуется аутентификация
//   - 405 Method Not Allowed - неверный метод
func (h *URLHandler) GetUserURLs(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Only GET requests are allowed", http.StatusMethodNotAllowed)
		return
	}

	userID, err := middleware.GetUserIDFromContext(r.Context())
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	urls, err := h.urlService.GetUserURLsByUserID(userID)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	if len(urls) == 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	response := make([]dto.UserURLsResponseDto, 0, len(urls))
	for _, url := range urls {
		response = append(response, dto.UserURLsResponseDto{
			ShortURL:    h.baseURL + url.ShortURL,
			OriginalURL: url.OriginalURL,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		log.Printf("failed to encode response: %v", err)
		return
	}
}
