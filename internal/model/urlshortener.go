package model

import (
	"math/rand"
	"strings"
)

type URLShortener struct {
	URLID   string
	LongURL string
}

func NewShortURL(longURL string) *URLShortener {
	urlID := generatedURLID()

	return &URLShortener{
		URLID:   urlID,
		LongURL: longURL,
	}
}

func generatedURLID() string {

	chars := []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ" +
		"abcdefghijklmnopqrstuvwxyz" +
		"0123456789")
	length := 8
	var b strings.Builder
	for i := 0; i < length; i++ {
		b.WriteRune(chars[rand.Intn(len(chars))])
	}
	ID := b.String()

	return ID
}
