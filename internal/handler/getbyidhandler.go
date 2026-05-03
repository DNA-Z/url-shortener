package handler

import (
	"net/http"
	"strings"
)

func (h *URLHandler) GetByIDGet(res http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		http.Error(res, "Only GET requests are allowed!", http.StatusBadRequest)
		return
	}

	id := req.PathValue("id")

	if id == "" {
		http.Error(res, "ID not provided", http.StatusBadRequest)
		return
	}

	result, err := h.urlService.GetByID(id)

	if err != nil {
		http.Error(res, "URL not found", http.StatusNotFound)
		return
	}

	if !strings.HasPrefix(result, "http://") && !strings.HasPrefix(result, "https://") {
		http.Error(res, "Invalid redirect URL: missing scheme", http.StatusInternalServerError)
		return
	}

	res.Header().Set("Location", result)
	res.WriteHeader(http.StatusTemporaryRedirect)
}
