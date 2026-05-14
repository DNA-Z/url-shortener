package storage

import (
	"github.com/DNA-Z/url-shortener/internal/model"
)

type URLStorage interface {
	Save(url *model.URL) error
	Saves(urls []model.URL) error
	Get(shortURL string) (string, error)
	//GetUserURLs(userID uuid.UUID) ([]dto.UserURLsResponseDto, error)
	LoadAll() (map[string]string, error)
}
