package storage

import (
	"github.com/DNA-Z/url-shortener/internal/dto"
	"github.com/DNA-Z/url-shortener/internal/model"
	"github.com/google/uuid"
)

type URLStorage interface {
	Save(url *model.URL) error
	Saves(urls []model.URL) error
	Get(shortURL string) (dto.GetByIDDto, error)
	GetUserURLs(userID uuid.UUID) ([]dto.UserURLsResponseDto, error)
	LoadAll() (map[string]string, error)
	DeleteUserURLs(userID uuid.UUID, shortURLs []string) error
	GetStats() (dto.StatsDto, error)
}
