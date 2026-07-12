package service

import (
	"testing"

	"github.com/DNA-Z/url-shortener/internal/model"
	"github.com/DNA-Z/url-shortener/internal/storage"
	"github.com/google/uuid"
)

func BenchmarkShorten(b *testing.B) {
	store := storage.NewMemoryStorage()
	urlService := &URLStorage{
		storage: store,
		isDB:    false,
	}

	userID := uuid.New()
	originalURL := "https://example.com/very/long/url/that/needs/to/be/shortened"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = urlService.Shorten(userID, originalURL)
	}
}

func BenchmarkGetOriginURL(b *testing.B) {

	store := storage.NewMemoryStorage()
	urlService := &URLStorage{
		storage: store,
		isDB:    false,
	}

	userID := uuid.New()
	originalURL := "https://example.com/test/url"
	shortURL, _ := urlService.Shorten(userID, originalURL)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = urlService.GetOriginURLByShortURL(shortURL)
	}
}

func BenchmarkNewShortURL(b *testing.B) {
	userID := uuid.New()
	originalURL := "https://example.com/benchmark/url"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = model.NewShortURL(userID, originalURL)
	}
}

func BenchmarkShortenDifferentSizes(b *testing.B) {
	sizes := []int{10, 50, 100, 500, 1000}

	for _, size := range sizes {
		b.Run("Size_"+string(rune(size)), func(b *testing.B) {

			store := storage.NewMemoryStorage()
			urlService := &URLStorage{
				storage: store,
				isDB:    false,
			}

			userID := uuid.New()
			url := "https://example.com/" + string(make([]byte, size))

			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_, _ = urlService.Shorten(userID, url)
			}
		})
	}
}
