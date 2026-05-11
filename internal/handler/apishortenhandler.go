package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/DNA-Z/url-shortener/internal/auth"
	"github.com/DNA-Z/url-shortener/internal/dto"
	cerrors "github.com/DNA-Z/url-shortener/internal/errors"
)

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

	userID := ""

	if auth.IsAuthEnabled() {
		userID, _ = r.Context().Value("userID").(string)
	}

	shortUrl, err := h.urlService.Shorten(request.URL, userID)

	var conflictErr *cerrors.ConflictError
	if err != nil && !errors.As(err, &conflictErr) {
		http.Error(w, "URL shorten error", http.StatusBadRequest)
		return
	}
	if errors.As(err, &conflictErr) {
		httpStatus = http.StatusConflict
		fmt.Println("Status:", conflictErr.Status)
	}

	response.ShortURL = baseAddress + shortUrl

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(httpStatus)

	enc := json.NewEncoder(w)
	if err := enc.Encode(response); err != nil {
		http.Error(w, "Error encoding response", http.StatusBadRequest)
		return
	}
}
