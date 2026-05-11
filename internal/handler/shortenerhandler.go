package handler

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/DNA-Z/url-shortener/internal/auth"
	cerrors "github.com/DNA-Z/url-shortener/internal/errors"
	"github.com/google/uuid"
)

func (h *URLHandler) ShortenerPost(w http.ResponseWriter, r *http.Request) {
	log.Print("Shortener URL")
	var baseAddress = h.baseURL

	if !strings.HasSuffix(baseAddress, "/") {
		baseAddress += "/"
	}

	if r.Method != http.MethodPost {
		http.Error(w, "Only POST requests are allowed!", http.StatusMethodNotAllowed)
		return
	}

	body, err := io.ReadAll(r.Body)
	r.Body.Close()

	if err != nil {
		http.Error(w, "Error reading request body: "+err.Error(), http.StatusBadRequest)
		return
	}

	url := string(body)

	if url == "" {
		http.Error(w, "Request url is required", http.StatusBadRequest)
		return
	}

	userID := uuid.Nil

	if auth.IsAuthEnabled() {
		if userIDStr, ok := r.Context().Value("userID").(string); ok && userIDStr != "" {
			if parsed, err := uuid.Parse(userIDStr); err == nil {
				userID = parsed
			}
		}
	}

	shortURL, err := h.urlService.Shorten(url, userID)

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
