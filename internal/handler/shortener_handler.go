package handler

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/DNA-Z/url-shortener/internal/auth"
	cerrors "github.com/DNA-Z/url-shortener/internal/errors"
	"github.com/DNA-Z/url-shortener/internal/middleware"
	"github.com/google/uuid"
)

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
		userIDStr, ok := middleware.GetUserIDFromContext(r.Context())
		if !ok || userIDStr == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		userID, err = uuid.Parse(userIDStr)
		if err != nil {
			http.Error(w, "Invalid user ID", http.StatusInternalServerError)
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
}
