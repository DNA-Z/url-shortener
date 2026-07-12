package audit

import (
	"encoding/json"
	"os"
	"sync"
)

type FileObserver struct {
	file   *os.File
	mu     sync.Mutex
	closed bool
}

func NewFileObserver(filePath string) (*FileObserver, error) {
	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, err
	}

	return &FileObserver{
		file:   file,
		closed: false,
	}, nil
}

func (f *FileObserver) Notify(event AuditEvent) error {
	f.mu.Lock()
	defer f.mu.Unlock()

	if f.closed {
		return nil
	}

	data, err := json.Marshal(event)
	if err != nil {
		return err
	}

	_, err = f.file.Write(append(data, '\n'))
	return err
}

func (f *FileObserver) Close() error {
	f.mu.Lock()
	defer f.mu.Unlock()

	if f.closed {
		return nil
	}

	f.closed = true
	return f.file.Close()
}
