package handler

import (
	"io"
	"net/http"
	"strings"
)

func (h *URLHandler) ShortenerPost(res http.ResponseWriter, req *http.Request) {
	var baseAddress = h.baseURL

	if !strings.HasSuffix(baseAddress, "/") {
		baseAddress += "/"
	}

	if req.Method != http.MethodPost {
		http.Error(res, "Only POST requests are allowed!", http.StatusMethodNotAllowed)
		return
	}

	body, err := io.ReadAll(req.Body)
	req.Body.Close()

	if err != nil {
		http.Error(res, "Error reading request body: "+err.Error(), http.StatusBadRequest)
		return
	}

	url := string(body)

	if url == "" {
		http.Error(res, "Request url is required", http.StatusBadRequest)
		return
	}

	shortURL, err := h.urlService.Shorten(url)
	if err != nil {
		http.Error(res, "Error shortening URL: "+err.Error(), http.StatusBadRequest)
	}

	result := baseAddress + shortURL

	res.WriteHeader(http.StatusCreated)
	res.Write([]byte(result))
}
