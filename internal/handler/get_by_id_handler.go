package handler

import (
	"log"
	"net/http"
	"strings"
)

func (h *URLHandler) GetByIDGet(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Only GET requests are allowed!", http.StatusBadRequest)
		return
	}

	id := r.PathValue("id")
	log.Printf("found short URL: %v", id)

	var baseAddress = h.baseURL
	if !strings.HasSuffix(baseAddress, "/") {
		baseAddress += "/"
	}
	log.Printf("base address: %v", h.baseURL)

	if id == "" {
		http.Error(w, "ID not provided", http.StatusBadRequest)
		return
	}

	url, err := h.urlService.GetOriginURLByShortURL(id)
	if err != nil {
		http.Error(w, "URL not found", http.StatusBadRequest)
		return
	}
	if url.IsDeleted {
		http.Error(w, "URL is gone", http.StatusGone)
	}

	var result string
	result = url.OriginalUrl

	w.Header().Set("Location", result)
	w.WriteHeader(http.StatusTemporaryRedirect)
}
