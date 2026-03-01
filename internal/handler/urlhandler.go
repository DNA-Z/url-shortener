package handler

import (
	"github.com/DNA-Z/url-shortener/internal/service"
)

type URLHandler struct {
	urlService *service.URL
}

func NewURLHandler(svc *service.URL) *URLHandler {
	return &URLHandler{urlService: svc}
}
