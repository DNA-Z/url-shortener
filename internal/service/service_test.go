package service

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestURL_Shorten(t *testing.T) {
	service := NewURL()

	t.Run("shortens a new URL and returns unique ID", func(t *testing.T) {
		url := "https://example.com"
		id := service.Shorten(url)

		assert.NotEmpty(t, id)
		longURL, err := service.GetByID(id)
		require.NoError(t, err)
		assert.Equal(t, url, longURL)
	})

	t.Run("returns the same ID when shortening the same URL twice", func(t *testing.T) {
		url := "https://golang.org"
		firstID := service.Shorten(url)
		secondID := service.Shorten(url)

		assert.Equal(t, firstID, secondID)
	})
}

func TestURL_GetByID(t *testing.T) {
	service := NewURL()
	id := service.Shorten("https://github.com")

	t.Run("returns original URL for valid ID", func(t *testing.T) {
		result, err := service.GetByID(id)
		require.NoError(t, err)
		assert.Equal(t, "https://github.com", result)
	})

	t.Run("returns error for non-existent ID", func(t *testing.T) {
		result, err := service.GetByID("unknown123")
		assert.Error(t, err)
		assert.Empty(t, result)
		assert.Equal(t, "URL unknown123 not found", err.Error())
	})
}

func TestURL_urlExists(t *testing.T) {
	service := NewURL()
	service.Shorten("https://reused.com")
	service.Shorten("https://new.com")

	t.Run("returns true and correct ID if URL already exists", func(t *testing.T) {
		exists, id := service.urlExists("https://reused.com")
		assert.True(t, exists)
		assert.NotEmpty(t, id)

		// Проверим, что по этому ID действительно лежит нужный URL
		longURL, err := service.GetByID(id)
		require.NoError(t, err)
		assert.Equal(t, "https://reused.com", longURL)
	})

	t.Run("returns false and empty ID if URL does not exist", func(t *testing.T) {
		exists, id := service.urlExists("https://notfound.com")
		assert.False(t, exists)
		assert.Empty(t, id)
	})
}
