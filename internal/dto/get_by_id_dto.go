package dto

type GetByIDDto struct {
	OriginalURL string `json:"original-url"`
	IsDeleted   bool   `json:"is_deleted"`
}
