package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/DNA-Z/url-shortener/internal/dto"
)

func (h *URLHandler) ShortUrlPost(w http.ResponseWriter, req *http.Request) {
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
	var buf bytes.Buffer

	_, err := buf.ReadFrom(req.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = json.Unmarshal(buf.Bytes(), &request)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	shortUrl, err := h.urlService.GetShortUrl(request.URL)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	response.ShortURL = baseAddress + shortUrl
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	resp, err := json.Marshal(response.ShortURL)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(resp)
}
