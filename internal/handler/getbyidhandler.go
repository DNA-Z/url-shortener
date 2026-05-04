package handler

import (
	"log"
	"net/http"
)

func (h *URLHandler) GetByIDGet(res http.ResponseWriter, req *http.Request) {
	log.Printf("Get URL: %v by baseURL: %v", req.URL, h.baseURL)

	if req.Method != http.MethodGet {
		http.Error(res, "Only GET requests are allowed!", http.StatusBadRequest)
		return
	}

	id := req.PathValue("id")
	log.Printf("found short URL: %v", id)

	if id == "" {
		http.Error(res, "ID not provided", http.StatusBadRequest)
		return
	}

	result, err := h.urlService.GetByID(id)
	if err != nil {
		http.Error(res, "URL not found", http.StatusNotFound)
		return
	}

	res.Header().Set("Location", result)
	res.WriteHeader(http.StatusTemporaryRedirect)
}
