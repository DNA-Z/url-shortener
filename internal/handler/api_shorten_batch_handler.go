package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/DNA-Z/url-shortener/internal/auth"
	"github.com/DNA-Z/url-shortener/internal/dto"
	"github.com/DNA-Z/url-shortener/internal/middleware"
	"github.com/google/uuid"
)

// ShortenBatchPost обрабатывает POST /api/shorten/batch запросы для пакетного сокращения URL.
//
// Позволяет сократить несколько URL за один запрос. Все URL обрабатываются
// в рамках одной сессии пользователя.
//
// Формат запроса: application/json с массивом объектов, содержащих correlation_id и original_url.
// Формат ответа: application/json с массивом объектов, содержащих correlation_id и short_url.
//
// Пример запроса:
//
//	POST /api/shorten/batch HTTP/1.1
//	Content-Type: application/json
//	Authorization: Bearer <token>
//
//	[
//	    {
//	        "correlation_id": "1",
//	        "original_url": "https://example.com/very/long/url/1"
//	    },
//	    {
//	        "correlation_id": "2",
//	        "original_url": "https://example.com/very/long/url/2"
//	    }
//	]
//
// Пример ответа:
//
//	HTTP/1.1 201 Created
//	Content-Type: application/json
//
//	[
//	    {
//	        "correlation_id": "1",
//	        "short_url": "http://localhost:8080/abc123"
//	    },
//	    {
//	        "correlation_id": "2",
//	        "short_url": "http://localhost:8080/def456"
//	    }
//	]
//
// Возможные статусы:
//   - 201 Created - все URL успешно созданы
//   - 400 Bad Request - неверный JSON или пустой запрос
//   - 401 Unauthorized - требуется аутентификация
//   - 500 Internal Server Error - ошибка обработки
func (h *URLHandler) ShortenBatchPost(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST requests are allowed!", http.StatusMethodNotAllowed)
		return
	}
	var baseAddress = h.baseURL
	if !strings.HasSuffix(baseAddress, "/") {
		baseAddress += "/"
	}

	var request []dto.BatchRequestDto
	dec := json.NewDecoder(r.Body)
	if err := dec.Decode(&request); err != nil {
		http.Error(w, "Cannot decode request JSON body", http.StatusBadRequest)
		return
	}

	userID := uuid.Nil

	if auth.IsAuthEnabled() {
		var err error
		userID, err = middleware.GetUserIDFromContext(r.Context())
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
	}

	response, err := h.urlService.Batch(userID, request, baseAddress)
	if err != nil {
		http.Error(w, "Batch processing failed", http.StatusInternalServerError)
		return
	}

	for i := range response {
		log.Printf("Batch processing response short URL: %v", response[i].ShortURL)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	enc := json.NewEncoder(w)
	if err := enc.Encode(response); err != nil {
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
		return
	}
}
