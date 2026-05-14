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

func (u *URLStorage) Shorten(userID uuid.UUID, originalURL string) (string, error) {

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

	newURL, err := model.NewShortURL(userID, originalURL)

	if err := u.storage.Save(newURL); err != nil {
		return "", err
	}
	return newURL.ShortURL, nil
}

func (u *URLStorage) GetOriginURLByShortURL(shortURL string) (string, error) {
	log.Printf("Get by id called with param='%s'", shortURL)

	url, err := u.storage.Get(shortURL)
	if err != nil {
		log.Printf("Get URL failed for ID=%s: %v", url, err)
		return "", err
	}

	log.Printf("Result URL: %v", url)
	return url, nil
}

func (u *URLStorage) GetUserURLsByUserID(userIDStr string) ([]dto.UserURLsResponseDto, error) {
	log.Printf("Get users URLs by user id called with userID='%s'", userIDStr)

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		log.Print("Invalid user ID format")
		return nil, err
	}

	result, err := u.storage.GetUserURLs(userID)
	if err != nil {
		log.Printf("Get users URLs by user id failed for ID=%s: %v", result, err)
		return nil, err
	}

	log.Printf("Result users URLs by user id: %v", result)
	return result, nil
}

func (u *URLStorage) Batch(userID uuid.UUID, request []dto.BatchRequestDto, baseAddress string) (response []dto.BatchResponseDto, err error) {
	response = make([]dto.BatchResponseDto, 0, len(request))

	for _, req := range request {
		shortURL, err := u.Shorten(userID, req.OriginalURL)

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
