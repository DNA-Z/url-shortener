package service

import (
	"fmt"

	. "github.com/DNA-Z/url-shortener/internal/model"
)

type Url struct {
	Urls map[string]string
}

func NewUrl() *Url {
	return &Url{
		Urls: make(map[string]string),
	}
}

func (u *Url) Shorten(url string) string {
	newUrl := NewShortUrl(url)

	u.Urls[newUrl.UrlID] = newUrl.LongUrl

	return newUrl.UrlID
}

func (u *Url) GetById(urlID string) (string, error) {
	foundUrl, ok := u.Urls[urlID]
	if !ok {
		return "", fmt.Errorf("URL %v not found", urlID)
	}

	return foundUrl, nil
}
