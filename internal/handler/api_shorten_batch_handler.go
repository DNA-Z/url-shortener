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

	var userID uuid.UUID

	if auth.IsAuthEnabled() {
		userID, ok := middleware.GetUserIDFromContext(r.Context())
		if !ok || userID == "" {
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
