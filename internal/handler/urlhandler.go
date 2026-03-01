package handler

import (
	"github.com/DNA-Z/url-shortener/internal/service"
)

type UrlHandler struct {
	urlService *service.Url
}

func NewUrlHandler(svc *service.Url) *UrlHandler {
	return &UrlHandler{urlService: svc}
}
