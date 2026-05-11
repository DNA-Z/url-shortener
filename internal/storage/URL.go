package storage

import "github.com/google/uuid"

type URL struct {
	UserID      uuid.UUID
	ShortURL    string
	OriginalURL string
	IsDeleted   bool
}
