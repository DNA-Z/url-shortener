package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/DNA-Z/url-shortener/internal/audit"
	"github.com/DNA-Z/url-shortener/internal/auth"
	"github.com/DNA-Z/url-shortener/internal/dto"
	cerrors "github.com/DNA-Z/url-shortener/internal/errors"
	"github.com/DNA-Z/url-shortener/internal/middleware"
	"github.com/google/uuid"
)

// ShortenURLPost обрабатывает POST /api/shorten запросы для сокращения URL.
//
// Формат запроса: application/json с полем "url".
// Формат ответа: application/json с полем "short_url".
//
// Пример запроса:
//
//	POST /api/shorten HTTP/1.1
//	Content-Type: application/json
//
//	{"url": "https://practicum.ru"}
//
// Пример ответа:
//
//	HTTP/1.1 201 Created
//	Content-Type: application/json
//
//	{"short_url": "http://localhost:8080/iGz4syDL"}
//
// Возможные статусы:
//   - 201 Created - URL успешно создан
//   - 409 Conflict - URL уже существует
//   - 400 Bad Request - неверный JSON или пустой URL
//   - 401 Unauthorized - требуется аутентификация
func (h *URLHandler) ShortenURLPost(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST requests are allowed!", http.StatusMethodNotAllowed)
		return
	}
	var baseAddress = h.baseURL
	if !strings.HasSuffix(baseAddress, "/") {
		baseAddress += "/"
	}

	var request dto.URLRequestDto
	var response dto.URLResponseDto

	dec := json.NewDecoder(r.Body)
	if err := dec.Decode(&request); err != nil {
		http.Error(w, "Cannot decode request JSON body", http.StatusBadRequest)
		return
	}
	if request.URL == "" {
		http.Error(w, "URL in JSON is empty", http.StatusBadRequest)
		return
	}
	httpStatus := http.StatusCreated

	userID := uuid.Nil

	if auth.IsAuthEnabled() {
		var err error
		userID, err = middleware.GetUserIDFromContext(r.Context())
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
	}

	shortURL, err := h.urlService.Shorten(userID, request.URL)

	var conflictErr *cerrors.ConflictError
	if err != nil && !errors.As(err, &conflictErr) {
		http.Error(w, "URL shorten error", http.StatusBadRequest)
		return
	}
	if errors.As(err, &conflictErr) {
		httpStatus = http.StatusConflict
		fmt.Println("Status:", conflictErr.Status)
	}

	response.ShortURL = baseAddress + shortURL

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(httpStatus)

	enc := json.NewEncoder(w)
	if err := enc.Encode(response); err != nil {
		http.Error(w, "Error encoding response", http.StatusBadRequest)
		return
	}

	if h.publisher != nil {
		userID, _ = middleware.GetUserIDFromContext(r.Context())

		event := audit.NewEvent(audit.Shorten, userID, request.URL)
		h.publisher.Publish(event)
	}
}
