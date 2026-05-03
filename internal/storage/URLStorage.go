package storage

import (
	"github.com/DNA-Z/url-shortener/internal/model"
)

type URLStorage interface {
	Save(url *model.URLDto) error
	Saves(urls []model.URLDto) error
	Get(shortURL string) (string, error)
	LoadAll() (map[string]string, error)
}
