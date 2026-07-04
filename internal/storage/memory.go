package storage

import (
	"fmt"
	"sync"

	"github.com/DNA-Z/url-shortener/internal/dto"
	"github.com/DNA-Z/url-shortener/internal/model"
	"github.com/google/uuid"
)

type MemoryStorage struct {
	data []model.URL
	mu   sync.RWMutex
}

func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{
		data: make([]model.URL, 0),
	}
}

func (m *MemoryStorage) Save(url *model.URL) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.data = append(m.data, model.URL{
		UserID:      url.UserID,
		ShortURL:    url.ShortURL,
		OriginalURL: url.OriginalURL,
		IsDeleted:   false,
	})

	return nil
}

func (m *MemoryStorage) Saves(urls []model.URL) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	for _, url := range urls {
		m.data = append(m.data, model.URL{
			UserID:      url.UserID,
			ShortURL:    url.ShortURL,
			OriginalURL: url.OriginalURL,
			IsDeleted:   false,
		})
	}

	return nil
}

func (m *MemoryStorage) Get(shortURL string) (dto.GetByIDDto, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	for _, url := range m.data {
		if url.ShortURL == shortURL {
			return dto.GetByIDDto{
				OriginalUrl: url.OriginalURL,
				IsDeleted:   url.IsDeleted,
			}, nil
		}
	}

	return dto.GetByIDDto{}, fmt.Errorf("URL not found for short URL: %s", shortURL)
}

func (m *MemoryStorage) GetUserURLs(userID uuid.UUID) ([]dto.UserURLsResponseDto, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	var result []dto.UserURLsResponseDto
	for _, entry := range m.data {
		if entry.UserID == userID && !entry.IsDeleted {
			result = append(result, dto.UserURLsResponseDto{
				ShortURL:    entry.ShortURL,
				OriginalURL: entry.OriginalURL,
			})
		}
	}

	return result, nil
}

func (m *MemoryStorage) DeleteUserURLs(userID uuid.UUID, shortURLs []string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	urlsToDelete := make(map[string]bool)
	for _, shortURL := range shortURLs {
		urlsToDelete[shortURL] = true
	}

	for i := range m.data {
		if m.data[i].UserID == userID && urlsToDelete[m.data[i].ShortURL] {
			m.data[i].IsDeleted = true
		}
	}

	return nil
}

func (m *MemoryStorage) LoadAll() (map[string]string, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	clone := make(map[string]string)
	for _, url := range m.data {
		if !url.IsDeleted {
			clone[url.ShortURL] = url.OriginalURL
		}
	}
	return clone, nil
}
