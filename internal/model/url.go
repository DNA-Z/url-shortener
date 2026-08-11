// Package model содержит основную богатую модель, содержащую логику создания короткого URL
package model

import (
	"math/rand"

	"github.com/google/uuid"
)

type URL struct {
	UUID        uuid.UUID `json:"uuid"`
	UserID      uuid.UUID `json:"user_id"`
	ShortURL    string    `json:"short_url"`
	OriginalURL string    `json:"original_url"`
	IsDeleted   bool      `json:"is_deleted"`
}

func generateShortCode() string {
	const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"
	const length = 8

	b := make([]byte, length)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}

func NewShortURL(userID uuid.UUID, originalURL string) (*URL, error) {
	return &URL{
		UUID:        uuid.New(),
		UserID:      userID,
		ShortURL:    generateShortCode(),
		OriginalURL: originalURL,
		IsDeleted:   false,
	}, nil
}
