package storage

import (
	"sync"

	"github.com/DNA-Z/url-shortener/internal/model"
)

type MemoryStorage struct {
	data map[string]string
	mu   sync.RWMutex
}

func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{
		data: make(map[string]string),
	}
}

func (m *MemoryStorage) Save(url *model.URLDto) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.data[url.ShortURL] = url.OriginalURL
	return nil
}

func (m *MemoryStorage) Get(shortURL string) (string, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	url, exists := m.data[shortURL]
	return url, exists
}

func (m *MemoryStorage) LoadAll() (map[string]string, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	clone := make(map[string]string)
	for k, v := range m.data {
		clone[k] = v
	}
	return clone, nil
}
