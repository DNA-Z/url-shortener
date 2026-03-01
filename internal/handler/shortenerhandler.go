package handler

import (
	"io"
	"net/http"

	"github.com/DNA-Z/url-shortener/internal/service"
)

func ShortenerPost(res http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		http.Error(res, "The method is not supported", http.StatusMethodNotAllowed)
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

	result := service.ShortenUrl(url)

	res.WriteHeader(http.StatusCreated)
	res.Write([]byte(result))
}
