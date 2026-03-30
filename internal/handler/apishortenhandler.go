package handler

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/DNA-Z/url-shortener/internal/dto"
)

func (h *URLHandler) ShortenURLPost(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		http.Error(w, "Only POST requests are allowed!", http.StatusMethodNotAllowed)
		return
	}
	var baseAddress = h.baseURL
	if !strings.HasSuffix(baseAddress, "/") {
		baseAddress += "/"
	}

	var request dto.URLRequestDto
	var response dto.URLResponseDto

	dec := json.NewDecoder(req.Body)
	if err := dec.Decode(&request); err != nil {
		http.Error(w, "Cannot decode request JSON body", http.StatusBadRequest)
		return
	}
	if request.URL == "" {
		http.Error(w, "URL in JSON is empty", http.StatusBadRequest)
		return
	}

	shortUrl := h.urlService.Shorten(request.URL)
	response.ShortURL = baseAddress + shortUrl

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	enc := json.NewEncoder(w)
	if err := enc.Encode(response); err != nil {
		http.Error(w, "Error encoding response", http.StatusBadRequest)
		return
	}
}
