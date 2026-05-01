package model

import (
	"math/rand"
	"strings"

	"github.com/google/uuid"
)

// Единая структура для хранения URL
type URLDto struct {
	UUID        uuid.UUID `json:"uuid"`
	ShortURL    string    `json:"short_url"`
	OriginalURL string    `json:"original_url"`
}

// Генерирует короткий ID (например, "aB3x9kLm")
func generateShortCode() string {
	chars := []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ" +
		"abcdefghijklmnopqrstuvwxyz" +
		"0123456789")
	length := 8
	var b strings.Builder
	for i := 0; i < length; i++ {
		b.WriteRune(chars[rand.Intn(len(chars))])
	}
	return b.String()
}

// NewShortURL создаёт новый URLDto с уникальным UUID и коротким URL
func NewShortURL(originalURL string) *URLDto {
	return &URLDto{
		UUID:        uuid.New(),
		ShortURL:    generateShortCode(),
		OriginalURL: originalURL,
	}
}
