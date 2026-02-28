package shortenerhandler

import (
	"encoding/json"
	"net/http"

	"github.com/DNA-Z/url-shortener/internal/service"
)

type ShortenerHandler struct {
	Message string `json:"message"`
}

func ShortenerPost(res http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		http.Error(res, "The method is not supported", http.StatusMethodNotAllowed)
		return
	}

	var request ShortenerHandler
	err := json.NewDecoder(req.Body).Decode(&request)
	if err != nil {
		http.Error(res, "Error parsing JSON: "+err.Error(), http.StatusBadRequest)
		return
	}
	defer req.Body.Close()

	if request.Message == "" {
		http.Error(res, "Message field is required", http.StatusBadRequest)
		return
	}

	result := service.ShortenUrl(request.Message)

	res.WriteHeader(http.StatusCreated)
	res.Write([]byte(result))
}
