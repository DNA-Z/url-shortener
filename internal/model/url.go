package model

import (
	"math/rand"
	"strings"

	"github.com/google/uuid"
)

type URLDto struct {
	UUID        uuid.UUID `json:"uuid"`
	ShortURL    string    `json:"short_url"`
	OriginalURL string    `json:"original_url"`
}

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

func NewShortURL(originalURL string) *URLDto {
	return &URLDto{
		UUID:        uuid.New(),
		ShortURL:    generateShortCode(),
		OriginalURL: originalURL,
	}
}
