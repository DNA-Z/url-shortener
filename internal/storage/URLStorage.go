package storage

import (
	"github.com/DNA-Z/url-shortener/internal/dto"
	"github.com/DNA-Z/url-shortener/internal/model"
	"github.com/google/uuid"
)

type URLStorage interface {
	Save(url *model.URLDto) error
	Saves(urls []model.URLDto) error
	Get(shortURL string) (string, error)
	GetUserURLs(userID uuid.UUID) ([]dto.UserURLsResponseDto, error)
	LoadAll() (map[string]string, error)
}
