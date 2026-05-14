package handler

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/DNA-Z/url-shortener/internal/auth"
	cerrors "github.com/DNA-Z/url-shortener/internal/errors"
	"github.com/google/uuid"
)

func (h *URLHandler) ShortenerPost(res http.ResponseWriter, req *http.Request) {
	var baseAddress = h.baseURL

	if !strings.HasSuffix(baseAddress, "/") {
		baseAddress += "/"
	}

	if req.Method != http.MethodPost {
		http.Error(res, "Only POST requests are allowed!", http.StatusMethodNotAllowed)
		return
	}

	body, err := io.ReadAll(req.Body)
	req.Body.Close()

	if err != nil {
		http.Error(res, "Error reading request body: "+err.Error(), http.StatusBadRequest)
		return
	}

	url := string(body)

	if url == "" {
		http.Error(res, "Request url is required", http.StatusBadRequest)
		return
	}

	userID := uuid.Nil

	if auth.IsAuthEnabled() {
		if userIDStr, ok := req.Context().Value("userID").(string); ok && userIDStr != "" {
			if parsed, err := uuid.Parse(userIDStr); err == nil {
				userID = parsed
			}
		}
	}

	shortURL, err := h.urlService.Shorten(userID, url)

	httpStatus := http.StatusCreated
	var conflictErr *cerrors.ConflictError
	if err != nil && !errors.As(err, &conflictErr) {
		http.Error(res, "Error shortening URL: "+err.Error(), http.StatusBadRequest)
	}
	if errors.As(err, &conflictErr) {
		httpStatus = http.StatusConflict
		fmt.Println("Status:", conflictErr.Status)
	}

	result := baseAddress + shortURL

	res.WriteHeader(httpStatus)
	res.Write([]byte(result))
}
