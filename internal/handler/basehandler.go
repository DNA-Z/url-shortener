package handler

import (
	"github.com/DNA-Z/url-shortener/internal/service"
)

type URLHandler struct {
	urlService    *service.URLStorages
	serverAddress string
	baseURL       string
}

func NewURLHandler(svc *service.URLStorages, serverAddress string, baseURL string) *URLHandler {
	return &URLHandler{
		urlService:    svc,
		serverAddress: serverAddress,
		baseURL:       baseURL,
	}
}
