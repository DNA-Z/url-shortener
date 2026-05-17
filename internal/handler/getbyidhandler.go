package handler

import (
	"log"
	"net/http"
	"strings"
)

func (h *URLHandler) GetByIDGet(res http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		http.Error(res, "Only GET requests are allowed!", http.StatusBadRequest)
		return
	}

	id := req.PathValue("id")
	log.Printf("found short URL: %v", id)

	var baseAddress = h.baseURL
	if !strings.HasSuffix(baseAddress, "/") {
		baseAddress += "/"
	}
	log.Printf("base address: %v", h.baseURL)

	if id == "" {
		http.Error(res, "ID not provided", http.StatusBadRequest)
		return
	}

	url, err := h.urlService.GetOriginURLByShortURL(id)
	if err != nil {
		http.Error(res, "URL not found", http.StatusBadRequest)
		return
	}
	if url.IsDeleted {
		http.Error(res, "URL is gone", http.StatusGone)
	}

	res.Header().Set("Location", url.ShortURL)
	res.WriteHeader(http.StatusTemporaryRedirect)
}
