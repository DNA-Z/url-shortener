package model

//
//import (
//	"testing"
//
//	"github.com/google/uuid"
//	"github.com/stretchr/testify/assert"
//)
//
//func TestNewShortURL(t *testing.T) {
//	t.Run("creates a new URLDto with non-empty ShortURL and correct OriginalURL", func(t *testing.T) {
//		originalURL := "https://example.com/very/long/path"
//		shortener, _ := NewShortURL(originalURL, uuid.Nil)
//
//		assert.NotEmpty(t, shortener.ShortURL)
//		assert.Equal(t, originalURL, shortener.OriginalURL)
//		assert.Len(t, shortener.ShortURL, 8) // Ожидаем длину 8 символов
//	})
//
//	t.Run("each call returns different ShortURL for same OriginalURL", func(t *testing.T) {
//		originalURL := "https://example.com"
//
//		first, _ := NewShortURL(originalURL, uuid.Nil)
//		second, _ := NewShortURL(originalURL, uuid.Nil)
//
//		assert.Equal(t, originalURL, first.OriginalURL)
//		assert.Equal(t, originalURL, second.OriginalURL)
//		assert.NotEqual(t, first.ShortURL, second.ShortURL, "Each generated ShortURL should be unique")
//	})
//
//	t.Run("handles empty OriginalURL correctly", func(t *testing.T) {
//		shortener, _ := NewShortURL("", uuid.Nil)
//
//		assert.NotEmpty(t, shortener.ShortURL)
//		assert.Empty(t, shortener.OriginalURL)
//	})
//}
//
//func TestGenerateShortCode(t *testing.T) {
//	t.Run("generates code of exactly 8 characters", func(t *testing.T) {
//		code := generateShortCode()
//		assert.Len(t, code, 8)
//	})
//
//	t.Run("generated code contains only valid characters (alphanumeric)", func(t *testing.T) {
//		code := generateShortCode()
//		validChars := "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"
//		validSet := make(map[rune]bool)
//		for _, c := range validChars {
//			validSet[c] = true
//		}
//
//		for _, c := range code {
//			assert.True(t, validSet[c], "Character %q in code %q is not in the allowed set", c, code)
//		}
//	})
//
//	t.Run("multiple calls generate different codes", func(t *testing.T) {
//		codes := make(map[string]bool)
//		const count = 100
//		for i := 0; i < count; i++ {
//			code := generateShortCode()
//			assert.False(t, codes[code], "Generated code collision detected: %s", code)
//			codes[code] = true
//		}
//		assert.Len(t, codes, count, "Expected %d unique codes", count)
//	})
//}
