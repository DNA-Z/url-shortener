package service

import (
	"fmt"

	"github.com/DNA-Z/url-shortener/internal/model"
)

type URL struct {
	Urls map[string]string
}

func NewURL() *URL {
	return &URL{
		Urls: make(map[string]string),
	}
}

func (u *URL) Shorten(url string) string {

	isURLExist, id := u.urlExists(url)

	if isURLExist {
		return id
	}

	newURL := model.NewShortURL(url)

	u.Urls[newURL.URLID] = newURL.LongURL

	return newURL.URLID
}

func (u *URL) GetByID(urlID string) (string, error) {
	foundURL, ok := u.Urls[urlID]
	if !ok {
		return "", fmt.Errorf("URL %v not found", urlID)
	}

	return foundURL, nil
}

func (u *URL) urlExists(url string) (bool, string) {
	for key, value := range u.Urls {
		if value == url {
			return true, key
		}
	}
	return false, ""
}
