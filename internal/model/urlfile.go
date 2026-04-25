package model

import "github.com/google/uuid"

type URLFile struct {
	UUID        uuid.UUID `json:"uuid"`
	ShortURL    string    `json:"short_url"`
	OriginalURL string    `json:"original_url"`
}

func NewURLFile(shortURL string, originalURL string) (*URLFile, error) {
	urlFileID, err := uuid.NewUUID()
	if err != nil {
		return nil, err
	}

	return &URLFile{
		UUID:        urlFileID,
		ShortURL:    shortURL,
		OriginalURL: originalURL,
	}, nil
}
