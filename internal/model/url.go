package model

import (
	"math/rand"
	"strings"

	"github.com/google/uuid"
)

type URL struct {
	UUID uuid.UUID `json:"uuid"`
	//UserID      uuid.UUID `json:"user_id"`
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
	//IsDeleted   bool      `json:"is_deleted"`
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

func NewShortURL(originalURL string) (*URL, error) {
	return &URL{
		UUID: uuid.New(),
		//UserID:      userID,
		ShortURL:    generateShortCode(),
		OriginalURL: originalURL,
		//IsDeleted:   false,
	}, nil
}
