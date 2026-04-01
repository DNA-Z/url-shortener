package service

import (
	"fmt"

	"github.com/DNA-Z/url-shortener/internal/infrastructure"
	"github.com/DNA-Z/url-shortener/internal/model"
)

type URL struct {
	consumer infrastructure.URLConsumer
	producer infrastructure.URLProducer
	URLs     map[string]string
}

func NewURL(consumer infrastructure.URLConsumer, producer infrastructure.URLProducer) *URL {
	urlService := &URL{
		consumer: consumer,
		producer: producer,
		URLs:     make(map[string]string),
	}

	urlService.loadAllURLToMap()

	return urlService
}

func (u *URL) Shorten(url string) (string, error) {

	isURLExist, id := u.urlExists(url)

	if isURLExist {
		return id, nil
	}

	newURL := model.NewShortURL(url)

	file, err := model.NewURLFile(newURL.URLID, newURL.LongURL)
	if err != nil {
		return "", err
	}

	err = u.producer.WriteURL(file)
	if err != nil {
		return "", err
	}

	u.URLs[newURL.URLID] = newURL.LongURL

	return newURL.URLID, nil
}

func (u *URL) GetByID(urlID string) (string, error) {
	foundURL, ok := u.URLs[urlID]
	if !ok {
		return "", fmt.Errorf("URL %v not found", urlID)
	}

	return foundURL, nil
}

func (u *URL) loadAllURLToMap() {
	files, err := u.consumer.ReadURL()
	if err != nil {
		return
	}

	for _, url := range files {
		u.URLs[url.ShortURL] = url.OriginalURL
	}
}

func (u *URL) urlExists(url string) (bool, string) {
	for key, value := range u.URLs {
		if value == url {
			return true, key
		}
	}
	return false, ""
}
