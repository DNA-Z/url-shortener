package handler

import (
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/DNA-Z/url-shortener/internal/middleware"
	"github.com/google/uuid"
)

// DeleteUserURLs обрабатывает DELETE /api/user/urls запросы для удаления URL пользователя.
// Удаляет (помечает как удаленные) URL, принадлежащие текущему пользователю.
//
// Формат запроса: application/json с массивом коротких идентификаторов.
//
// Пример запроса:
//
//	DELETE /api/user/urls HTTP/1.1
//	Content-Type: application/json
//	Authorization: Bearer <token>
//
//	["abc123", "def456", "ghi789"]
//
// Пример ответа:
//
//	HTTP/1.1 202 Accepted
//
// Возможные статусы:
//   - 202 Accepted - запрос принят на удаление
//   - 400 Bad Request - неверный JSON или пустой список
//   - 401 Unauthorized - требуется аутентификация
//   - 405 Method Not Allowed - неверный метод
func (h *URLHandler) DeleteUserURLs(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Only Delete requests are allowed", http.StatusMethodNotAllowed)
		return
	}

	userID, err := middleware.GetUserIDFromContext(r.Context())
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	var shortURLs []string
	if err := json.Unmarshal(body, &shortURLs); err != nil {
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}
	if len(shortURLs) == 0 {
		http.Error(w, "URL IDs list is empty", http.StatusBadRequest)
		return
	}

	go func(usrID uuid.UUID, urls []string) {
		if err := h.urlService.DeleteUserURLs(usrID, urls); err != nil {
			log.Printf("Failed to delete URLs: %v", err)
		}
	}(userID, shortURLs)

	w.WriteHeader(http.StatusAccepted)
}
