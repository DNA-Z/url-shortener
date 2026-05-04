package handler

import (
	"log"
	"net/http"
	"strings"
)

func (h *URLHandler) GetByIDGet(res http.ResponseWriter, req *http.Request) {
	log.Printf("Get URL: %v by baseURL: %v", req.URL, h.baseURL)

	if req.Method != http.MethodGet {
		http.Error(res, "Only GET requests are allowed!", http.StatusBadRequest)
		return
	}

	if h.normalizeAndRedirectIfNeeded(res, req) {
		return
	}

	var baseAddress = h.baseURL
	if !strings.HasSuffix(baseAddress, "/") {
		baseAddress += "/"
	}
	log.Printf("base address: %v", baseAddress)

	id := req.PathValue("id")
	log.Printf("found short URL: %v", id)

	if id == "" {
		http.Error(res, "ID not provided", http.StatusBadRequest)
		return
	}

	result, err := h.urlService.GetByID(id)

	if err != nil {
		http.Error(res, "URL not found", http.StatusBadRequest)
		return
	}

	res.Header().Set("Location", result)
	res.WriteHeader(http.StatusTemporaryRedirect)
}

func (h *URLHandler) normalizeAndRedirectIfNeeded(res http.ResponseWriter, req *http.Request) bool {
	path := req.URL.Path
	id := req.PathValue("id")

	if id == "" || strings.HasPrefix(path, "http://") || strings.HasPrefix(path, "https://") {
		return false
	}

	if !strings.Contains(path, "://") {
		fullURL := "http://localhost:8080" + path
		log.Printf("Normalizing URL: %s -> %s", path, fullURL)

		res.Header().Set("Location", fullURL)
		res.WriteHeader(http.StatusTemporaryRedirect)
		return true
	}

	return false
}
