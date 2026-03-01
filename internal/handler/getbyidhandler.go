package handler

import (
	"net/http"
)

func (h *UrlHandler) GetByIdGet(res http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		http.Error(res, "Only GET requests are allowed!", http.StatusBadRequest)
	}

	id := req.PathValue("id")

	if id == "" {
		http.Error(res, "ID not provided", http.StatusBadRequest)
	}

	result, err := h.urlService.GetById(id)

	if err != nil {
		http.Error(res, "URL not found", http.StatusBadRequest)
	}

	//http.Redirect(res, req, result, http.StatusTemporaryRedirect)

	res.Header().Set("Location", result)
	res.WriteHeader(http.StatusTemporaryRedirect)
}
