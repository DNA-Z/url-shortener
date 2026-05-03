package storage

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"sync"

	"github.com/DNA-Z/url-shortener/internal/model"
)

type FileStorage struct {
	file    *os.File
	writer  *bufio.Writer
	scanner *bufio.Scanner
	data    map[string]string
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
		data:    make(map[string]string),
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
		var urlFile model.URLDto
		if err := json.Unmarshal(f.scanner.Bytes(), &urlFile); err != nil {
			return err
		}
		f.data[urlFile.ShortURL] = urlFile.OriginalURL
	}
	return f.scanner.Err()
}

func (f *FileStorage) Save(url *model.URLDto) error {
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

	f.data[url.ShortURL] = url.OriginalURL
	return nil
}

func (f *FileStorage) Saves(urls []model.URLDto) error {
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

		err = f.writer.WriteByte('n')
		if err != nil {
			return err
		}

		f.data[url.ShortURL] = url.OriginalURL
	}
	return f.writer.Flush()
}

func (f *FileStorage) Get(shortURL string) (string, error) {
	f.mu.RLock()
	defer f.mu.RUnlock()
	url, exists := f.data[shortURL]
	if exists == false {
		return "", fmt.Errorf("URL not found for short URL: %s", shortURL)
	}
	return url, nil
}

func (f *FileStorage) LoadAll() (map[string]string, error) {
	f.mu.RLock()
	defer f.mu.RUnlock()
	clone := make(map[string]string)
	for k, v := range f.data {
		clone[k] = v
	}
	return clone, nil
}
