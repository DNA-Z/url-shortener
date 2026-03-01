package service

import (
	. "github.com/DNA-Z/url-shortener/internal/model"
)

func ShortenUrl(url string) string {

	newUrl := NewUrl(url)

	return newUrl.UrlID
}

// func GetById(urlID string) (string, error) {

// }
