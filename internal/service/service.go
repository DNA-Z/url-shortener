package service

import (
	"log"

	"github.com/DNA-Z/url-shortener/internal/config"
	"github.com/DNA-Z/url-shortener/internal/dto"
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
			log.Printf("Using database storage")
			return newURLService(store)
		}
		log.Printf("Failed to connect to DB: %v\n", err)
	}

	if cfg.FileStoragePath != "" {
		store, err = storage.NewFileStorage(cfg.FileStoragePath)
		if err == nil {
			log.Println("Using file storage")
			return newURLService(store)
		}
		log.Printf("Failed to open file storage: %v\n", err)
	}

	store = storage.NewMemoryStorage()
	log.Println("Using in-memory storage")
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
	log.Printf("Get by id called with param='%s'", shortURL)

	url, err := u.storage.Get(shortURL)
	if err != nil {
		log.Printf("Get URL failed for ID=%s: %v", url, err)
		return "", err
	}

	log.Printf("Result URL: %v", url)
	return url, nil
}

func (u *URL) Batch(request []dto.BatchRequestDto) (response []dto.BatchResponseDto, err error) {
	response = make([]dto.BatchResponseDto, 0, len(request))

	for _, req := range request {
		shortURL, err := u.Shorten(req.OriginalURL)

		log.Printf("shortened from '%+q' -> '%+q'\n", req.OriginalURL, shortURL)

		if err != nil {
			log.Printf("Failed to shorten URL for ID=%s, OriginalURL=%s: %v", req.ID, req.OriginalURL, err)
			continue
		}
		response = append(response, dto.BatchResponseDto{
			ID:       req.ID,
			ShortURL: "http://localhost:8080/" + shortURL,
		})
	}

	return response, err
}

func newURLService(store storage.URLStorage) (*URL, error) {
	urlService := &URL{storage: store}

	if data, err := store.LoadAll(); err == nil {
		_ = data
	}

	return urlService, nil
}
