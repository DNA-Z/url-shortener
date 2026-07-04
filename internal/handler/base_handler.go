package handler

import (
	"github.com/DNA-Z/url-shortener/internal/audit"
	"github.com/DNA-Z/url-shortener/internal/service"
)

type URLHandler struct {
	urlService    *service.URLStorage
	serverAddress string
	baseURL       string
	publisher     audit.IPublisher
}

func NewURLHandler(svc *service.URLStorage, serverAddress string, baseURL string, publisher audit.IPublisher) *URLHandler {
	return &URLHandler{
		urlService:    svc,
		serverAddress: serverAddress,
		baseURL:       baseURL,
		publisher:     publisher,
	}
}
