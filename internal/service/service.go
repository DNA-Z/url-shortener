package service

import (
	"log"
	"net/http"

	"github.com/DNA-Z/url-shortener/internal/config"
	"github.com/DNA-Z/url-shortener/internal/dto"
	"github.com/DNA-Z/url-shortener/internal/errors"
	"github.com/DNA-Z/url-shortener/internal/model"
	"github.com/DNA-Z/url-shortener/internal/storage"
	"github.com/google/uuid"
)

type URLStorage struct {
	storage storage.URLStorage
	isDB    bool
}

func NewURL(cfg *config.Options) (*URLStorage, error) {
	var store storage.URLStorage
	var err error
	isDB := false

	if cfg.ConnectionString != "" {
		store, err = storage.NewDBStorage(cfg.ConnectionString)
		if err == nil {
			log.Printf("Using database storage")
			isDB = true
			return newURLService(store, isDB)
		}
		log.Printf("Failed to connect to DB: %v\n", err)
	}

	if cfg.FileStoragePath != "" {
		store, err = storage.NewFileStorage(cfg.FileStoragePath)
		if err == nil {
			log.Println("Using file storage")
			return newURLService(store, isDB)
		}
		log.Printf("Failed to open file storage: %v\n", err)
	}

	store = storage.NewMemoryStorage()
	log.Println("Using in-memory storage")
	return newURLService(store, isDB)
}

func (u *URLStorage) Shorten(originalURL string, userID string) (string, error) {

	data, err := u.storage.LoadAll()
	if err != nil {
		return "", err
	}

	for short, long := range data {
		if long == originalURL {
			if u.isDB {
				return short, &errors.ConflictError{
					Status: http.StatusConflict,
					URL:    originalURL,
				}
			}
			return short, nil
		}
	}

	newURL, err := model.NewShortURL(originalURL, userID)
	if err != nil {
		return "", err
	}

	if err := u.storage.Save(newURL); err != nil {
		return "", err
	}
	return newURL.ShortURL, nil
}

func (u *URLStorage) GetByID(shortURL string) (string, error) {
	log.Printf("Get by id called with param='%s'", shortURL)

	url, err := u.storage.Get(shortURL)
	if err != nil {
		log.Printf("Get URL failed for ID=%s: %v", url, err)
		return "", err
	}

	log.Printf("Result URL: %v", url)
	return url, nil
}

func (u *URLStorage) GetUserUrls(userID string) ([]dto.UserURLsResponseDto, error) {
	log.Printf("Get urls by user id called with param='%s'", userID)

	parsedUUID, err := uuid.Parse(userID)
	if err != nil {
		log.Printf("Parse error UUID: %v\n", err)
	}

	urls, err := u.storage.GetUserURLs(parsedUUID)
	if err != nil {
		log.Printf("Get URLs failed for userID=%s: %v", userID, err)
		return nil, err
	}

	return urls, nil
}

func (u *URLStorage) Batch(request []dto.BatchRequestDto, baseAddress string, userID string) (response []dto.BatchResponseDto, err error) {
	response = make([]dto.BatchResponseDto, 0, len(request))

	for _, req := range request {
		shortURL, err := u.Shorten(req.OriginalURL, userID)

		log.Printf("shortened from '%+q' -> '%+q'\n", req.OriginalURL, shortURL)

		if err != nil {
			log.Printf("Failed to shorten URL for ID=%s, OriginalURL=%s: %v", req.ID, req.OriginalURL, err)
			continue
		}
		response = append(response, dto.BatchResponseDto{
			ID:       req.ID,
			ShortURL: baseAddress + shortURL,
		})
	}

	return response, err
}

func newURLService(store storage.URLStorage, isDB bool) (*URLStorage, error) {
	urlService := &URLStorage{storage: store, isDB: isDB}

	if data, err := store.LoadAll(); err == nil {
		_ = data
	}

	return urlService, nil
}
