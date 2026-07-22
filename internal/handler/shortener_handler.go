package handler

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/DNA-Z/url-shortener/internal/audit"
	"github.com/DNA-Z/url-shortener/internal/auth"
	cerrors "github.com/DNA-Z/url-shortener/internal/errors"
	"github.com/DNA-Z/url-shortener/internal/middleware"
	"github.com/google/uuid"
)

// ShortenerPost обрабатывает POST / запросы для сокращения URL.
//
// Формат запроса: text/plain с URL в теле.
// Формат ответа: текстовый - сокращенный URL.
//
// Пример запроса:
//
//	POST / HTTP/1.1
//	Content-Type: text/plain
//
//	https://practicum.ru
//
// Пример ответа:
//
//	HTTP/1.1 201 Created
//	Content-Type: text/plain
//
//	http://localhost:8080/iGz4syDL
//
// Возможные статусы:
//   - 201 Created - URL успешно создан
//   - 409 Conflict - URL уже существует (возвращается существующий короткий URL)
//   - 400 Bad Request - неверный формат запроса
//   - 401 Unauthorized - требуется аутентификация
func (h *URLHandler) ShortenerPost(w http.ResponseWriter, r *http.Request) {
	var baseAddress = h.baseURL

	if !strings.HasSuffix(baseAddress, "/") {
		baseAddress += "/"
	}

	if r.Method != http.MethodPost {
		http.Error(w, "Only POST requests are allowed!", http.StatusMethodNotAllowed)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body: "+err.Error(), http.StatusBadRequest)
		return
	}
	r.Body.Close()

	url := string(body)

	if url == "" {
		http.Error(w, "Request url is required", http.StatusBadRequest)
		return
	}

	userID := uuid.Nil

	if auth.IsAuthEnabled() {
		userID, err = middleware.GetUserIDFromContext(r.Context())
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
	}

	shortURL, err := h.urlService.Shorten(userID, url)

	httpStatus := http.StatusCreated
	var conflictErr *cerrors.ConflictError
	if err != nil && !errors.As(err, &conflictErr) {
		http.Error(w, "Error shortening URL: "+err.Error(), http.StatusBadRequest)
	}
	if errors.As(err, &conflictErr) {
		httpStatus = http.StatusConflict
		fmt.Println("Status:", conflictErr.Status)
	}

	result := baseAddress + shortURL

	w.WriteHeader(httpStatus)
	w.Write([]byte(result))

	if h.publisher != nil {
		userID, _ = middleware.GetUserIDFromContext(r.Context())

		event := audit.NewEvent(audit.Shorten, userID, url)
		h.publisher.Publish(event)
	}
}
