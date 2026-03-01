package model

import (
	"math/rand"
	"strings"
)

type UrlShortener struct {
	UrlID   string
	LongUrl string
}

func NewShortUrl(longUrl string) *UrlShortener {
	urlID := generatedUrlID()

	return &UrlShortener{
		UrlID:   urlID,
		LongUrl: longUrl,
	}
}

func generatedUrlID() string {

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
