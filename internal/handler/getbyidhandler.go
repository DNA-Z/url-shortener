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

	normalizeRequestURL(req)
	log.Printf("Normalized URL: scheme=%s host=%s path=%s",
		req.URL.Scheme, req.URL.Host, req.URL.Path)

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

func normalizeRequestURL(req *http.Request) {
	if req.URL.Scheme == "" {
		req.URL.Scheme = "http"
	}

	if req.URL.Host == "" {
		req.URL.Host = "localhost:8080"
	}

	if req.URL.Path == "" || !strings.HasPrefix(req.URL.Path, "/") {
		if req.URL.Path != "" {
			req.URL.Path = "/" + req.URL.Path
		} else {
			req.URL.Path = "/"
		}
	}
}
