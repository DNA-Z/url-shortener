package service

import (
	"fmt"

	"github.com/DNA-Z/url-shortener/internal/config"
	"github.com/DNA-Z/url-shortener/internal/model"
	"github.com/DNA-Z/url-shortener/internal/storage"
)

type URL struct {
	storage storage.URLStorage
}

func NewURL(cfg *config.Options) (*URL, error) {
	var store storage.URLStorage
	var err error

	if cfg.ConnectionString != "" {
		store, err = storage.NewDBStorage(cfg.ConnectionString)
		if err == nil {
			fmt.Println("Using database storage")
			return newURLService(store)
		}
		fmt.Printf("Failed to connect to DB: %v\n", err)
	}

	if cfg.FileStoragePath != "" {
		store, err = storage.NewFileStorage(cfg.FileStoragePath)
		if err == nil {
			fmt.Println("Using file storage")
			return newURLService(store)
		}
		fmt.Printf("Failed to open file storage: %v\n", err)
	}

	store = storage.NewMemoryStorage()
	fmt.Println("Using in-memory storage")
	return newURLService(store)
}

func (u *URL) Shorten(originalURL string) (string, error) {

	data, err := u.storage.LoadAll()
	if err != nil {
		return "", err
	}

	for short, long := range data {
		if long == originalURL {
			return short, nil
		}
	}

	newURL := model.NewShortURL(originalURL)

	if err := u.storage.Save(newURL); err != nil {
		return "", err
	}
	return newURL.ShortURL, nil
}

func (u *URL) GetByID(shortURL string) (string, error) {
	if url, ok := u.storage.Get(shortURL); ok {
		return url, nil
	}

	return "", fmt.Errorf("URL not found")
}

func newURLService(store storage.URLStorage) (*URL, error) {
	urlService := &URL{storage: store}

	if data, err := store.LoadAll(); err == nil {
		_ = data
	}

	return urlService, nil
}
