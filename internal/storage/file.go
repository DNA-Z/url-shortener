package storage

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"sync"

	"github.com/DNA-Z/url-shortener/internal/dto"
	"github.com/DNA-Z/url-shortener/internal/model"
	"github.com/google/uuid"
)

type FileStorage struct {
	file    *os.File
	writer  *bufio.Writer
	scanner *bufio.Scanner
	data    []model.URL
	mu      sync.RWMutex
	path    string
}

func NewFileStorage(path string) (*FileStorage, error) {
	file, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return nil, err
	}

	storage := &FileStorage{
		file:    file,
		writer:  bufio.NewWriter(file),
		scanner: bufio.NewScanner(file),
		data:    make([]model.URL, 0),
		path:    path,
	}

	if err := storage.load(); err != nil {
		file.Close()
		return nil, err
	}

	return storage, nil
}

func (f *FileStorage) load() error {
	_, err := f.file.Seek(0, 0)
	if err != nil {
		return err
	}
	f.scanner = bufio.NewScanner(f.file)

	for f.scanner.Scan() {
		var urlFile model.URL
		if err := json.Unmarshal(f.scanner.Bytes(), &urlFile); err != nil {
			return err
		}
		f.data = append(f.data, model.URL{
			ShortURL:    urlFile.ShortURL,
			OriginalURL: urlFile.OriginalURL,
			UserID:      urlFile.UserID,
			IsDeleted:   urlFile.IsDeleted,
		})
	}
	return f.scanner.Err()
}

func (f *FileStorage) Save(url *model.URL) error {
	f.mu.Lock()
	defer f.mu.Unlock()

	data, err := json.Marshal(url)
	if err != nil {
		return err
	}

	_, err = f.writer.Write(data)
	if err != nil {
		return err
	}
	err = f.writer.WriteByte('\n')
	if err != nil {
		return err
	}
	err = f.writer.Flush()
	if err != nil {
		return err
	}

	f.data = append(f.data, model.URL{
		ShortURL:    url.ShortURL,
		OriginalURL: url.OriginalURL,
		UserID:      url.UserID,
		IsDeleted:   url.IsDeleted,
	})
	return nil
}

func (f *FileStorage) Saves(urls []model.URL) error {
	f.mu.Lock()
	defer f.mu.Unlock()

	for _, url := range urls {
		data, err := json.Marshal(url)
		if err != nil {
			return err
		}

		_, err = f.writer.Write(data)
		if err != nil {
			return err
		}

		err = f.writer.WriteByte('\n')
		if err != nil {
			return err
		}

		f.data = append(f.data, model.URL{
			ShortURL:    url.ShortURL,
			OriginalURL: url.OriginalURL,
			UserID:      url.UserID,
			IsDeleted:   url.IsDeleted,
		})
	}
	return f.writer.Flush()
}

func (f *FileStorage) Get(shortURL string) (dto.GetByIDDto, error) {
	f.mu.RLock()
	defer f.mu.RUnlock()

	for _, url := range f.data {
		if url.ShortURL == shortURL {
			return dto.GetByIDDto{
				OriginalURL: url.OriginalURL,
				IsDeleted:   url.IsDeleted,
			}, nil
		}
	}
	return dto.GetByIDDto{}, fmt.Errorf("URL not found for short URL: %s", shortURL)
}

func (f *FileStorage) GetUserURLs(userID uuid.UUID) ([]dto.UserURLsResponseDto, error) {
	f.mu.RLock()
	defer f.mu.RUnlock()

	var result []dto.UserURLsResponseDto

	for _, url := range f.data {
		if url.UserID == userID && !url.IsDeleted {
			result = append(result, dto.UserURLsResponseDto{
				ShortURL:    url.ShortURL,
				OriginalURL: url.OriginalURL,
			})
		}
	}

	if len(result) == 0 {
		return nil, nil
	}
	return result, nil
}

func (f *FileStorage) DeleteUserURLs(userID uuid.UUID, shortURLs []string) error {
	f.mu.Lock()
	defer f.mu.Unlock()

	urlsToDelete := make(map[string]bool)
	for _, shortURL := range shortURLs {
		urlsToDelete[shortURL] = true
	}

	for i, url := range f.data {
		if url.UserID == userID && urlsToDelete[url.ShortURL] && !url.IsDeleted {
			f.data[i].IsDeleted = true
		}
	}

	return f.saveAllToFile()
}

func (f *FileStorage) LoadAll() (map[string]string, error) {
	f.mu.RLock()
	defer f.mu.RUnlock()

	clone := make(map[string]string)
	for _, url := range f.data {
		if !url.IsDeleted {
			clone[url.ShortURL] = url.OriginalURL
		}
	}
	return clone, nil
}

func (f *FileStorage) saveAllToFile() error {
	if err := f.file.Truncate(0); err != nil {
		return err
	}
	if _, err := f.file.Seek(0, 0); err != nil {
		return err
	}

	f.writer = bufio.NewWriter(f.file)

	for _, url := range f.data {
		data, err := json.Marshal(url)
		if err != nil {
			return err
		}

		if _, err := f.writer.Write(data); err != nil {
			return err
		}
		if err := f.writer.WriteByte('\n'); err != nil {
			return err
		}
	}

	return f.writer.Flush()
}
