package handler

import (
	"github.com/DNA-Z/url-shortener/internal/service"
)

type URLHandler struct {
	urlService    *service.URLStorage
	serverAddress string
	baseURL       string
}

func NewURLHandler(svc *service.URLStorage, serverAddress string, baseURL string) *URLHandler {
	return &URLHandler{
		urlService:    svc,
		serverAddress: serverAddress,
		baseURL:       baseURL,
	}
}
