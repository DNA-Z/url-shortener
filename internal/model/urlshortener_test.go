package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewShortURL(t *testing.T) {
	t.Run("creates a new URLShortener with non-empty URLID and correct LongURL", func(t *testing.T) {
		longURL := "https://example.com/very/long/path"
		shortener := NewShortURL(longURL)

		assert.NotEmpty(t, shortener.URLID)
		assert.Equal(t, longURL, shortener.LongURL)
		assert.Len(t, shortener.URLID, 8) // Мы ожидаем длину 8 символов
	})

	t.Run("each call returns different URLID for same LongURL", func(t *testing.T) {
		longURL := "https://example.com"

		first := NewShortURL(longURL)
		second := NewShortURL(longURL)

		assert.Equal(t, longURL, first.LongURL)
		assert.Equal(t, longURL, second.LongURL)
		assert.NotEqual(t, first.URLID, second.URLID, "Each generated ID should be unique")
	})

	t.Run("handles empty LongURL correctly", func(t *testing.T) {
		shortener := NewShortURL("")

		assert.NotEmpty(t, shortener.URLID)
		assert.Empty(t, shortener.LongURL)
	})
}

func TestGeneratedURLID(t *testing.T) {
	t.Run("generates ID of exactly 8 characters", func(t *testing.T) {
		id := generatedURLID()
		assert.Len(t, id, 8)
	})

	t.Run("generated ID contains only valid characters (alphanumeric)", func(t *testing.T) {
		id := generatedURLID()
		validChars := "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"
		validSet := make(map[rune]bool)
		for _, c := range validChars {
			validSet[c] = true
		}

		for _, c := range id {
			assert.True(t, validSet[c], "Character %q in ID %q is not in the allowed set", c, id)
		}
	})

	t.Run("multiple calls generate different IDs", func(t *testing.T) {
		ids := make(map[string]bool)
		const count = 100
		for i := 0; i < count; i++ {
			id := generatedURLID()
			assert.False(t, ids[id], "Generated ID collision detected: %s", id)
			ids[id] = true
		}
		assert.Len(t, ids, count, "Expected %d unique IDs", count)
	})
}
