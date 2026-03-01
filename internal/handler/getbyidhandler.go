package handler

import (
	"net/http"
)

func (h *UrlHandler) GetByIdGet(res http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		http.Error(res, "Only GET requests are allowed!", http.StatusMethodNotAllowed)
		return
	}

	id := req.PathValue("id")

	if id == "" {
		http.Error(res, "ID not provided", http.StatusBadRequest)
		return
	}

	result, err := h.urlService.GetById(id)

	if err != nil {
		http.Error(res, "URL not found", http.StatusNotFound)
		return
	}

	http.Redirect(res, req, result, http.StatusTemporaryRedirect)
}
