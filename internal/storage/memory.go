package storage

import (
	"fmt"
	"sync"

	"github.com/DNA-Z/url-shortener/internal/dto"
	"github.com/DNA-Z/url-shortener/internal/model"
	"github.com/google/uuid"
)

type MemoryStorage struct {
	urls map[string]model.URL
	mu   sync.RWMutex
}

func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{
		urls: make(map[string]model.URL),
	}
}

func (m *MemoryStorage) Save(url *model.URL) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.urls[url.ShortURL] = model.URL{
		UserID:      url.UserID,
		ShortURL:    url.ShortURL,
		OriginalURL: url.OriginalURL,
		IsDeleted:   false,
	}

	return nil
}

func (m *MemoryStorage) Saves(urls []model.URL) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	for _, url := range urls {
		m.urls[url.ShortURL] = model.URL{
			UserID:      url.UserID,
			ShortURL:    url.ShortURL,
			OriginalURL: url.OriginalURL,
			IsDeleted:   false,
		}
	}

	return nil
}

func (m *MemoryStorage) Get(shortURL string) (dto.GetByIDDto, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if url, exists := m.urls[shortURL]; exists {
		return dto.GetByIDDto{
			OriginalURL: url.OriginalURL,
			IsDeleted:   url.IsDeleted,
		}, nil
	}

	return dto.GetByIDDto{}, fmt.Errorf("URL not found for short URL: %s", shortURL)
}

func (m *MemoryStorage) GetUserURLs(userID uuid.UUID) ([]dto.UserURLsResponseDto, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	var result []dto.UserURLsResponseDto
	for _, url := range m.urls {
		if url.UserID == userID && !url.IsDeleted {
			result = append(result, dto.UserURLsResponseDto{
				ShortURL:    url.ShortURL,
				OriginalURL: url.OriginalURL,
			})
		}
	}

	return result, nil
}

func (m *MemoryStorage) DeleteUserURLs(userID uuid.UUID, shortURLs []string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	urlsToDelete := make(map[string]bool, len(shortURLs))
	for _, shortURL := range shortURLs {
		urlsToDelete[shortURL] = true
	}

	for shortURL := range urlsToDelete {
		if url, exists := m.urls[shortURL]; exists && url.UserID == userID {
			url.IsDeleted = true
			m.urls[shortURL] = url
		}
	}

	return nil
}

func (m *MemoryStorage) LoadAll() (map[string]string, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	clone := make(map[string]string, len(m.urls))
	for shortURL, url := range m.urls {
		if !url.IsDeleted {
			clone[shortURL] = url.OriginalURL
		}
	}
	return clone, nil
}

func (m *MemoryStorage) GetStats() (dto.StatsDto, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	var stats dto.StatsDto
	uniqueUsers := make(map[uuid.UUID]bool)

	for _, url := range m.urls {
		if !url.IsDeleted {
			stats.URLs++
			if url.UserID != uuid.Nil {
				uniqueUsers[url.UserID] = true
			}
		}
	}

	stats.Users = len(uniqueUsers)
	return stats, nil
}
