package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/DNA-Z/url-shortener/internal/dto"
)

func (h *URLHandler) ShortenBatchPost(res http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		http.Error(res, "Only POST requests are allowed!", http.StatusMethodNotAllowed)
		return
	}
	var baseAddress = h.baseURL
	if !strings.HasSuffix(baseAddress, "/") {
		baseAddress += "/"
	}

	var request []dto.BatchRequestDto
	dec := json.NewDecoder(req.Body)
	if err := dec.Decode(&request); err != nil {
		http.Error(res, "Cannot decode request JSON body", http.StatusBadRequest)
		return
	}

	response, err := h.urlService.Batch(request)
	if err != nil {
		http.Error(res, "Batch processing failed", http.StatusInternalServerError)
		return
	}

	for i := range response {
		log.Printf("Batch processing response short URL: %v", response[i].ShortURL)
	}

	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(http.StatusCreated)

	enc := json.NewEncoder(res)
	if err := enc.Encode(response); err != nil {
		http.Error(res, "Error encoding response", http.StatusInternalServerError)
		return
	}
}
